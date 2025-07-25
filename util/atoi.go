package util

import (
	"fmt"
	"strconv"
)

func AtoiToUint(str string) (uint, error) {
	// Parse the string to an int
	intVal, err := strconv.Atoi(str)
	if err != nil || intVal < 0 {
		return 0, fmt.Errorf("invalid string: %s", str)
	}
	// Cast to uint
	return uint(intVal), nil
}
