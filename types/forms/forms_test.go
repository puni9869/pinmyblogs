package forms

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func validate(obj any) error {
	return binding.Validator.ValidateStruct(obj) //nolint:wrapcheck
}

// --- SignUp ---

func TestSignUp_Valid(t *testing.T) {
	f := SignUp{Email: "user@example.com", Password: "secret", ConfirmPassword: "secret"}
	if err := validate(f); err != nil {
		t.Errorf("valid SignUp should pass: %v", err)
	}
}

func TestSignUp_MissingEmail(t *testing.T) {
	f := SignUp{Password: "secret", ConfirmPassword: "secret"}
	if err := validate(f); err == nil {
		t.Error("SignUp without email should fail")
	}
}

func TestSignUp_InvalidEmail(t *testing.T) {
	f := SignUp{Email: "not-an-email", Password: "secret", ConfirmPassword: "secret"}
	if err := validate(f); err == nil {
		t.Error("SignUp with invalid email should fail")
	}
}

func TestSignUp_MissingPassword(t *testing.T) {
	f := SignUp{Email: "user@example.com", ConfirmPassword: "secret"}
	if err := validate(f); err == nil {
		t.Error("SignUp without password should fail")
	}
}

func TestSignUp_MissingConfirmPassword(t *testing.T) {
	f := SignUp{Email: "user@example.com", Password: "secret"}
	if err := validate(f); err == nil {
		t.Error("SignUp without confirm_password should fail")
	}
}

// --- Reset ---

func TestReset_Valid(t *testing.T) {
	f := Reset{Email: "user@example.com"}
	if err := validate(f); err != nil {
		t.Errorf("valid Reset should pass: %v", err)
	}
}

func TestReset_MissingEmail(t *testing.T) {
	f := Reset{}
	if err := validate(f); err == nil {
		t.Error("Reset without email should fail")
	}
}

func TestReset_InvalidEmail(t *testing.T) {
	f := Reset{Email: "bad"}
	if err := validate(f); err == nil {
		t.Error("Reset with invalid email should fail")
	}
}

// --- ResetPassword ---

func TestResetPassword_Valid(t *testing.T) {
	f := ResetPassword{Hash: "abc", Password: "newpass", ConfirmPassword: "newpass"}
	if err := validate(f); err != nil {
		t.Errorf("valid ResetPassword should pass: %v", err)
	}
}

func TestResetPassword_MissingHash(t *testing.T) {
	f := ResetPassword{Password: "newpass", ConfirmPassword: "newpass"}
	if err := validate(f); err == nil {
		t.Error("ResetPassword without hash should fail")
	}
}

// --- WeblinkRequest ---

func TestWeblinkRequest_Valid(t *testing.T) {
	f := WeblinkRequest{Url: "https://go.dev", Tag: "golang"}
	if err := validate(f); err != nil {
		t.Errorf("valid WeblinkRequest should pass: %v", err)
	}
}

func TestWeblinkRequest_MissingUrl(t *testing.T) {
	f := WeblinkRequest{Tag: "test"}
	if err := validate(f); err == nil {
		t.Error("WeblinkRequest without url should fail")
	}
}

func TestWeblinkRequest_MissingTag(t *testing.T) {
	f := WeblinkRequest{Url: "https://go.dev"}
	if err := validate(f); err == nil {
		t.Error("WeblinkRequest without tag should fail")
	}
}

// --- BulkAction ---

func TestBulkAction_Valid(t *testing.T) {
	f := BulkAction{IDs: []string{"1", "2"}, Action: "bulk-delete"}
	if err := validate(f); err != nil {
		t.Errorf("valid BulkAction should pass: %v", err)
	}
}

func TestBulkAction_MissingIDs(t *testing.T) {
	f := BulkAction{Action: "bulk-delete"}
	if err := validate(f); err == nil {
		t.Error("BulkAction without IDs should fail")
	}
}

func TestBulkAction_MissingAction(t *testing.T) {
	f := BulkAction{IDs: []string{"1"}}
	if err := validate(f); err == nil {
		t.Error("BulkAction without action should fail")
	}
}

// --- Prefs ---

func TestPrefs_Valid(t *testing.T) {
	f := Prefs{Action: "sideNav", Value: "hide"}
	if err := validate(f); err != nil {
		t.Errorf("valid Prefs should pass: %v", err)
	}
}

func TestPrefs_ActionTooShort(t *testing.T) {
	f := Prefs{Action: "a", Value: "hide"}
	if err := validate(f); err == nil {
		t.Error("Prefs with 1-char action should fail (min=2)")
	}
}

func TestPrefs_ValueTooLong(t *testing.T) {
	f := Prefs{Action: "sideNav", Value: "this-value-is-way-too-long-for-the-max"}
	if err := validate(f); err == nil {
		t.Error("Prefs with long value should fail (max=20)")
	}
}

// --- JoinWaitList ---

func TestJoinWaitList_Valid(t *testing.T) {
	f := JoinWaitList{Email: "user@example.com"}
	if err := validate(f); err != nil {
		t.Errorf("valid JoinWaitList should pass: %v", err)
	}
}

func TestJoinWaitList_MissingEmail(t *testing.T) {
	f := JoinWaitList{}
	if err := validate(f); err == nil {
		t.Error("JoinWaitList without email should fail")
	}
}

// --- AccountEnable ---

func TestAccountEnable_Valid(t *testing.T) {
	f := AccountEnable{Hash: "550e8400-e29b-41d4-a716-446655440000"}
	if err := validate(f); err != nil {
		t.Errorf("valid AccountEnable should pass: %v", err)
	}
}

func TestAccountEnable_MissingHash(t *testing.T) {
	f := AccountEnable{}
	if err := validate(f); err == nil {
		t.Error("AccountEnable without hash should fail")
	}
}

func TestAccountEnable_InvalidUUID(t *testing.T) {
	f := AccountEnable{Hash: "not-a-uuid"}
	if err := validate(f); err == nil {
		t.Error("AccountEnable with non-uuid4 hash should fail")
	}
}

// --- Binding tag check (ensure gin binding tags are parseable) ---

func TestAllForms_BindingTagsParseable(_ *testing.T) {
	// This ensures the binding tags don't have typos that would silently fail
	forms := []any{
		SignUp{},
		Reset{},
		ResetPassword{},
		WeblinkRequest{},
		BulkAction{},
		Prefs{},
		JoinWaitList{},
		AccountEnable{},
	}
	_ = binding.Form // ensure gin binding is importable
	for _, f := range forms {
		_ = validate(f) // just making sure validate doesn't panic
	}
}
