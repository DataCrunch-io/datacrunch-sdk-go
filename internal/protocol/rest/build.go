package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol"
)

const (
	floatNaN    = "NaN"
	floatInf    = "Infinity"
	floatNegInf = "-Infinity"
)

// Whether the byte value can be sent without escaping in Datacrunch URLs
var noEscape [256]bool

var errValueNotSet = fmt.Errorf("value not set")

var byteSliceType = reflect.TypeOf([]byte{})

func init() {
	for i := 0; i < len(noEscape); i++ {
		// Datacrunch expects every character except these to be escaped
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '-' ||
			i == '.' ||
			i == '_' ||
			i == '~'
	}
}

// BuildHandler is a named request handler for building rest protocol requests
var BuildHandler = request.NamedHandler{Name: "rest.Build", Fn: Build}

// Build builds the REST component of a service request.
func Build(r *request.Request) {
	if r.ParamsFilled() {
		v := reflect.ValueOf(r.Params).Elem()
		buildLocationElements(r, v, false)
		buildBody(r, v)
	}
}

// BuildAsGET builds the REST component of a service request with the ability to hoist
// data from the body.
func BuildAsGET(r *request.Request) {
	if r.ParamsFilled() {
		v := reflect.ValueOf(r.Params).Elem()
		buildLocationElements(r, v, true)
		buildBody(r, v)
	}
}

func buildLocationElements(r *request.Request, v reflect.Value, buildGETQuery bool) {
	query := r.HTTPRequest.URL.Query()

	// Setup the raw path to match the base path pattern. This is needed
	// so that when the path is mutated a custom escaped version can be
	// stored in RawPath that will be used by the Go client.
	r.HTTPRequest.URL.RawPath = r.HTTPRequest.URL.Path

	for i := 0; i < v.NumField(); i++ {
		m := v.Field(i)
		if n := v.Type().Field(i).Name; n[0:1] == strings.ToLower(n[0:1]) {
			continue
		}

		if m.IsValid() {
			field := v.Type().Field(i)
			name := field.Tag.Get("locationName")
			if name == "" {
				name = field.Name
			}
			if kind := m.Kind(); kind == reflect.Ptr {
				m = m.Elem()
			} else if kind == reflect.Interface {
				if !m.Elem().IsValid() {
					continue
				}
			}
			if !m.IsValid() {
				continue
			}
			if field.Tag.Get("ignore") != "" {
				continue
			}

			var err error
			switch field.Tag.Get("location") {
			case "headers": // header maps
				err = buildHeaderMap(&r.HTTPRequest.Header, m, field.Tag)
			case "header":
				err = buildHeader(&r.HTTPRequest.Header, m, name, field.Tag)
			case "uri":
				err = buildURI(r.HTTPRequest.URL, m, name, field.Tag)
			case "querystring":
				err = buildQueryString(query, m, name, field.Tag)
			default:
				if buildGETQuery {
					err = buildQueryString(query, m, name, field.Tag)
				}
			}
			r.Error = err
		}
		if r.Error != nil {
			return
		}
	}

	r.HTTPRequest.URL.RawQuery = query.Encode()
}

func buildBody(r *request.Request, v reflect.Value) {
	params := v.Interface()
	if params == nil {
		return
	}

	switch body := params.(type) {
	case io.ReadSeeker:
		r.SetReaderBody(body)
	case []byte:
		r.SetBufferBody(body)
	case string:
		r.SetStringBody(body)
	default:
		// JSON marshal everything else
		if data, err := json.Marshal(params); err == nil {
			r.SetBufferBody(data)
			r.HTTPRequest.Header.Set("Content-Type", "application/json")
		} else {
			r.Error = dcerr.New(request.ErrCodeSerialization,
				"failed to encode REST request", err)
		}
	}
}

func buildHeader(header *http.Header, v reflect.Value, name string, tag reflect.StructTag) error {
	str, err := convertType(v, tag)
	if err == errValueNotSet {
		return nil
	} else if err != nil {
		return dcerr.New(request.ErrCodeSerialization, "failed to encode REST request", err)
	}

	name = strings.TrimSpace(name)
	str = strings.TrimSpace(str)

	header.Add(name, str)

	return nil
}

