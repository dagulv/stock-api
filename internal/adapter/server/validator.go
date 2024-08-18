package server

import "github.com/go-playground/validator/v10"

type Validate struct {
	Validator *validator.Validate
}

func (v *Validate) Validate(i any) (err error) {
	return v.Validator.Struct(i)
}
