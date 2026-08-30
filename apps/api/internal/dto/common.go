package dto

type EmptyPayload struct{}

func (p *EmptyPayload) Validate() error { return nil }
