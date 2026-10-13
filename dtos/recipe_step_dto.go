package dtos

type RecipeStepDto struct {
	StepNo   int     `json:"stepNo"`
	Section  *string `json:"section,omitempty"`
	Text     string  `json:"text"`
	Prep     bool    `json:"prep"`
	PrepLead *string `json:"prepLead,omitempty"`
}
