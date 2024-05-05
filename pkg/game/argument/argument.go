package argument

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type Argument interface {
	GetType() Type
	IsMandatory() bool
	IsDisabled() bool
	Disable()
	Enable()
	GetSeperator(Seperator) Seperator
	GetArgument() string
	GetName() string
	Parse(Seperator) (string, error)
	NormalizeState()
	ValidateLazy() error
	Reset()
}

type ArgumentWithValue[T any] interface {
	Argument
	GetDefault() T
	GetValue() T
}

type NumberValue interface {
	int | float32
}

type ArgumentWithNumberValue[T NumberValue] interface {
	ArgumentWithValue[T]
	GetMinValue() T
	GetMaxValue() T
}

type Arguments struct {
	Seperator *Seperator `json:"seperator,omitempty" yaml:"seperator,omitempty"`
	Arguments []Argument `json:"items,omitempty" yaml:"items,omitempty"`
}

type tmpYAMLArguments struct {
	Seperator *Seperator  `json:"seperator" yaml:"seperator"`
	Arguments []yaml.Node `json:"items" yaml:"items"`
}

type tmpJSONArguments struct {
	Seperator *Seperator        `json:"seperator" yaml:"seperator"`
	Arguments []json.RawMessage `json:"items" yaml:"items"`
}

// TODO YAML and JSON unmarshl can be unified in one function by using their Decode interfaces
// we need to declare a new interface Decode and assign it the decode functions --> then YAML and
// JSON Decoder implement that interface and it can be given as input to the function
// function: unmarshl(Argument, Decoder) (Argument, err)
func (arguments *Arguments) UnmarshalYAML(value *yaml.Node) (err error) {
	var tmpArguments tmpYAMLArguments
	err = value.Decode(&tmpArguments)
	if err != nil {
		return
	}
	arguments.Seperator = tmpArguments.Seperator
	for _, arg := range tmpArguments.Arguments {
		var item Base
		err = arg.Decode(&item)
		if err != nil {
			return
		}
		switch item.Type {
		case TYPE_BOOLEAN:
			var boolArgument Boolean
			err = arg.Decode(&boolArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &boolArgument)
		case TYPE_STRING:
			var strArgument String
			err = arg.Decode(&strArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &strArgument)
		case TYPE_INTEGER:
			var intArgument Integer
			err = arg.Decode(&intArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &intArgument)
		case TYPE_FLOAT:
			var floatArgument Float
			err = arg.Decode(&floatArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &floatArgument)
		case TYPE_ENUM:
			var enumArgument Enum
			err = arg.Decode(&enumArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &enumArgument)
		case TYPE_CONNECT:
			var connectArgument Connect
			err = arg.Decode(&connectArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &connectArgument)
		case TYPE_BASE:
			arguments.Arguments = append(arguments.Arguments, &item)
		}
	}
	return
}

func (arguments *Arguments) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var tmpArguments tmpJSONArguments
	err = json.Unmarshal(data, &tmpArguments)
	if err != nil {
		return
	}
	arguments.Seperator = tmpArguments.Seperator
	for _, arg := range tmpArguments.Arguments {
		var item Base
		err = json.Unmarshal(arg, &item)
		if err != nil {
			return
		}
		switch item.Type {
		case TYPE_BOOLEAN:
			var boolArgument Boolean
			err = json.Unmarshal(arg, &boolArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &boolArgument)
		case TYPE_STRING:
			var strArgument String
			err = json.Unmarshal(arg, &strArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &strArgument)
		case TYPE_INTEGER:
			var intArgument Integer
			err = json.Unmarshal(arg, &intArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &intArgument)
		case TYPE_FLOAT:
			var floatArgument Float
			err = json.Unmarshal(arg, &floatArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &floatArgument)
		case TYPE_ENUM:
			var enumArgument Enum
			err = json.Unmarshal(arg, &enumArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &enumArgument)
		case TYPE_CONNECT:
			var connectArgument Connect
			err = json.Unmarshal(arg, &connectArgument)
			if err != nil {
				return
			}
			arguments.Arguments = append(arguments.Arguments, &connectArgument)
		case TYPE_BASE:
			arguments.Arguments = append(arguments.Arguments, &item)
		}
	}
	return
}

func (arguments *Arguments) Parse() (args []string, err error) {
	tmpArgs := []string{}
	tmpSep := []string{}
	tmpArguments := ""
	for _, argument := range arguments.Arguments {
		if argument.IsDisabled() || argument.GetType() == TYPE_CONNECT {
			continue
		}
		arg, err := argument.Parse(*arguments.Seperator)
		if err != nil {
			return args, err
		}
		tmpArgs = append(tmpArgs, arg)
		seperator := argument.GetSeperator(*arguments.Seperator)
		tmpSep = append(tmpSep, *seperator.Arguments)
	}

	for index, tmpArg := range tmpArgs {
		if index+1 < len(tmpArgs) {
			tmpArguments = tmpArguments + tmpArg + tmpSep[index]
		} else {
			tmpArguments = tmpArguments + tmpArg
		}
	}
	if len(tmpArguments) > 0 {
		args = strings.Split(tmpArguments, " ")
	}
	return
}

func (arguments *Arguments) ValidateLazy() (err error) {
	for _, argument := range arguments.Arguments {
		validateErr := argument.ValidateLazy()
		if validateErr != nil {
			validateErr = fmt.Errorf("%s: %w", argument.GetName(), validateErr)
		}
		err = errors.Join(err, validateErr)
	}
	return
}

func (arguments *Arguments) NormalizeState() {
	if len(arguments.Arguments) > 0 && arguments.Seperator == nil {
		arguments.Seperator = &SEPERATOR_SPACE
	}
	for _, argument := range arguments.Arguments {
		argument.NormalizeState()
	}
}

func (left Arguments) Equal(right Arguments) bool {
	return left.Seperator == right.Seperator &&
		slices.EqualFunc(left.Arguments, right.Arguments, func(l, r Argument) bool {
			if l.GetType() != r.GetType() {
				return false
			}
			if l.GetType() == TYPE_ENUM {
				return l.(*Enum).Equal(r.(*Enum))
			}
			return l == r
		})
}
