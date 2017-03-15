package curscraft

import (
	"errors"
	"fmt"
)

type Arg struct {
	Name     string
	Command  string
	Shortcut string
	Help     string
	Number   int
}

func (a *Arg) checkNumberSubArgument(args int) error {
	strErr := fmt.Sprintf("Argument invalid for command %v", a.Command)
	if a.Number != args && a.Number > 0 {
		return errors.New(strErr)
	}

	return nil
}

type Argts struct {
	Args []Arg
}

func (a *Argts) Help() {
	fmt.Println("Usage : ")
	for _, argument := range a.Args {
		fmt.Println(argument.Help)
	}
}

func (a *Argts) CheckArgument(args []string) (string, error) {
	var argName string
	var err error
	var find bool

	// No Arguments
	if len(args) == 1 {
		return argName, errors.New("Try a command.")
	}

	// Arguments
	for _, argument := range a.Args {
		if args[1] == argument.Command || args[1] == argument.Shortcut {
			find = true
			err = argument.checkNumberSubArgument(len(args) - 2)
			return argument.Name, err
		}
	}

	// Argument don't find
	if find == false {
		return argName, errors.New("No such command.")
	}

	return argName, err
}
