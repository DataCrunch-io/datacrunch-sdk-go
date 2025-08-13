package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/util"
)

// UnmarshalHandler is a named request handler for unmarshaling rest protocol requests
var UnmarshalHandler = request.NamedHandler{Name: "rest.Unmarshal", Fn: Unmarshal}

// UnmarshalMetaHandler is a named request handler for unmarshaling rest protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{Name: "rest.UnmarshalMeta", Fn: UnmarshalMeta}

// Unmarshal unmarshals the REST component of a response in a REST service.
func Unmarshal(r *request.Request) {
	if r.DataFilled() {
		v := reflect.Indirect(reflect.ValueOf(r.Data))
		if err := unmarshalBody(r, v); err != nil {
			r.Error = err
		}
	}
}

// UnmarshalMeta unmarshals the REST metadata of a response in a REST service
func UnmarshalMeta(r *request.Request) {

	if r.DataFilled() {
		if err := UnmarshalResponse(r.HTTPResponse, r.Data, false); err != nil {
			r.Error = err
		}
	}
}

// UnmarshalResponse attempts to unmarshal the REST response headers to
// the data type passed in. The type must be a pointer. An error is returned
// with any error unmarshaling the response into the target datatype.
func UnmarshalResponse(resp *http.Response, data interface{}, lowerCaseHeaderMaps bool) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	// Only unmarshal location elements for struct types
	if v.Kind() == reflect.Struct {
		return unmarshalLocationElements(resp, v, lowerCaseHeaderMaps)
	}
	// For non-struct types (like slices), there are no location elements to unmarshal
	return nil
}

func unmarshalBody(r *request.Request, v reflect.Value) error {
	if field, ok := v.Type().FieldByName("_"); ok {
		if payloadName := field.Tag.Get("payload"); payloadName != "" {
			pfield, _ := v.Type().FieldByName(payloadName)
			if ptag := pfield.Tag.Get("type"); ptag != "" && ptag != "structure" {
				payload := v.FieldByName(payloadName)
				if payload.IsValid() {
					switch payload.Interface().(type) {
					case []byte:
						defer func() {
							if err := r.HTTPResponse.Body.Close(); err != nil {
								r.Error = dcerr.New(request.ErrCodeSerialization, "failed to close response body", err)
							}
						}()
						b, err := io.ReadAll(r.HTTPResponse.Body)
						if err != nil {
							return dcerr.New(request.ErrCodeSerialization, "failed to decode REST response", err)
						}

						payload.Set(reflect.ValueOf(b))

					case *string:
						defer func() {
							if err := r.HTTPResponse.Body.Close(); err != nil {
								r.Error = dcerr.New(request.ErrCodeSerialization, "failed to close response body", err)
							}
						}()
						b, err := io.ReadAll(r.HTTPResponse.Body)
						if err != nil {
							return dcerr.New(request.ErrCodeSerialization, "failed to decode REST response", err)
						}

						str := string(b)
						payload.Set(reflect.ValueOf(&str))

					default:
						switch payload.Type().String() {
						case "io.ReadCloser":
							payload.Set(reflect.ValueOf(r.HTTPResponse.Body))

						case "io.ReadSeeker":
							b, err := io.ReadAll(r.HTTPResponse.Body)
							if err != nil {
								return dcerr.New(request.ErrCodeSerialization,
									"failed to read response body", err)
							}
							payload.Set(reflect.ValueOf(io.NopCloser(bytes.NewReader(b))))

						default:
							if _, err := io.Copy(io.Discard, r.HTTPResponse.Body); err != nil {
								// Log the error but continue with cleanup
								_ = err // Suppress unused variable warning
							}
							if err := r.HTTPResponse.Body.Close(); err != nil {
								// Log the error but continue with cleanup
								_ = err // Suppress unused variable warning
							}
							return dcerr.New(request.ErrCodeSerialization,
								"failed to decode REST response",
								fmt.Errorf("unknown payload type %s", payload.Type()))
						}
					}
				}
			}
		}
	}

	return nil
}

