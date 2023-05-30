package queue

import "encoding/json"

type QueueDTO struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	ID       int    `json:"id"`
}

func (q *QueueDTO) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

func (q *QueueDTO) Unmarshal(data []byte) error {
	return json.Unmarshal(data, q)
}