func buildHeaderMap(header *http.Header, v reflect.Value, tag reflect.StructTag) error {
	prefix := tag.Get("locationName")
	for _, key := range v.MapKeys() {
		str, err := convertType(v.MapIndex(key), tag)
		if err == errValueNotSet {
			continue
		} else if err != nil {
			return dcerr.New(request.ErrCodeSerialization, "failed to encode REST request", err)

		}
		keyStr := strings.TrimSpace(key.String())
		str = strings.TrimSpace(str)

		header.Add(prefix+keyStr, str)
	}
	return nil
}

func buildURI(u *url.URL, v reflect.Value, name string, tag reflect.StructTag) error {
	value, err := convertType(v, tag)
	if err == errValueNotSet {
		return nil
	} else if err != nil {
		return dcerr.New(request.ErrCodeSerialization, "failed to encode REST request", err)
	}

	u.Path = strings.ReplaceAll(u.Path, "{"+name+"}", EscapePath(value))

	u.RawPath = strings.ReplaceAll(u.RawPath, "{"+name+"}", EscapePath(value))

	return nil
}

func buildQueryString(query url.Values, v reflect.Value, name string, tag reflect.StructTag) error {
	switch value := v.Interface().(type) {
	case []*string:
		for _, item := range value {
			query.Add(name, *item)
		}
	case map[string]*string:
		for key, item := range value {
			query.Add(key, *item)
		}
	case map[string][]*string:
		for key, items := range value {
			for _, item := range items {
				query.Add(key, *item)
			}
		}
	default:
		str, err := convertType(v, tag)
		if err == errValueNotSet {
			return nil
		} else if err != nil {
			return dcerr.New(request.ErrCodeSerialization, "failed to encode REST request", err)
		}
		query.Set(name, str)
	}

	return nil
}

func EscapePath(path string) string {
	return url.PathEscape(path)
}

func convertType(v reflect.Value, tag reflect.StructTag) (str string, err error) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return "", errValueNotSet
	}

	switch value := v.Interface().(type) {
	case string:
		// if the value is a string and the tag has suppressedJSONValue=true and location=header, then encode the value to base64
		if tag.Get("suppressedJSONValue") == "true" && tag.Get("location") == "header" {
			value = base64.StdEncoding.EncodeToString([]byte(value))
		}
		str = value
	case []*string:
		if tag.Get("location") != "header" || tag.Get("enum") == "" {
			return "", fmt.Errorf("%T is only supported with location header and enum shapes", value)
		}
		if len(value) == 0 {
			return "", errValueNotSet
		}

		buff := &bytes.Buffer{}
		for i, sv := range value {
			if sv == nil || len(*sv) == 0 {
				continue
			}
			if i != 0 {
				buff.WriteRune(',')
			}
			item := *sv
			if strings.Contains(item, `,`) || strings.Contains(item, `"`) {
				item = strconv.Quote(item)
			}
			buff.WriteString(item)
		}
		str = buff.String()
	case []byte:
		str = base64.StdEncoding.EncodeToString(value)
	case bool:
		str = strconv.FormatBool(value)
	case int:
		str = strconv.FormatInt(int64(value), 10)
	case int64:
		str = strconv.FormatInt(value, 10)
	case float64:
		switch {
		case math.IsNaN(value):
			str = floatNaN
		case math.IsInf(value, 1):
			str = floatInf
		case math.IsInf(value, -1):
			str = floatNegInf
		default:
			str = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case time.Time:
		format := tag.Get("timestampFormat")
		if len(format) == 0 {
			format = "2006-01-02T15:04:05Z" // default to ISO8601
		}
		str = value.Format(format)
	case map[string]interface{}:
		if len(value) == 0 {
			return "", errValueNotSet
		}
		escaping := protocol.NoEscape
		if tag.Get("location") == "header" {
			escaping = protocol.Base64Escape
		}
		str, err = protocol.EncodeJSONValue(value, escaping)
		if err != nil {
			return "", fmt.Errorf("unable to encode JSONValue, %v", err)
		}
	default:
		err := fmt.Errorf("unsupported value for param %v (%s)", v.Interface(), v.Type())
		return "", err
	}

	return str, nil
}
