package cmd

import (
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli := NewCLI([]string{"test"})

	if len(cli.args) != 1 {
		t.Errorf("expected 1, got %d", len(cli.args))
	}

	if len(cli.commands) != 0 {
		t.Errorf("expected 0, got %d", len(cli.commands))
	}
}

func TestAddCommand(t *testing.T) {
	cli := NewCLI([]string{"test"})

	cli.AddCommand(NewCommand("test", func(args []string) {}))

	if len(cli.commands) != 1 {
		t.Errorf("expected 1, got %d", len(cli.commands))
	}
}

func TestRun(t *testing.T) {
	cli := NewCLI([]string{})

	type Test struct {
		Args        []string
		Executed    bool
		RecivedArgs []string

		ExpectedExecuted bool
		ExpectedArgs     []string
	}

	tests := make(map[string]*Test)

	tests["foo"] = &Test{
		Args:        []string{"foo", "test"},
		Executed:    true,
		RecivedArgs: []string{"test"},

		ExpectedExecuted: true,
		ExpectedArgs:     []string{"test"},
	}
	tests["bar"] = &Test{
		Args:        []string{},
		Executed:    false,
		RecivedArgs: []string{"test"},

		ExpectedExecuted: false,
		ExpectedArgs:     []string{},
	}

	MakeTestHandler := func(name string) func([]string) {
		return func(args []string) {
			tests[name].Executed = true
			tests[name].RecivedArgs = args
		}
	}

	for name, _ := range tests {
		cli.AddCommand(NewCommand(name, MakeTestHandler(name)))
	}

	cli.Run()

	for name, test := range tests {
		if test.Executed {
			if !test.ExpectedExecuted {
				t.Errorf("command %s excepted to be executed, but not", name)
			}

			if len(test.RecivedArgs) != len(test.ExpectedArgs) {
				t.Errorf("command %s excepted args %v, but got %v", name, test.ExpectedArgs, test.RecivedArgs)
			}

			for i, arg := range test.ExpectedArgs {
				if arg != test.RecivedArgs[i] {
					t.Errorf("command %s excepted args %d: %v, but got %v", name, i, test.ExpectedArgs, test.RecivedArgs)
				}
			}

		} else {
			if test.ExpectedExecuted {
				t.Errorf("command %s excepted not to be executed, but was", name)
			}
		}
	}
}
