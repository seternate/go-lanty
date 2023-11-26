package argument

import (
	"errors"
	"strconv"
)

var _ ArgumentWithNumberValue[int] = (*Integer)(nil)

type Integer struct {
	Base     `yaml:",inline"`
	Default  int `json:"default" yaml:"default"`
	Value    int `json:"value" yaml:"value"`
	MinValue int `json:"minvalue" yaml:"minvalue"`
	MaxValue int `json:"maxvalue" yaml:"maxvalue"`
}

func (argument *Integer) GetDefault() int {
	return argument.Default
}

func (argument *Integer) GetValue() int {
	return argument.Value
}

func (argument *Integer) GetMinValue() int {
	return argument.MinValue
}

func (argument *Integer) GetMaxValue() int {
	return argument.MaxValue
}

func (argument *Integer) Parse(sep Seperator) (arg string, err error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	arg = argument.GetArgument() + *argument.GetSeperator(sep).ArgumentValue + strconv.Itoa(argument.GetValue())
	return
}

func (argument *Integer) NormalizeState() {
	argument.Value = argument.Default
}

func (argument *Integer) ValidateLazy() (err error) {
	if argument.Type != TYPE_INTEGER {
		err = errors.Join(err, errors.New("Integer is not of type INTEGER"))
	}
	if argument.GetDefault() < argument.GetMinValue() {
		err = errors.Join(err, errors.New("default value is lower than minvalue"))
	}
	if argument.GetDefault() > argument.GetMaxValue() {
		err = errors.Join(err, errors.New("default value is higher than maxvalue"))
	}
	return
}

func (argument *Integer) Reset() {
	argument.Value = argument.Default
}
