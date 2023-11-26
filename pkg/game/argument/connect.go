package argument

import (
	"errors"
	"strings"
)

var _ Argument = (*Connect)(nil)

type Connect struct {
	Type     Type   `json:"type" yaml:"type"`
	Argument string `json:"argument" yaml:"argument"`
}

func (argument *Connect) GetType() Type {
	return argument.Type
}

func (argument *Connect) IsMandatory() bool {
	return true
}

func (argument *Connect) IsDisabled() bool {
	return false
}

func (argument *Connect) Disable() {
}

func (argument *Connect) Enable() {
}

func (argument *Connect) GetSeperator(seperator Seperator) (sep Seperator) {
	return Seperator{}
}

func (argument *Connect) GetArgument() string {
	return argument.Argument
}

func (argument *Connect) GetName() string {
	return "Connect Argument"
}

func (argument *Connect) Parse(sep Seperator) (string, error) {
	if argument.IsDisabled() && !argument.IsMandatory() {
		return "", nil
	}
	return argument.GetArgument(), nil
}

func (argument *Connect) NormalizeState() {}

func (argument *Connect) ValidateLazy() (err error) {
	if argument.Type != TYPE_CONNECT {
		err = errors.Join(err, errors.New("Connect is not of type CONNECT"))
	}
	return
}

func (argument *Connect) ParseWithIP(ip string) (string, error) {
	arg, err := argument.Parse(Seperator{})
	if err != nil {
		return "", err
	}
	return strings.Replace(arg, "?", ip, 1), nil
}

func (argument *Connect) Reset() {}
