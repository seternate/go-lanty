package argument

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

func TypeFromString(s string) (Type, error) {
	switch s {
	case TYPE_BOOLEAN.slug:
		return TYPE_BOOLEAN, nil
	case TYPE_STRING.slug:
		return TYPE_STRING, nil
	case TYPE_INTEGER.slug:
		return TYPE_INTEGER, nil
	case TYPE_FLOAT.slug:
		return TYPE_FLOAT, nil
	case TYPE_ENUM.slug:
		return TYPE_ENUM, nil
	case TYPE_BASE.slug:
		return TYPE_BASE, nil
	case TYPE_CONNECT.slug:
		return TYPE_CONNECT, nil
	}
	return TYPE_UNDEFINED, errors.New("unknown type: " + s)
}

var (
	TYPE_UNDEFINED = Type{""}
	TYPE_BOOLEAN   = Type{"boolean"}
	TYPE_STRING    = Type{"string"}
	TYPE_INTEGER   = Type{"integer"}
	TYPE_FLOAT     = Type{"float"}
	TYPE_ENUM      = Type{"enum"}
	TYPE_BASE      = Type{"base"}
	TYPE_CONNECT   = Type{"connect"}
)

type Type struct {
	slug string
}

func (t Type) String() string {
	return t.slug
}

func (t *Type) UnmarshalYAML(value *yaml.Node) (err error) {
	var field string
	err = value.Decode(&field)
	if err != nil {
		return
	}
	*t, err = TypeFromString(field)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.New("error unmarshal yaml at "+strconv.Itoa(value.Line)+":"+strconv.Itoa(value.Column)), err)
	}
	return
}

func (t *Type) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var field string
	err = json.Unmarshal(data, &field)
	if err != nil {
		return
	}
	*t, err = TypeFromString(field)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.New("error unmarshal json"), err)
	}
	return
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
