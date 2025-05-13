package sprout

import "encoding/json"

type Response struct {
	Success bool   `json:"isSuccess"`
	Message string `json:"message"`
}

func NewResponse(body []byte) (*Response, error) {
	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
