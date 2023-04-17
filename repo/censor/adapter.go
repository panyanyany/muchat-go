package censor

type Adapter interface {
	MakeTextAuditing(id string, text string) (*TextAuditingResult, error)
}

type TextAuditingResult struct {
	Safe         bool
	FilteredText string
}
