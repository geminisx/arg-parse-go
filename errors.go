package goat

var errorStack = map[int]string{
	0:"[ERROR] '%s' is a standalone",
	1:"[ERROR] '%s' does not accept sub-commands",
	2:"[ERROR] '%s' doesnt't accept '%s' sub-command",
	3:"[ERROR] '%s' does not accept values",
} 
