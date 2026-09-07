package dto

type EmptyPayload struct{}

func (p *EmptyPayload) Validate() error { return nil }

type MessageResponse struct {
	Message string `json:"message"`
}