func unmarshalLocationElements(resp *http.Response, v reflect.Value, lowerCaseHeaderMaps bool) error {
	for i := 0; i < v.NumField(); i++ {
		m, field := v.Field(i), v.Type().Field(i)
		if n := field.Name; n[0:1] == strings.ToLower(n[0:1]) {
			continue
		}

		if m.IsValid() {
			name := field.Tag.Get("locationName")
			if name == "" {
				name = field.Name
			}

			switch field.Tag.Get("location") {
			case "statusCode":
				unmarshalStatusCode(m, resp.StatusCode)

			case "header":
				err := unmarshalHeader(m, resp.Header.Get(name), field.Tag)
				if err != nil {
					return dcerr.New(request.ErrCodeSerialization, "failed to decode REST response", err)
				}

			case "headers":
				prefix := field.Tag.Get("locationName")
				err := unmarshalHeaderMap(m, resp.Header, prefix, lowerCaseHeaderMaps)
				if err != nil {
					return dcerr.New(request.ErrCodeSerialization, "failed to decode REST response", err)
				}
			}
		}
	}

	return nil
}

func unmarshalStatusCode(v reflect.Value, statusCode int) {
	if !v.IsValid() {
		return
	}

	switch v.Interface().(type) {
	case *int64:
		s := int64(statusCode)
		v.Set(reflect.ValueOf(&s))
	}
}

func unmarshalHeaderMap(r reflect.Value, headers http.Header, prefix string, normalize bool) error {
	if len(headers) == 0 {
		return nil
	}
	switch r.Interface().(type) {
	case map[string]*string: // we only support string map value types
		out := map[string]*string{}
		for k, v := range headers {
			if util.HasPrefixFold(k, prefix) {
				if normalize {
					k = strings.ToLower(k)
				} else {
					k = http.CanonicalHeaderKey(k)
				}
				out[k[len(prefix):]] = &v[0]
			}
		}
		if len(out) != 0 {
			r.Set(reflect.ValueOf(out))
		}

	}
	return nil
}

func unmarshalHeader(v reflect.Value, header string, tag reflect.StructTag) error {
	switch tag.Get("type") {
	case "jsonvalue":
		if len(header) == 0 {
			return nil
		}
	case "blob":
		if len(header) == 0 {
			return nil
		}
	default:
		if !v.IsValid() || (header == "" && (!v.IsNil() && v.Elem().Kind() != reflect.String)) {
			return nil
		}
	}

	switch v.Interface().(type) {
	case *string:
		if tag.Get("suppressedJSONValue") == "true" && tag.Get("location") == "header" {
			b, err := base64.StdEncoding.DecodeString(header)
			if err != nil {
				return fmt.Errorf("failed to decode JSONValue, %v", err)
			}
			header = string(b)
		}
		v.Set(reflect.ValueOf(&header))
	case []byte:
		b, err := base64.StdEncoding.DecodeString(header)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(b))
	case *bool:
		b, err := strconv.ParseBool(header)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(&b))
	case *int64:
		i, err := strconv.ParseInt(header, 10, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(&i))
	case *float64:
		var f float64
		switch {
		case strings.EqualFold(header, floatNaN):
			f = math.NaN()
		case strings.EqualFold(header, floatInf):
			f = math.Inf(1)
		case strings.EqualFold(header, floatNegInf):
			f = math.Inf(-1)
		default:
			var err error
			f, err = strconv.ParseFloat(header, 64)
			if err != nil {
				return err
			}
		}
		v.Set(reflect.ValueOf(&f))
	case *time.Time:
		format := "2006-01-02T15:04:05Z" // default to ISO8601
		t, err := time.Parse(format, header)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(&t))
	case map[string]interface{}:
		escaping := protocol.NoEscape
		if tag.Get("location") == "header" {
			escaping = protocol.Base64Escape
		}
		m, err := protocol.DecodeJSONValue(header, escaping)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(m))
	default:
		err := fmt.Errorf("unsupported value for param %v (%s)", v.Interface(), v.Type())
		return err
	}
	return nil
}
