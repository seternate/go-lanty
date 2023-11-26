package argument

import (
	"errors"
	"strconv"
)

var _ ArgumentWithNumberValue[float32] = (*Float)(nil)

type Float struct {
	Base     `yaml:",inline"`
	Default  float32 `json:"default" yaml:"default"`
	Value    float32 `json:"value" yaml:"value"`
	MinValue float32 `json:"minvalue" yaml:"minvalue"`
	MaxValue float32 `json:"maxvalue" yaml:"maxvalue"`
}

func (argument *Float) GetDefault() float32 {
	return argument.Default
}

func (argument *Float) GetValue() float32 {
	return argument.Value
}

func (argument *Float) GetMinValue() float32 {
	return argument.MinValue
}

func (argument *Float) GetMaxValue() float32 {
	return argument.MaxValue
}

func (argument *Float) Parse(sep Seperator) (arg string, err error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	arg = argument.GetArgument() + *argument.GetSeperator(sep).ArgumentValue + strconv.FormatFloat(float64(argument.GetValue()), 'f', -1, 32)
	return
}

func (argument *Float) NormalizeState() {
	argument.Value = argument.Default
}

func (argument *Float) ValidateLazy() (err error) {
	if argument.Type != TYPE_FLOAT {
		err = errors.Join(err, errors.New("Float is not of type FLOAT"))
	}
	if argument.GetDefault() < argument.GetMinValue() {
		err = errors.Join(err, errors.New("default value is lower than minvalue"))
	}
	if argument.GetDefault() > argument.GetMaxValue() {
		err = errors.Join(err, errors.New("default value is higher than maxvalue"))
	}
	return
}

func (argument *Float) Reset() {
	argument.Value = argument.Default
}
