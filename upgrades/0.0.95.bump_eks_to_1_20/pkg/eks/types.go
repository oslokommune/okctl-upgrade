package eks

import "errors"

var errNothingToDo = errors.New("nothing to do")

const (
	eksSourceVersion = "1.19"
	eksTargetVersion = "1.20"
)
