package eval

import (
	"bytes"
	"text/template"
)

// Context holds the variables available for evaluation.
type Context map[string]interface{}

// Evaluate interpolates {{vars}} in a string using the provided context.
func Evaluate(input string, ctx Context) (string, error) {
	t, err := template.New("eval").Parse(input)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}

	return buf.String(), nil
}
