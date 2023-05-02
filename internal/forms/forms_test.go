package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	has := form.Has("whatever", r)
	if has {
		t.Error("Form shows that it has a field when it does NOT !")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	has = form.Has("a", r)
	if has {
		t.Error("Form shows that it has a field when it does NOT !!")
	}
}

func TestForm_Minlength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.MinLength("xfield", 10)
	if form.Valid() {
		t.Error("form shows minlength for a non-existent field")
	}
	isErr := form.Errors.Get("xfield")
	if isErr == "" {
		t.Error("Shoudl have an error, but we did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some value")
	form = New(postedValues)
	form.MinLength("some_field", 1000)
	if !form.Valid() {
		t.Error("form shows minlength of 1000 met when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abcd1234")
	form = New(postedValues)
	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("form shows minlength of 1 is not met when data is way longer")
	}

	isErr = form.Errors.Get("another_field")
	if isErr != "" {
		t.Error("Shoudl NOT have an error, but we did get one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows a VALID EMAIL for a non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "bryan@kouhadi.com")
	form = New(postedValues)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Got an invalid email error when we should not")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "bryan?kouhadi")
	form = New(postedValues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got Valid for an invalid email")
	}
}
