package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func TestGetForm_ReturnsNilWhenNotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	got := GetForm(c)
	if got != nil {
		t.Errorf("GetForm() = %v, want nil", got)
	}
}

func TestGetForm_ReturnsValueWhenSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expected := "test-form-value"
	c.Set(formKey, expected)

	got := GetForm(c)
	if got != expected {
		t.Errorf("GetForm() = %v, want %v", got, expected)
	}
}

func TestGetForm_ReturnsStructWhenSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	type MyForm struct {
		Name string
	}
	expected := &MyForm{Name: "Alice"}
	c.Set(formKey, expected)

	got := GetForm(c)
	form, ok := got.(*MyForm)
	if !ok {
		t.Fatalf("GetForm() type = %T, want *MyForm", got)
	}
	if form.Name != "Alice" {
		t.Errorf("form.Name = %q, want Alice", form.Name)
	}
}

func TestGetContext_ReturnsNilWhenNotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	got := GetContext(c)
	if got != nil {
		t.Errorf("GetContext() = %v, want nil", got)
	}
}

func TestGetContext_ReturnsMapWhenSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expected := gin.H{"HasError": false, "Email_Error": ""}
	c.Set(contextKey, map[string]any(expected))

	got := GetContext(c)
	if got == nil {
		t.Fatal("GetContext() = nil, want non-nil")
	}
	if got["HasError"] != false {
		t.Errorf("HasError = %v, want false", got["HasError"])
	}
}

func TestBind_FormContentType(t *testing.T) {
	type TestForm struct {
		Name string `form:"name" binding:"required"`
	}

	r := gin.New()
	r.POST("/test", Bind(TestForm{}), func(c *gin.Context) {
		form := GetForm(c)
		if form == nil {
			c.String(http.StatusInternalServerError, "no form")
			return
		}
		f := form.(*TestForm)
		c.String(http.StatusOK, f.Name)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("name=Alice"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if w.Body.String() != "Alice" {
		t.Errorf("body = %q, want Alice", w.Body.String())
	}
}

func TestBind_JSONContentType(t *testing.T) {
	type TestForm struct {
		Name string `json:"name" binding:"required"`
	}

	r := gin.New()
	r.POST("/test", Bind(TestForm{}), func(c *gin.Context) {
		form := GetForm(c)
		f := form.(*TestForm)
		c.String(http.StatusOK, f.Name)
	})

	w := httptest.NewRecorder()
	body := `{"name":"Bob"}`
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", binding.MIMEJSON)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if w.Body.String() != "Bob" {
		t.Errorf("body = %q, want Bob", w.Body.String())
	}
}

func TestBind_SetsContextWithHasErrorFalse(t *testing.T) {
	type TestForm struct {
		Value string `form:"value" binding:"required"`
	}

	r := gin.New()
	r.POST("/test", Bind(TestForm{}), func(c *gin.Context) {
		ctx := GetContext(c)
		if ctx == nil {
			c.String(http.StatusInternalServerError, "no context")
			return
		}
		if ctx["HasError"] != false {
			c.String(http.StatusInternalServerError, "HasError not false")
			return
		}
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("value=test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	if w.Body.String() != "ok" {
		t.Errorf("body = %q, want ok", w.Body.String())
	}
}

func TestBind_ValidationError_SetsErrorContext(t *testing.T) {
	type TestForm struct {
		Email string `form:"email" binding:"required,email"`
	}

	r := gin.New()
	r.POST("/test", Bind(TestForm{}), func(c *gin.Context) {
		ctx := GetContext(c)
		if ctx == nil {
			c.String(http.StatusInternalServerError, "no context")
			return
		}
		// When validation fails, Errorf replaces the context
		if _, ok := ctx["Email_HasError"]; ok {
			c.String(http.StatusOK, "has_error")
			return
		}
		c.String(http.StatusOK, "no_error")
	})

	// Send empty form — email is required so validation should fail
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	if w.Body.String() != "has_error" {
		t.Errorf("body = %q, want has_error (validation should fail)", w.Body.String())
	}
}

func TestBind_FillsFormFieldsInContext(t *testing.T) {
	type TestForm struct {
		Name string `form:"name" binding:"required"`
		Age  string `form:"age"`
	}

	r := gin.New()
	r.POST("/test", Bind(TestForm{}), func(c *gin.Context) {
		ctx := GetContext(c)
		if ctx["name"] != "Charlie" {
			c.String(http.StatusInternalServerError, "name wrong")
			return
		}
		if ctx["age"] != "25" {
			c.String(http.StatusInternalServerError, "age wrong")
			return
		}
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("name=Charlie&age=25"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	if w.Body.String() != "ok" {
		t.Errorf("body = %q, want ok", w.Body.String())
	}
}
