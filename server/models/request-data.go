package models

import "encoding/json"

type RequestData struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data" default:""`
}

func (r *RequestData) GetData(data interface{}) error {
	err := json.Unmarshal(r.Data, &data)
	if err != nil {
		return err
	}
	return nil
}
