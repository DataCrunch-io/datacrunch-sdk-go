// Example: JSON service using global logger
package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// JSONService - handles JSON operations
type JSONService struct{}

func NewJSONService() *JSONService {
	return &JSONService{}
}

// UnmarshalResponse - unmarshal with logging
func (j *JSONService) UnmarshalResponse(data []byte, target interface{}) error {
	targetType := reflect.TypeOf(target)
	
	logger.Debug("Starting JSON unmarshal",
		"target_type", targetType.String(),
		"data_size", len(data),
	)

	start := time.Now()
	
	err := json.Unmarshal(data, target)
	duration := time.Since(start)
	
	if err != nil {
		logger.Error("JSON unmarshal failed",
			"target_type", targetType.String(),
			"error", err,
			"data_sample", j.getSafeDataSample(data),
		)
		return err
	}

	// Log success with performance metrics
	logger.Debug("JSON unmarshal completed",
		"target_type", targetType.String(),
		"duration", duration,
		"data_size", len(data),
		"parse_speed_mb_per_sec", j.calculateParseSpeed(len(data), duration),
	)

	// Log field mapping info for structs
	if targetType.Kind() == reflect.Ptr && targetType.Elem().Kind() == reflect.Struct {
		j.logFieldMappingInfo(targetType.Elem(), target)
	}

	return nil
}

func (j *JSONService) getSafeDataSample(data []byte) string {
	const maxSampleSize = 200
	sample := string(data)
	
	if len(sample) > maxSampleSize {
		sample = sample[:maxSampleSize] + "... (truncated)"
	}
	
	return sample
}

func (j *JSONService) calculateParseSpeed(bytes int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	return float64(bytes) / duration.Seconds() / (1024 * 1024) // MB/s
}

func (j *JSONService) logFieldMappingInfo(structType reflect.Type, target interface{}) {
	// Count populated fields
	populated := 0
	total := structType.NumField()
	
	targetValue := reflect.ValueOf(target).Elem()
	
	for i := 0; i < total; i++ {
		field := targetValue.Field(i)
		if !field.IsZero() {
			populated++
		}
	}
	
	logger.Debug("Struct field population",
		"struct_type", structType.String(),
		"populated_fields", populated,
		"total_fields", total,
		"population_rate", float64(populated)/float64(total)*100,
	)
}

// Example with error handling
func (j *JSONService) UnmarshalInstanceResponse(jsonData string) (*Instance, error) {
	logger.Info("Unmarshaling instance response")
	
	var instance Instance
	err := j.UnmarshalResponse([]byte(jsonData), &instance)
	if err != nil {
		return nil, err
	}
	
	// Validate result
	if instance.ID == "" {
		logger.Warn("Unmarshaled instance missing ID field")
	}
	
	logger.Info("Instance unmarshaling completed",
		"instance_id", instance.ID,
		"instance_name", instance.Name,
	)
	
	return &instance, nil
}

// Example with array response
func (j *JSONService) UnmarshalInstanceListResponse(jsonData string) ([]*Instance, error) {
	logger.Info("Unmarshaling instance list response")
	
	var instances []*Instance
	err := j.UnmarshalResponse([]byte(jsonData), &instances)
	if err != nil {
		return nil, err
	}
	
	logger.Info("Instance list unmarshaling completed",
		"instance_count", len(instances),
	)
	
	return instances, nil
}

// Example usage and testing
func (j *JSONService) ExampleUsage() {
	// Valid JSON
	validJSON := `{"ID": "inst-123", "Name": "my-server", "Type": "small"}`
	instance, err := j.UnmarshalInstanceResponse(validJSON)
	if err == nil {
		logger.Info("Successfully parsed instance", "instance", instance)
	}
	
	// Invalid JSON (will show error logging)
	invalidJSON := `{"ID": "inst-123", "Name": "my-server", "Type":}` // Missing value
	_, err = j.UnmarshalInstanceResponse(invalidJSON)
	if err != nil {
		// Error already logged by UnmarshalResponse
		logger.Info("Expected error occurred during invalid JSON parsing")
	}
	
	// Array response
	arrayJSON := `[{"ID": "inst-1", "Name": "server-1"}, {"ID": "inst-2", "Name": "server-2"}]`
	instances, err := j.UnmarshalInstanceListResponse(arrayJSON)
	if err == nil {
		logger.Info("Successfully parsed instance list", "count", len(instances))
	}
}