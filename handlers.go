package main 

import (
	"slices"
	"strings"
)



func isCommand (parameter string) bool {
	return strings.HasPrefix(parameter, "-")
}

func isValidCommand (y []string, z string) bool {
	return slices.Contains(y, z) 
}
