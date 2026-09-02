package repositories_test

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.boilerplate/src/modules/auth/repositories"
)

func TestUserRepository_SaveAvatar(t *testing.T) {
	dir := t.TempDir()
	repo := repositories.NewUserRepository(nil, dir, "http://localhost:8080")
	file := formFile(t, "avatar.png", []byte{0x89, 0x50, 0x4e, 0x47})

	url, err := repo.SaveAvatar(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "http://localhost:8080/uploads/avatars/") || !strings.HasSuffix(url, ".png") {
		t.Fatalf("unexpected url %s", url)
	}

	entries, err := os.ReadDir(filepath.Join(dir, "avatars"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("expected saved file, err=%v entries=%d", err, len(entries))
	}
}

func TestUserRepository_SaveAvatar_RejectsBadType(t *testing.T) {
	repo := repositories.NewUserRepository(nil, t.TempDir(), "http://localhost:8080")
	file := formFile(t, "notes.txt", []byte("hello"))
	if _, err := repo.SaveAvatar(file); err == nil {
		t.Fatal("expected invalid type")
	}
}

func formFile(t *testing.T, filename string, data []byte) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("picture", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })

	files := form.File["picture"]
	if len(files) == 0 {
		t.Fatal("missing picture part")
	}
	return files[0]
}
