package argument

import (
	"errors"
)

type BooleanValue struct {
	Type  BooleanType `json:"type" yaml:"type"`
	True  string      `json:"true" yaml:"true"`
	False string      `json:"false" yaml:"false"`
}

var (
	BOOLEANVALUE_BOOL          = BooleanValue{BOOLEANTYPE_BOOL, "true", "false"}
	BOOLEANVALUE_BOOLUPPER     = BooleanValue{BOOLEANTYPE_BOOLUPPER, "True", "False"}
	BOOLEANVALUE_BOOLUPPERFULL = BooleanValue{BOOLEANTYPE_BOOLUPPERFULL, "TRUE", "FALSE"}
	BOOLEANVALUE_INTEGER       = BooleanValue{BOOLEANTYPE_INTEGER, "1", "0"}
)

var _ ArgumentWithValue[bool] = (*Boolean)(nil)

type Boolean struct {
	Base    `yaml:",inline"`
	Default bool          `json:"default" yaml:"default"`
	Value   bool          `json:"value" yaml:"value"`
	Values  *BooleanValue `json:"values" yaml:"values"`
}

func (argument *Boolean) GetDefault() bool {
	return argument.Default
}

func (argument *Boolean) GetValue() bool {
	return argument.Value
}

func (argument *Boolean) Parse(sep Seperator) (arg string, err error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	arg = argument.GetArgument() + *argument.GetSeperator(sep).ArgumentValue
	if argument.GetValue() {
		arg = arg + argument.Values.True
	} else {
		arg = arg + argument.Values.False
	}
	return
}

func (argument *Boolean) NormalizeState() {
	if argument.Values != nil {
		switch argument.Values.Type {
		case BOOLEANTYPE_BOOL:
			argument.Values = &BOOLEANVALUE_BOOL
		case BOOLEANTYPE_BOOLUPPER:
			argument.Values = &BOOLEANVALUE_BOOLUPPER
		case BOOLEANTYPE_BOOLUPPERFULL:
			argument.Values = &BOOLEANVALUE_BOOLUPPERFULL
		case BOOLEANTYPE_INTEGER:
			argument.Values = &BOOLEANVALUE_INTEGER
		}
	} else {
		argument.Values = &BOOLEANVALUE_BOOL
	}
	argument.Value = argument.Default
}

func (argument *Boolean) ValidateLazy() (err error) {
	if argument.Type != TYPE_BOOLEAN {
		err = errors.Join(err, errors.New("Boolean is not of type BOOLEAN"))
	}
	return
}

func (argument *Boolean) Reset() {
	argument.Value = argument.Default
}
