package dtos

type TypeDto struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Intensity *int   `json:"intensity,omitempty"`
}
