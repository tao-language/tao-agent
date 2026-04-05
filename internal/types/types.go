package types

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Kind string

const (
	KindNull    Kind = "Null"
	KindBoolean Kind = "Boolean"
	KindNumber  Kind = "Number"
	KindString  Kind = "String"
	KindLiteral Kind = "Literal"
	KindList    Kind = "List"
	KindTuple   Kind = "Tuple"
	KindRecord  Kind = "Record"
	KindUnion   Kind = "Union"
	KindResult  Kind = "Result"
)

// Definition describes a Tao type and its properties.
type Definition struct {
	Type         Kind                   `yaml:"type"`
	Description  string                 `yaml:"description,omitempty"`
	Default      interface{}            `yaml:"default,omitempty"`
	Value        interface{}            `yaml:"value,omitempty"`        // For Literal
	Items        *Definition            `yaml:"items,omitempty"`        // For List
	TupleItems   []*Definition          `yaml:"tuple_items,omitempty"`  // For Tuple
	Fields       map[string]*Definition `yaml:"fields,omitempty"`       // For Record
	Alternatives []*Definition          `yaml:"alternatives,omitempty"` // For Union
	Ok           *Definition            `yaml:"ok,omitempty"`           // For Result
	Error        *Definition            `yaml:"error,omitempty"`        // For Result
}

// UnmarshalYAML implements custom YAML unmarshaling for sugar syntax.
func (d *Definition) UnmarshalYAML(value *yaml.Node) error {
	// Handle sugar syntax like "String"
	if value.Kind == yaml.ScalarNode {
		d.Type = Kind(value.Value)
		return nil
	}

	// For more complex cases, use a temporary struct to avoid recursion
	type alias Definition
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*d = Definition(aux)

	// Defaults for Result
	if d.Type == KindResult {
		if d.Ok == nil {
			d.Ok = &Definition{Type: KindString}
		}
		if d.Error == nil {
			d.Error = &Definition{Type: KindString}
		}
	}

	return nil
}

// Validate checks if a value matches the type definition.
func (d *Definition) Validate(v interface{}) error {
	if v == nil {
		if d.Type == KindNull {
			return nil
		}
		return fmt.Errorf("expected %s, got null", d.Type)
	}

	switch d.Type {
	case KindBoolean:
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("expected Boolean, got %T", v)
		}
	case KindNumber:
		switch v.(type) {
		case float64, int, int64:
			return nil
		default:
			return fmt.Errorf("expected Number, got %T", v)
		}
	case KindString:
		if _, ok := v.(string); !ok {
			return fmt.Errorf("expected String, got %T", v)
		}
	case KindLiteral:
		if v != d.Value {
			return fmt.Errorf("expected Literal %v, got %v", d.Value, v)
		}
	case KindList:
		list, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("expected List, got %T", v)
		}
		if d.Items != nil {
			for i, item := range list {
				if err := d.Items.Validate(item); err != nil {
					return fmt.Errorf("invalid list item at index %d: %w", i, err)
				}
			}
		}
	case KindRecord:
		record, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected Record, got %T", v)
		}
		for name, fieldDef := range d.Fields {
			val, exists := record[name]
			if !exists {
				if fieldDef.Default == nil {
					return fmt.Errorf("missing required field: %s", name)
				}
				continue
			}
			if err := fieldDef.Validate(val); err != nil {
				return fmt.Errorf("invalid field %s: %w", name, err)
			}
		}
	}
	return nil
}
