package converter

import (
	"encoding/json"
	"fmt"

	commonpb "go.temporal.io/api/common/v1"
)

const (
	// maxJSONPayloadSize is the maximum size for a JSON payload during deserialization
	// to prevent memory exhaustion attacks. This is set to 10MB as a reasonable default
	// while the gRPC max payload size is 128MB.
	maxJSONPayloadSize = 10 * 1024 * 1024 // 10MB
)

// JSONPayloadConverter converts to/from JSON.
type JSONPayloadConverter struct {
}

// NewJSONPayloadConverter creates a new instance of JSONPayloadConverter.
func NewJSONPayloadConverter() *JSONPayloadConverter {
	return &JSONPayloadConverter{}
}

// ToPayload converts a single value to a payload.
func (c *JSONPayloadConverter) ToPayload(value interface{}) (*commonpb.Payload, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnableToEncode, err)
	}
	return newPayload(data, c), nil
}

// FromPayload converts a single payload to a value.
func (c *JSONPayloadConverter) FromPayload(payload *commonpb.Payload, valuePtr interface{}) error {
	data := payload.GetData()
	// Security: Check payload size to prevent memory exhaustion
	if len(data) > maxJSONPayloadSize {
		return fmt.Errorf("%w: payload size %d exceeds maximum allowed size %d",
			ErrUnableToDecode, len(data), maxJSONPayloadSize)
	}
	err := json.Unmarshal(data, valuePtr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnableToDecode, err)
	}
	return nil
}

// ToString converts a payload object into a human-readable string.
func (c *JSONPayloadConverter) ToString(payload *commonpb.Payload) string {
	return string(payload.GetData())
}

// Encoding returns MetadataEncodingJSON.
func (c *JSONPayloadConverter) Encoding() string {
	return MetadataEncodingJSON
}
