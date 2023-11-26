package argument

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

func BooleanTypeFromString(s string) (BooleanType, error) {
	switch s {
	case BOOLEANTYPE_BOOL.slug:
		return BOOLEANTYPE_BOOL, nil
	case BOOLEANTYPE_BOOLUPPER.slug:
		return BOOLEANTYPE_BOOLUPPER, nil
	case BOOLEANTYPE_BOOLUPPERFULL.slug:
		return BOOLEANTYPE_BOOLUPPERFULL, nil
	case BOOLEANTYPE_INTEGER.slug:
		return BOOLEANTYPE_INTEGER, nil
	case BOOLEANTYPE_CUSTOM.slug:
		return BOOLEANTYPE_CUSTOM, nil
	}
	return BOOLEANTYPE_UNDEFINED, errors.New("unknown booleantype: " + s)
}

var (
	BOOLEANTYPE_UNDEFINED     = BooleanType{""}
	BOOLEANTYPE_BOOL          = BooleanType{"bool"}
	BOOLEANTYPE_BOOLUPPER     = BooleanType{"boolupper"}
	BOOLEANTYPE_BOOLUPPERFULL = BooleanType{"boolupperfull"}
	BOOLEANTYPE_INTEGER       = BooleanType{"integer"}
	BOOLEANTYPE_CUSTOM        = BooleanType{"custom"}
)

type BooleanType struct {
	slug string
}

func (t BooleanType) String() string {
	return t.slug
}

func (t *BooleanType) UnmarshalYAML(value *yaml.Node) (err error) {
	var field string
	err = value.Decode(&field)
	if err != nil {
		return
	}
	*t, err = BooleanTypeFromString(field)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.New("error unmarshal yaml at "+strconv.Itoa(value.Line)+":"+strconv.Itoa(value.Column)), err)
	}
	return
}

func (t *BooleanType) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var field string
	err = json.Unmarshal(data, &field)
	if err != nil {
		return
	}
	*t, err = BooleanTypeFromString(field)
	if err != nil {
		err = fmt.Errorf("%w: %w", errors.New("error unmarshal json"), err)
	}
	return
}

func (t BooleanType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
