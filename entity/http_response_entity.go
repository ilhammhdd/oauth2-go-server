package entity

import "github.com/ilhammhdd/go-toolkit/errorkit"

type ResponseBodyTemplate struct {
	RegexNoMatchMsgs     *map[string][]string      `json:"regex_no_match,omitempty"`
	Message              *string                   `json:"message,omitempty"`
	Errs                 []error                   `json:"errors,omitempty"`
	DetailedErrs         []*errorkit.DetailedError `json:"detailed_errors,omitempty"`
	FlatRegexNoMatchMsgs []string                  `json:"flatten_regex_no_match,omitempty"`
}

func (rbt *ResponseBodyTemplate) FlattenErrorPageDataRegexNoMatchMsgs(errorPageDataID uint64) (uint32, []interface{}) {
	if rbt.RegexNoMatchMsgs == nil {
		return 0, nil
	}
	var flatRegexNoMatchMsgs []interface{}
	var flatLength uint32
	for key, val := range *rbt.RegexNoMatchMsgs {
		for idx := range val {
			flatRegexNoMatchMsgs = append(flatRegexNoMatchMsgs, key, idx, val[idx], errorPageDataID)
			flatLength++
		}
	}
	return flatLength, flatRegexNoMatchMsgs
}
