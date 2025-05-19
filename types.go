package goat

//top level structures
type Tree struct {
	Root     *Root
	Cursor   int
	NodeFlag bool
}

type Root struct {
	Command          Command
	Value            []string
	Nodes            []Node
	Error            []error
}

type Node struct {
	Command          Command
	Value            []string
	Error            []error
}

type Command struct {
	TLName 			 string
	FQsubCommandName string
	AcceptsCommands  AcceptsCommands
	AcceptsValues    AcceptsValues
}

//Parameter specifiers 
type AcceptsCommands struct {
	Bool bool
	commands []string
}

type AcceptsValues struct {
	Bool bool
	Types []Types
}

type Types struct {
	typeString string
	typeArrayString bool
}