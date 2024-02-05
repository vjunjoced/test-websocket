package models

import "encoding/json"

type ResponseData struct {
	Action  string          `json:"action"`
	Data    json.RawMessage `json:"data,omitempty"`
	Success bool            `json:"success,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func (r *ResponseData) SetData(data interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Data = raw
	return nil
}
