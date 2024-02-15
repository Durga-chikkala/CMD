
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Command struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	Action      func(flags map[string]string) error
}

type CLI struct {
	Commands map[string]*Command
}

func NewCLI() CLI {
	return CLI{
		Commands: make(map[string]*Command),
	}
}

// RegisterCommand registers a new command
func (c *CLI) RegisterCommand(cmd Command) {
	c.Commands[cmd.Name] = &cmd
}

// AddCommand adds a command to the CLI
func (cli *CLI) AddCommand(cmd *Command) {
	cli.Commands[cmd.Name] = cmd
}

// Run executes the command-line interface
func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	cmdName := os.Args[1]
	for _, cmd := range cli.Commands {
		if cmd.Name == cmdName {
			cmd.Flags.Parse(os.Args[2:])
			flagMap := make(map[string]string)
			cmd.Flags.Visit(func(f *flag.Flag) {  // TODO: ADD TYPE CHECKING LOGIC HERE
				flagMap[f.Name] = f.Value.String()
			})
			err := cmd.Action(flagMap)
			if err != nil {
				log.Fatalf("Error executing command '%s': %v\n", cmd.Name, err)
			}
			return
		}
	}

	fmt.Printf("Error: Unknown command '%s'\n", cmdName)
	cli.printUsage()
	os.Exit(1)
}

const (
	STRING = "string"
	INT    = "int"
	BOOL   = "bool"
)

type Flags struct {
	Name        string
	Type        string
	Description string
}

func (cli *CLI) AddFlags(command string, flags []Flags) {
	helloFlags := flag.NewFlagSet(cli.Commands[command].Name, flag.ExitOnError)
	for _, value := range flags {
		switch value.Type {
		case STRING:
			helloFlags.String(value.Name, "", value.Description)
		case INT:
			helloFlags.Int(value.Name, 0, value.Description)
		case BOOL:
			helloFlags.Bool(value.Name, false, value.Description)
		}
	}

	cli.Commands[command].Flags = helloFlags
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage: mytool <command> [options]")
	fmt.Println("Available commands:")
	for _, cmd := range cli.Commands {
		fmt.Printf("  %s\t%s\n", cmd.Name, cmd.Description)
	}
}

// Example function to demonstrate command action
func helloAction(flags map[string]string) error {
	name, ok := flags["name"]
	if !ok {
		name = "world"
	}
	fmt.Printf("Hello, %s!\n", name)
	return nil
}

func main() {
	 cli:=NewCLI()
	 
	 cli.RegisterCommand(Command{
		 Name:        "hello",
		 Description: "Print hello world",
		 Action: func(flags map[string]string) error {
			 return helloAction(flags)
		 },
	 })

	cli.AddFlags("hello", []Flags{{Name: "name", Type: STRING, Description: "Tell your Name"}})

	cli.Run()
}
