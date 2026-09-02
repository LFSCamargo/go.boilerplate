package reactemail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Renderer struct {
	emailsDir     string
	assetsBaseURL string
	nodeBin       string
	npmBin        string
}

func NewRenderer(emailsDir string, assetsBaseURL string) (*Renderer, error) {
	if emailsDir == "" {
		emailsDir = "emails"
	}

	absDir, err := filepath.Abs(emailsDir)
	if err != nil {
		return nil, fmt.Errorf("resolve emails dir: %w", err)
	}

	if _, err := os.Stat(filepath.Join(absDir, "package.json")); err != nil {
		return nil, fmt.Errorf("emails workspace not found at %s: %w", absDir, err)
	}

	nodeBin, err := exec.LookPath("node")
	if err != nil {
		return nil, fmt.Errorf("node not found in PATH: %w", err)
	}

	npmBin, err := exec.LookPath("npm")
	if err != nil {
		return nil, fmt.Errorf("npm not found in PATH: %w", err)
	}

	return &Renderer{
		emailsDir:     absDir,
		assetsBaseURL: assetsBaseURL,
		nodeBin:       nodeBin,
		npmBin:        npmBin,
	}, nil
}

func (r *Renderer) Render(template string, props map[string]any) (string, error) {
	if template == "" {
		return "", fmt.Errorf("template name is required")
	}

	propsJSON, err := json.Marshal(props)
	if err != nil {
		return "", fmt.Errorf("marshal props: %w", err)
	}

	renderScript := filepath.Join(r.emailsDir, "node_modules", ".bin", "tsx")
	if _, err := os.Stat(renderScript); err != nil {
		renderScript = ""
	}

	var cmd *exec.Cmd
	if renderScript != "" {
		cmd = exec.Command(renderScript, "render.ts", template, string(propsJSON))
	} else {
		cmd = exec.Command(r.npmBin, "run", "render", "--", template, string(propsJSON))
	}
	cmd.Dir = r.emailsDir
	cmd.Env = os.Environ()
	if r.assetsBaseURL != "" {
		cmd.Env = append(cmd.Env, "EMAIL_ASSETS_BASE_URL="+r.assetsBaseURL)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("react-email render: %w: %s", err, stderr.String())
	}

	html := stdout.String()
	if html == "" {
		return "", fmt.Errorf("react-email render returned empty html")
	}

	return html, nil
}
