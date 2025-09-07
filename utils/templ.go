package utils

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

// RenderToString renders a templ template to a buffer and returns a string of that render
func RenderToString(component templ.Component, ctx context.Context) (string, error) {
	var buf bytes.Buffer
	err := component.Render(ctx, &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
