package eval

import (
	"bytes"
	"regexp"
	"text/template"
)

// Context holds the variables available for evaluation.
type Context map[string]interface{}

var varRegex = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\s*\}\}`)

// Evaluate interpolates {{vars}} in a string using the provided context.
func Evaluate(input string, ctx Context) (string, error) {
	// Convert {{var}} to {{.var}} for text/template compatibility
	transformed := varRegex.ReplaceAllString(input, "{{.$1}}")

	t, err := template.New("eval").Parse(transformed)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		// If template execution fails, return original input to avoid breaking things
		return transformed, err
	}

	return buf.String(), nil
}
