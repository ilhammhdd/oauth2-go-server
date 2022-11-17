package usecase

import "github.com/ilhammhdd/go-toolkit/errorkit"

type LoginUserByEmail struct {
	errDescGen errorkit.ErrDescGenerator
}

func (lue LoginUserByEmail) VerifyPassword() {

}
