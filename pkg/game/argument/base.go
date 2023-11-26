package argument

import "errors"

var _ Argument = (*Base)(nil)

type Base struct {
	Type      Type       `json:"type" yaml:"type"`
	Mandatory bool       `json:"mandatory" yaml:"mandatory"`
	Disabled  bool       `json:"disabled" yaml:"disabled"`
	Seperator *Seperator `json:"seperator,omitempty" yaml:"seperator,omitempty"`
	Argument  string     `json:"argument" yaml:"argument"`
	Name      string     `json:"name" yaml:"name"`
}

func (argument *Base) GetType() Type {
	return argument.Type
}

func (argument *Base) IsMandatory() bool {
	return argument.Mandatory
}

func (argument *Base) IsDisabled() bool {
	return argument.Disabled
}

func (argument *Base) Disable() {
	argument.Disabled = true
}

func (argument *Base) Enable() {
	argument.Disabled = false
}

func (argument *Base) GetSeperator(seperator Seperator) (sep Seperator) {
	if argument.Seperator != nil && argument.Seperator.ArgumentValue != nil {
		sep.ArgumentValue = new(string)
		*sep.ArgumentValue = *argument.Seperator.ArgumentValue
	} else if seperator.ArgumentValue != nil {
		sep.ArgumentValue = new(string)
		*sep.ArgumentValue = *seperator.ArgumentValue
	}
	if argument.Seperator != nil && argument.Seperator.Arguments != nil {
		sep.Arguments = new(string)
		*sep.Arguments = *argument.Seperator.Arguments
	} else if seperator.Arguments != nil {
		sep.Arguments = new(string)
		*sep.Arguments = *seperator.Arguments
	}
	return
}

func (argument *Base) GetArgument() string {
	return argument.Argument
}

func (argument *Base) GetName() string {
	return argument.Name
}

func (argument *Base) Parse(sep Seperator) (string, error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	return argument.GetArgument(), nil
}

func (argument *Base) NormalizeState() {}

func (argument *Base) ValidateLazy() (err error) {
	if argument.Type != TYPE_BASE {
		err = errors.Join(err, errors.New("Base is not of type BASE"))
	}
	return
}

func (argument *Base) Reset() {}
