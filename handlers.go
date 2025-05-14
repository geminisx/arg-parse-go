package main 

import (
	"slices"
	"strings"
)

func (structure *Tree) NodeErrorHandler(error error, node Node) error {
	node.Error = append(node.Error, error)
	structure.Root.Nodes = append(structure.Root.Nodes, node)
	structure.NodeFlag = true
	return error
}


func isCommand (parameter string) bool {
	return strings.HasPrefix(parameter, "-")
}

func isValidCommand (y []string, z string) bool {
	return slices.Contains(y, z) 
}
