package entity

type ResponseBodyTemplate struct {
	// RegexNoMatchMsgs []string `json:"regex_no_match,omitempty"`
	RegexNoMatchMsgs        *map[string][]string `json:"regex_no_match,omitempty"`
	Message                 string               `json:"message,omitempty"`
	Errs                    []error              `json:"errors,omitempty"`
	FlattenRegexNoMatchMsgs []string             `json:"flatten_regex_no_match,omitempty"`
}
