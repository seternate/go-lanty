package argument

import (
	"errors"
	"slices"
)

type EnumItem struct {
	Name    string `json:"name" yaml:"name"`
	Value   string `json:"value" yaml:"value"`
	Default bool   `json:"default,omitempty" yaml:"default,omitempty"`
}

var _ ArgumentWithValue[string] = (*Enum)(nil)

type Enum struct {
	Base  `yaml:",inline"`
	Value string     `json:"value" yaml:"value"`
	Items []EnumItem `json:"items" yaml:"items"`
}

func (argument *Enum) GetDefault() string {
	for _, item := range argument.Items {
		if item.Default {
			return item.Value
		}
	}
	return ""
}

func (argument *Enum) GetValue() string {
	return argument.Value
}

func (argument *Enum) Parse(sep Seperator) (arg string, err error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	arg = argument.GetArgument() + *argument.GetSeperator(sep).ArgumentValue + argument.GetValue()
	return
}

func (argument *Enum) NormalizeState() {
	argument.Value = argument.GetDefault()
}

func (argument *Enum) ValidateLazy() (err error) {
	if argument.Type != TYPE_ENUM {
		err = errors.Join(err, errors.New("Enum is not of type ENUM"))
	}
	if len(argument.GetDefault()) == 0 {
		err = errors.Join(err, errors.New("default value is empty"))
	}
	if len(argument.Items) == 0 {
		err = errors.Join(err, errors.New("no items specified in Enum"))
	}
	return
}

func (left *Enum) Equal(right *Enum) bool {
	return left.Base == right.Base &&
		left.Value == right.Value &&
		slices.Equal(left.Items, right.Items)
}

func (argument *Enum) Reset() {
	argument.Value = argument.GetDefault()
}
