package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatePool = sync.Pool{
		New: func() interface{} {
			v := validator.New(validator.WithRequiredStructEnabled())
			return v
		},
	}
)

func Get() *validator.Validate {
	return validatePool.Get().(*validator.Validate)
}

func Put(v *validator.Validate) {
	validatePool.Put(v)
}

func Validate(s interface{}) error {
	v := Get()
	defer Put(v)
	return v.Struct(s)
}
