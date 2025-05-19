package goat

import (
	"fmt"
	"reflect"
	"slices"
)

func (s *Tree) Main(args []string, commands []Command) {
	var rootFlag bool = false
	for _, i := range commands {
		if (i.TLName == args[0]){
			s.Root = &Root{Command: i}
			rootFlag = true
		}
	}

	if rootFlag {
		s.structuration(args, commands)
	}
	// can return structure if there are no errors
}

func (structure *Tree) structuration(args []string, commands []Command) {
	structure.NodeFlag = false

		root     := structure.Root
		rootName := structure.Root.Command.TLName

		acceptsCommands := root.Command.AcceptsCommands.Bool
		acceptsValues   := root.Command.AcceptsValues.Bool

		nodeSubCommands := root.Command.AcceptsCommands.commands
		
		args = args[1:]

		if !acceptsValues && !acceptsCommands && len(args) > 1 { structure.RootErrorHandler(fmt.Errorf(errorStack[0], rootName)); return }

		for structure.Cursor = 0; structure.Cursor < len(args); structure.Cursor++ {
			cursor := structure.Cursor
			if isCommand(args[cursor]) { 
				if acceptsCommands {
					if slices.Contains(nodeSubCommands, args[cursor]) { 
						// bug: design hierarchy systems
						for _, j := range commands {
							// first it matches cannot be intended command 
							if (args[cursor] == j.FQsubCommandName) {
								error := structure.nodeParse(args[cursor + 1:], j)
								if error != nil { fmt.Println(error);return  }
								break 
							} 
						}
					} else { structure.RootErrorHandler(fmt.Errorf(errorStack[2], root.Command.TLName, args[cursor]));break }
				} else { structure.RootErrorHandler(fmt.Errorf(errorStack[1], root.Command.TLName ));break }
			} else { 
				if !structure.NodeFlag && acceptsValues { root.Value = append(root.Value, args[cursor]) 
				} else { if len(args) > 1 { structure.RootErrorHandler(fmt.Errorf(errorStack[3], structure.Root.Command.TLName));break } }
			}
		}
}

func (structure *Tree) nodeParse(args []string, j Command) error {
	//structure handler ->+ accepts:
	node := Node{Command: j}
	
	//readable
	root := structure.Root
	subCommandValueLock := node.Command.AcceptsValues.Bool
	Types               := node.Command.AcceptsValues.Types
	rootCommands        := root.Command.AcceptsCommands.commands
	nodeCommands        := node.Command.AcceptsCommands.commands

	for cursor := range args {
		//is root child able to hold values?
		if subCommandValueLock {
			if isCommand(args[cursor]) { if !isValidCommand(nodeCommands, args[cursor]) {
				if slices.Contains(rootCommands, args[cursor]) {
					root.Nodes = append(root.Nodes, node)
					return nil
				}
				error := fmt.Errorf("[ERROR] '%s','%s' iteration doesn't recognize: '%s'",root.Command.TLName,node.Command.FQsubCommandName, args[cursor])
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
					error := fmt.Errorf("[ERROR] Command '%s' doesn't accept array values of type '%s'", node.Command.FQsubCommandName, j.typeString)
					return structure.NodeErrorHandler(error,node)
				}
				node.Value = append(node.Value, args[cursor])
			}
		} else {
			//does top level accept command?
			if slices.Contains(rootCommands, args[cursor]) {
				return nil
			} else {
				return fmt.Errorf("[ERROR] '%s' doesnt accept '%s' as a value",node.Command.FQsubCommandName, args[cursor])
			}
		}
	}

	root.Nodes = append(root.Nodes, node)
	structure.NodeFlag = true
	return nil
}

func (structure *Tree) RootErrorHandler(error error) {
	structure.Root.Error = append(structure.Root.Error, error)
}

func (structure *Tree) NodeErrorHandler(error error, node Node) error {
	node.Error = append(node.Error, error)
	structure.Root.Nodes = append(structure.Root.Nodes, node)
	structure.NodeFlag = true
	return error
}
