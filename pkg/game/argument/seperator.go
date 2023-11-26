package argument

type Seperator struct {
	Arguments     *string `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	ArgumentValue *string `json:"argumentvalue,omitempty" yaml:"argumentvalue,omitempty"`
}

var (
	SEPERATOR_SPACE = newSpaceSeperator()
)

func newSpaceSeperator() Seperator {
	space1 := new(string)
	space2 := new(string)
	*space1 = " "
	*space2 = " "
	return Seperator{space1, space2}
}
