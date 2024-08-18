package chat

import (
	"encoding/json"
	"errors"
	"fmt"
)

func TypeFromString(s string) (Type, error) {
	switch s {
	case TYPE_TEXT.slug:
		return TYPE_TEXT, nil
	case TYPE_FILE.slug:
		return TYPE_FILE, nil
	}
	return TYPE_UNDEFINED, errors.New("unknown type: " + s)
}

var (
	TYPE_UNDEFINED = Type{""}
	TYPE_TEXT      = Type{"text"}
	TYPE_FILE      = Type{"file"}
)

type Type struct {
	slug string
}

func (t Type) String() string {
	return t.slug
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
