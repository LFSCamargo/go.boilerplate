package reactemail_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go.boilerplate/src/mail/reactemail"
)

func findEmailsDir(t *testing.T) string {
	t.Helper()
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "emails", "package.json")); err == nil {
			return filepath.Join(root, "emails")
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Skip("emails workspace not found")
		}
		root = parent
	}
}

func TestNewRenderer_MissingWorkspace(t *testing.T) {
	_, err := reactemail.NewRenderer("/tmp/does-not-exist-emails", "")
	if err == nil {
		t.Fatal("expected error for missing emails workspace")
	}
}

func TestRenderer_Render_UnknownTemplate(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not installed")
	}

	emailsDir := findEmailsDir(t)
	if _, err := os.Stat(filepath.Join(emailsDir, "node_modules")); err != nil {
		t.Skip("run npm install in emails/")
	}

	renderer, err := reactemail.NewRenderer(emailsDir, "http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	_, err = renderer.Render("does-not-exist", map[string]any{})
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
}

func TestRenderer_Render_VerifyEmailIncludesGoBoilerplateAssets(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not installed")
	}

	emailsDir := findEmailsDir(t)
	if _, err := os.Stat(filepath.Join(emailsDir, "node_modules")); err != nil {
		t.Skip("run npm install in emails/")
	}

	renderer, err := reactemail.NewRenderer(emailsDir, "http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	html, err := renderer.Render("verify_email", map[string]any{
		"name":          "Ada",
		"code":          "123456",
		"companyName":   "Go Boilerplate",
		"expiryMinutes": 10,
		"verifyUrl":     "http://localhost:8080/api/v1/auth/verify-email?token=abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{"Go Boilerplate", "static/go_boilerplate/hero_image.png", "static/go_boilerplate/logo_mark.svg", "123456"} {
		if !bytesContains(html, needle) {
			t.Errorf("rendered html missing %q", needle)
		}
	}
}

func bytesContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || (len(s) > 0 && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()))
}
