package datatypes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type DataType interface {
	Name() string
	Validate(value interface{}) error
	Convert(value interface{}) (interface{}, error)
	MarshalJSON() ([]byte, error)
}

type IntegerType struct{}
type TextType struct{}
type BooleanType struct{}
type DateTimeType struct{}

func (t *IntegerType) Name() string  { return "INTEGER" }
func (t *TextType) Name() string     { return "TEXT" }
func (t *BooleanType) Name() string  { return "BOOLEAN" }
func (t *DateTimeType) Name() string { return "DATETIME" }

func (t *IntegerType) Validate(value interface{}) error {
	switch v := value.(type) {
	case int, int32, int64:
		return nil
	case float64:
		// Check if it's a whole number
		if v == float64(int(v)) {
			return nil
		}
		return fmt.Errorf("float value must be a whole number: %v", value)
	case string:
		_, err := strconv.ParseInt(v, 10, 64)
		return err
	default:
		return fmt.Errorf("invalid integer value: %v", value)
	}
}

func (t *TextType) Validate(value interface{}) error {
	switch value.(type) {
	case string:
		return nil
	default:
		return fmt.Errorf("invalid text value: %v", value)
	}
}

func (t *BooleanType) Validate(value interface{}) error {
	switch v := value.(type) {
	case bool:
		return nil
	case string:
		_, err := strconv.ParseBool(v)
		return err
	default:
		return fmt.Errorf("invalid boolean value: %v", value)
	}
}

func (t *DateTimeType) Validate(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		return nil
	case string:
		_, err := time.Parse(time.RFC3339, v)
		return err
	default:
		return fmt.Errorf("invalid datetime value: %v", value)
	}
}

// Convert functions
func (t *IntegerType) Convert(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case int, int32, int64:
		return v, nil
	case float64:
		// Check if it's a whole number
		if v == float64(int(v)) {
			return int(v), nil
		}
		return nil, fmt.Errorf("float value must be a whole number: %v", value)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		return int(i), nil
	default:
		return nil, fmt.Errorf("cannot convert to integer: %v", value)
	}
}

func (t *TextType) Convert(value interface{}) (interface{}, error) {
	return fmt.Sprintf("%v", value), nil
}

func (t *BooleanType) Convert(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return nil, fmt.Errorf("cannot convert to boolean: %v", value)
	}
}

func (t *DateTimeType) Convert(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return nil, fmt.Errorf("cannot convert to datetime: %v", value)
	}
}

func (t *IntegerType) MarshalJSON() ([]byte, error) {
	return json.Marshal("INTEGER")
}

func (t *TextType) MarshalJSON() ([]byte, error) {
	return json.Marshal("TEXT")
}

func (t *BooleanType) MarshalJSON() ([]byte, error) {
	return json.Marshal("BOOLEAN")
}

func (t *DateTimeType) MarshalJSON() ([]byte, error) {
	return json.Marshal("DATETIME")
}

// GetType returns the appropriate DataType for a type name
func GetType(typeName string) (DataType, error) {
	switch typeName {
	case "INTEGER", "INT":
		return &IntegerType{}, nil
	case "TEXT", "VARCHAR", "STRING":
		return &TextType{}, nil
	case "BOOLEAN", "BOOL":
		return &BooleanType{}, nil
	case "DATETIME", "TIMESTAMP":
		return &DateTimeType{}, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", typeName)
	}
}
