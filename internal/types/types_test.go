package types

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestUnmarshalSugar(t *testing.T) {
	data := `String`
	var def Definition
	if err := yaml.Unmarshal([]byte(data), &def); err != nil {
		t.Fatalf("failed to unmarshal sugar syntax: %v", err)
	}
	if def.Type != KindString {
		t.Errorf("expected KindString, got %s", def.Type)
	}
}

func TestValidateString(t *testing.T) {
	def := Definition{Type: KindString}
	if err := def.Validate("hello"); err != nil {
		t.Errorf("expected validation to pass: %v", err)
	}
	if err := def.Validate(123); err == nil {
		t.Errorf("expected validation to fail for non-string")
	}
}

func TestValidateNumber(t *testing.T) {
	def := Definition{Type: KindNumber}
	if err := def.Validate(42); err != nil {
		t.Errorf("expected validation to pass for int: %v", err)
	}
	if err := def.Validate(3.14); err != nil {
		t.Errorf("expected validation to pass for float: %v", err)
	}
	if err := def.Validate("42"); err == nil {
		t.Errorf("expected validation to fail for string")
	}
}

func TestValidateList(t *testing.T) {
	def := Definition{
		Type:  KindList,
		Items: &Definition{Type: KindString},
	}
	if err := def.Validate([]interface{}{"a", "b"}); err != nil {
		t.Errorf("expected list validation to pass: %v", err)
	}
	if err := def.Validate([]interface{}{"a", 1}); err == nil {
		t.Errorf("expected list validation to fail for mixed types")
	}
}
