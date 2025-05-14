package main

import (
	"fmt"
	"os"
	"reflect"
	"slices"
)

func (structure *Tree) main(args []string) {
	var rootFlag bool = false

	a := Command{
		Alias:          []string{"-a","--a"},
		AcceptsCommands:  AcceptsCommands{Bool: true, comandos: []string{"-b","--b","-c","--c","-3"}},
		AcceptsValues:  AcceptsValues{Bool: true, Types: []Types{{typeString: "string", typeArrayString: true}}},
		QualifiedName:  "a command",
	}
	b := Command{
		Alias:          []string{"-b","--b"},
		AcceptsCommands:  AcceptsCommands{Bool: true, comandos: []string{"-1","-2"}},
		AcceptsValues:  AcceptsValues{Bool: true, Types: []Types{{typeString: "string", typeArrayString: true}}},
		QualifiedName:  "b command",
	}
	c := Command{
		Alias:          []string{"-c","--c"},
		AcceptsCommands:  AcceptsCommands{Bool: false, comandos: []string{}},
		AcceptsValues:  AcceptsValues{Bool: false, Types: []Types{}},
		QualifiedName:  "c command",
	}
	_3 := Command{
		Alias:          []string{"-3","--3"},
		AcceptsCommands:  AcceptsCommands{Bool: false, comandos: []string{}},
		AcceptsValues:  AcceptsValues{Bool: true, Types: []Types{{typeString: "string", typeArrayString: false}}},
		QualifiedName:  "3 command",
	}

	commands := []Command{a, b, c,_3}

	for _, i := range commands {
		if (slices.Contains(i.Alias, args[0])){
			structure.Root = &Root{Command: i}
			rootFlag = true
		}
	}

	if rootFlag {
		structure.NodeFlag = false

		root := structure.Root

		acceptsCommands := root.Command.AcceptsValues.Bool
		acceptsValues   := root.Command.AcceptsCommands.Bool

		nodeSubCommands := root.Command.AcceptsCommands.comandos
		
		args = args[1:]

		if !acceptsValues {
			if len(args) > 1 {
				if len(args) > 1 { 
					structure.Root.Error = append(structure.Root.Error ,fmt.Errorf("[ERROR] '%s' does not accept values", structure.Root.Command.Alias))
					return
				} 
			} 
		}

		for structure.Cursor = 0; structure.Cursor < len(args); structure.Cursor++ {
			cursor := structure.Cursor
			if isCommand(args[cursor]) { 
				if acceptsCommands {
						if slices.Contains(nodeSubCommands, args[cursor]) { 
							for _, j := range commands { 
								if (slices.Contains(j.Alias, args[cursor])){
									error := structure.nodeParse(args[cursor + 1:], j)
									if error != nil { fmt.Println(error);return }
									break 
								}
							} 
						} else { 
							structure.Root.Error = append(structure.Root.Error, fmt.Errorf("[ERROR] '%s' does not accept '%s' as a subCommand", root.Command.Alias, args[cursor]))
							break
						}
				}
			} else { 
				if !structure.NodeFlag { 
					if acceptsValues { 
						root.Value = append(root.Value, args[cursor]) 
					} else { 
						if len(args) > 1 { 
							structure.Root.Error = append(structure.Root.Error ,fmt.Errorf("[ERROR] '%s' does not accept values", structure.Root.Command.Alias))
							break
						} 
					} 
				} 
			}
		}
	}
	// can return structure if there are no errors
}

func (structure *Tree) nodeParse(args []string, j Command) error {
	node := Node{Command: j}
	root := structure.Root

	subCommandValueLock := node.Command.AcceptsValues.Bool
	Types               := node.Command.AcceptsValues.Types
	rootCommands        := root.Command.AcceptsCommands.comandos
	nodeCommands        := node.Command.AcceptsCommands.comandos

	for cursor := range args {
		
		if subCommandValueLock {
			if isCommand(args[cursor]) { if !isValidCommand(nodeCommands, args[cursor]) {
				if slices.Contains(rootCommands, args[cursor]) {
					root.Nodes = append(root.Nodes, node)
					return nil
				}
				error := fmt.Errorf("[ERROR] '%s' + '%s' concat, doesn't recognize: '%s'",root.Command.Alias,node.Command.Alias, args[cursor])
				return structure.NodeErrorHandler(error,node)
				} }
				
			structure.Cursor += cursor

			for _, j := range Types {
				// not valid argument type
				if reflect.TypeOf(args[cursor]).String() != j.typeString {
					error := fmt.Errorf("[ERROR] Invalid data type passed'%s', node requires: '%s'", reflect.TypeOf(args[cursor]).String(), j.typeString)
					return structure.NodeErrorHandler(error,node)
				}
				// command doesn't accept arrays 
				if !j.typeArrayString && len(node.Value) < 1 {
					error := fmt.Errorf("[ERROR] Command '%s' doesn't accept array values of type '%s'", node.Command.Alias, j.typeString)
					return structure.NodeErrorHandler(error,node)
				}
				node.Value = append(node.Value, args[cursor])
			}
		} else {
			if slices.Contains(rootCommands, args[cursor]) {
				return nil
			} else {
				return nil
			}
		}
		
	}

	root.Nodes = append(root.Nodes, node)
	structure.NodeFlag = true
	return nil
}

func main() {
	x := Tree{}
	x.main(os.Args[1:])
	
	fmt.Println()
	fmt.Println("Command: ", x.Root.Command.QualifiedName)
	fmt.Println("Values: ", x.Root.Value)
	fmt.Println()
	fmt.Println("---Node-Components---")
	for _, i := range x.Root.Nodes {
		fmt.Println("Node Command: ", i.Command.QualifiedName)
		fmt.Println("Node Values: ",  i.Value)
		fmt.Println()
	}
	fmt.Println("---ROOT-ERRORS--------")
	for _, i := range x.Root.Error {
		fmt.Println(i)
	}

	fmt.Println("---NODE-ERRORS--------")
	for _, i := range x.Root.Nodes {
		for _, j := range i.Error {
			fmt.Println(j)
		}
	}
}
