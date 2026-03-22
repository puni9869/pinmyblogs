package formbinding

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TestFieldErrors_Error_KnownPhrases(t *testing.T) {
	f := new(FieldErrors)
	tests := []struct {
		phrase string
		want   string
	}{
		{"required", "Cannot be empty"},
		{"email", "is not a valid email address."},
		{"password_not_match", "Password not match"},
		{"alpha_dash_dot", "Mandate 1 letter, 1 digit, 7-15 characters, with 6 non-digits"},
	}
	for _, tt := range tests {
		t.Run(tt.phrase, func(t *testing.T) {
			got := f.Error(tt.phrase)
			if got != tt.want {
				t.Errorf("Error(%q) = %q, want %q", tt.phrase, got, tt.want)
			}
		})
	}
}

func TestFieldErrors_Error_UnknownPhrase(t *testing.T) {
	f := new(FieldErrors)
	got := f.Error("nonexistent_phrase")
	if got != "nonexistent_phrase" {
		t.Errorf("Error(unknown) = %q, want passthrough of phrase", got)
	}
}

func TestFieldErrors_IsValid(t *testing.T) {
	f := new(FieldErrors)
	// Regex is ^[A-Za-z\d\W_]{6,15}$ — any chars, 6-15 length
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"valid 7 chars", "Abc1234", true},
		{"valid with special chars", "P@ss1word", true},
		{"too short 3 chars", "Ab1", false},
		{"all letters 8 chars", "Abcdefgh", true},
		{"all digits 8 chars", "12345678", true},
		{"empty", "", false},
		{"exactly 6 chars", "abcdef", true},
		{"exactly 15 chars", "abcdefghijklmno", true},
		{"too long 16 chars", "abcdefghijklmnop", false},
		{"5 chars too short", "abcde", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsValid(tt.password)
			if got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestFillContext(t *testing.T) {
	type TestForm struct {
		Name  string `form:"name"`
		Email string `form:"email"`
		Age   int    `form:"-"`
	}

	form := &TestForm{Name: "Alice", Email: "alice@test.com", Age: 30}
	data := make(gin.H)
	FillContext(form, data)

	if data["name"] != "Alice" {
		t.Errorf("data[name] = %v, want Alice", data["name"])
	}
	if data["Name_HasError"] != false {
		t.Errorf("data[Name_HasError] = %v, want false", data["Name_HasError"])
	}
	if data["Name_Error"] != "" {
		t.Errorf("data[Name_Error] = %v, want empty", data["Name_Error"])
	}

	if data["email"] != "alice@test.com" {
		t.Errorf("data[email] = %v, want alice@test.com", data["email"])
	}

	if _, ok := data["Age_HasError"]; ok {
		t.Error("Age field with form:\"-\" should be skipped")
	}
}

func TestFillContext_NoFormTag_UsesFieldName(t *testing.T) {
	type SimpleForm struct {
		Username string
	}

	form := &SimpleForm{Username: "bob"}
	data := make(gin.H)
	FillContext(form, data)

	if data["Username"] != "bob" {
		t.Errorf("data[Username] = %v, want bob", data["Username"])
	}
}

func TestErrorf_NilErrors(t *testing.T) {
	var valErrs validator.ValidationErrors
	result := Errorf(make(gin.H), valErrs)
	if result != nil {
		t.Errorf("Errorf with empty errors should return nil, got %v", result)
	}
}

func TestErrorf_WithErrors(t *testing.T) {
	type testForm struct {
		Email string `validate:"required,email"`
		Name  string `validate:"required"`
	}

	v := validator.New()
	err := v.Struct(testForm{})
	if err == nil {
		t.Fatal("expected validation errors")
	}

	var valErrs validator.ValidationErrors
	if !errors.As(err, &valErrs) {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	data := make(gin.H)
	result := Errorf(data, valErrs)

	if result == nil {
		t.Fatal("Errorf should return non-nil for validation errors")
	}

	if _, ok := result["Email_HasError"]; !ok {
		t.Error("expected Email_HasError in result")
	}
	if result["Email_HasError"] != true {
		t.Errorf("Email_HasError = %v, want true", result["Email_HasError"])
	}
	if _, ok := result["Email_Error"]; !ok {
		t.Error("expected Email_Error in result")
	}

	if _, ok := result["Name_HasError"]; !ok {
		t.Error("expected Name_HasError in result")
	}
	if result["Name_HasError"] != true {
		t.Errorf("Name_HasError = %v, want true", result["Name_HasError"])
	}
}

func TestErrorf_SingleError(t *testing.T) {
	type testForm struct {
		Pwd string `validate:"required"` //nolint:gosec
	}

	v := validator.New()
	err := v.Struct(testForm{})

	var valErrs validator.ValidationErrors
	if !errors.As(err, &valErrs) {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	data := make(gin.H)
	result := Errorf(data, valErrs)

	if result == nil {
		t.Fatal("Errorf should return non-nil")
	}
	if result["Pwd_HasError"] != true {
		t.Errorf("Pwd_HasError = %v, want true", result["Pwd_HasError"])
	}

	errMsg, ok := result["Pwd_Error"].(string)
	if !ok {
		t.Fatal("Pwd_Error should be a string")
	}
	if errMsg == "" {
		t.Error("Pwd_Error should not be empty")
	}
	if len(errMsg) < 3 {
		t.Errorf("Pwd_Error seems too short: %q", errMsg)
	}
}

func TestFieldErrors_Error_AllKnownPhrases(t *testing.T) {
	f := new(FieldErrors)
	knownPhrases := []string{
		"required", "alpha_dash", "alpha_dash_dot", "size",
		"min_size", "max_size", "email", "url", "include",
		"password_not_match", "username", "unknown",
	}
	for _, phrase := range knownPhrases {
		t.Run(phrase, func(t *testing.T) {
			got := f.Error(phrase)
			if got == phrase {
				t.Errorf("Error(%q) returned the phrase itself, expected a mapped message", phrase)
			}
			if got == "" {
				t.Errorf("Error(%q) returned empty string", phrase)
			}
		})
	}
}

func TestFillContext_MultipleFieldTypes(t *testing.T) {
	type MixedForm struct {
		Title   string `form:"title"`
		Count   int    `form:"count"`
		Active  bool   `form:"active"`
		Ignored string `form:"-"`
	}

	form := &MixedForm{Title: "Test", Count: 42, Active: true, Ignored: "skip"}
	data := make(gin.H)
	FillContext(form, data)

	if data["title"] != "Test" {
		t.Errorf("title = %v, want Test", data["title"])
	}
	if data["count"] != 42 {
		t.Errorf("count = %v, want 42", data["count"])
	}
	if data["active"] != true {
		t.Errorf("active = %v, want true", data["active"])
	}
	if _, ok := data["Ignored_HasError"]; ok {
		t.Error("Ignored field should be skipped")
	}
}

func TestFillContext_EmptyStruct(t *testing.T) {
	type EmptyForm struct{}

	form := &EmptyForm{}
	data := make(gin.H)
	FillContext(form, data)

	if len(data) != 0 {
		t.Errorf("FillContext on empty struct should produce empty data, got %d entries", len(data))
	}
}
