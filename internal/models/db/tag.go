package db

import (
	"encoding/json"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// TODO comments
type TagDBResponse struct {
	Key        string `db:"key"`
	Value      string `db:"value"`
	InstanceID string `db:"instance_id"`
}

// TODO comments
func (t TagDBResponse) ToTagDTOResponse() *dto.TagDTOResponse {
	return &dto.TagDTOResponse{
		Key:        t.Key,
		Value:      t.Value,
		InstanceID: t.InstanceID,
	}
}

// TagDBResponses implements sql.Scanner interface for processing tags as JSONB
type TagDBResponses []TagDBResponse

// Scan for implementing sql.Scanner when reading instances from DB as InstanceDBResponse
func (t *TagDBResponses) Scan(value interface{}) error {
	if value == nil {
		*t = TagDBResponses{}
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	}

	// Expect array of objects with {key,value}
	if err := json.Unmarshal(b, t); err == nil {
		return nil
	}

	// Optional robustness: accept map[string]string and convert to KV slice
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("TagsKVList: invalid JSON: %w", err)
	}

	out := make(TagDBResponses, 0, len(m))
	for k, v := range m {
		out = append(out, TagDBResponse{Key: k, Value: fmt.Sprint(v)})
	}

	*t = out

	return nil
}
