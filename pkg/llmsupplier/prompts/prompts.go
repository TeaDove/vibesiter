package prompts

import (
	"bytes"
	"embed"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/cockroachdb/errors"
)

//go:embed *
var prompts embed.FS

type Prompts struct {
	Templates *template.Template
}

func NewPrompts() (*Prompts, error) {
	templates, err := template.ParseFS(prompts, "*.gohtml")
	if err != nil {
		return nil, errors.Wrap(err, "parse templates")
	}

	return &Prompts{templates}, nil
}

func (r *Prompts) Render(name string, data map[string]any) (string, error) {
	if !strings.HasSuffix(name, ".gohtml") {
		name += ".gohtml"
	}

	data["jsonMarshal"] = func(v any) (string, error) {
		marshalled, err := json.Marshal(v)

		return string(marshalled), err
	}

	var buf bytes.Buffer

	err := r.Templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", errors.Wrap(err, "render templates")
	}

	return buf.String(), nil
}
