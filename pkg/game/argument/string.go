package argument

import "errors"

var _ ArgumentWithValue[string] = (*String)(nil)

type String struct {
	Base    `yaml:",inline"`
	Default string `json:"default" yaml:"default"`
	Value   string `json:"value" yaml:"value"`
}

func (argument *String) GetDefault() string {
	return argument.Default
}

func (argument *String) GetValue() string {
	return argument.Value
}

func (argument *String) Parse(sep Seperator) (arg string, err error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	arg = argument.GetArgument() + *argument.GetSeperator(sep).ArgumentValue + argument.GetValue()
	return
}

func (argument *String) NormalizeState() {
	argument.Value = argument.Default
}

func (argument *String) ValidateLazy() (err error) {
	if argument.Type != TYPE_STRING {
		err = errors.Join(err, errors.New("String is not of type STRING"))
	}
	if len(argument.GetDefault()) == 0 {
		err = errors.Join(err, errors.New("default value is empty"))
	}
	return
}

func (argument *String) Reset() {
	argument.Value = argument.Default
}
