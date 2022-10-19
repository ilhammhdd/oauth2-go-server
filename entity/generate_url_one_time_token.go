package entity

import (
	"github.com/ilhammhdd/go-toolkit/errorkit"
)

type GenerateURLOneTimeToken struct {
	ClientID    string
	AccessToken string
	Path        string
	Query       string
	Fragment    string
	ErrDescGen  errorkit.ErrDescGenerator
}

func (guott *GenerateURLOneTimeToken) GenerateRedirectLocationURL() string {

	return ""
}
