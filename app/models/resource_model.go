package models

type Resource struct {
	Id     string                 `json:"id"`
	Type   string                 `json:"type"`
	Parent string                 `json:"parent"`
	Data   map[string]interface{} `json:"data"`
}
