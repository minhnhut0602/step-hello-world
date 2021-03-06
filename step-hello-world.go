package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coinbase/step/machine"
	"github.com/coinbase/step/utils/run"
)

func main() {
	var arg, command string
	switch len(os.Args) {
	case 1:
		fmt.Println("Starting Lambda")
		run.Lambda(StateMachine())
	case 2:
		command = os.Args[1]
		arg = ""
	case 3:
		command = os.Args[1]
		arg = os.Args[2]
	default:
		printUsage() // Print how to use and exit
	}

	switch command {
	case "json":
		run.JSON(StateMachine())
	case "exec":
		run.Exec(StateMachine())(&arg)
	default:
		printUsage() // Print how to use and exit
	}
}

func printUsage() {
	fmt.Println("Usage: step-hello-world <json|exec> <arg> (No args starts Lambda)")
	os.Exit(0)
}

// StateMachine returns the StateMachine
// replacing the `Resource` in Task states with the lambdaArn
func StateMachine() (*machine.StateMachine, error) {
	state_machine, err := machine.FromJSON([]byte(`{
    "Comment": "Hello World",
    "StartAt": "HelloFn",
    "States": {
      "HelloFn": {
        "Type": "Pass",
        "Result": "Hello",
        "ResultPath": "$.Task",
        "Next": "Hello"
      },
      "Hello": {
        "Type": "Task",
        "Comment": "Deploy Step Function",
        "End": true
      }
    }
  }`))

	if err != nil {
		return nil, err
	}

	state_machine.SetResourceFunction("Hello", HelloHandler)

	return state_machine, nil
}

////////////
// HANDLERS
////////////

type Hello struct {
	Greeting string
}

// Handlers must conform to function type
// func (context.Context, interface{}) (interface{}, error)
// The input is auto-unmarshalled and marsheled from and to JSON
func HelloHandler(_ context.Context, hello *Hello) (interface{}, error) {
	if hello.Greeting == "" {
		hello.Greeting = "Giday Mate"
	}
	return hello, nil
}
