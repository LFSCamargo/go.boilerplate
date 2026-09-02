package repositories

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const maxAvatarBytes = 2 << 20

var allowedAvatarExt = map[string]string{
	".jpg":  ".jpg",
	".jpeg": ".jpg",
	".png":  ".png",
	".webp": ".webp",
	".gif":  ".gif",
}

func (r *UserRepository) SaveAvatar(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}
	if r.uploadsDir == "" {
		return "", fmt.Errorf("uploads directory is not configured")
	}
	if file.Size > maxAvatarBytes {
		return "", fmt.Errorf("picture too large")
	}

	ext := allowedAvatarExt[strings.ToLower(filepath.Ext(file.Filename))]
	if ext == "" {
		return "", fmt.Errorf("picture must be jpeg, png, webp, or gif")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dir := filepath.Join(r.uploadsDir, "avatars")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name := uuid.NewString() + ext
	dst, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, io.LimitReader(src, maxAvatarBytes+1)); err != nil {
		return "", err
	}

	base := strings.TrimRight(r.publicURL, "/")
	return base + "/uploads/avatars/" + name, nil
}
