package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form create a custom Form struct, and embeds url.Values Object
type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0 // shorthand of if length is 0 return true otherwise return false
}

// New Initializes a Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if a form's field has a value(not empty)
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be Blank !")
		}
	}
}

// minLength checks the minimum length of the field's inserted value
func (f *Form) MinLength(field string, minLength int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < minLength {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d Characters long", minLength))
		return false
	}
	return true
}

// IsEmail checks the validity of email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address !")
	}
}
