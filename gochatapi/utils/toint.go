package utils

import "strconv"

// ToInt ...
func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
