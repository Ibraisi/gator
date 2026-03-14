package commands

import "testing"

func Test_Run_UnknownCommand(t *testing.T) {
	// Arrange
	cmds := New()

	// Act
	err := cmds.Run(&State{}, Command{Name: "nonexistent"})

	// Assert
	if err == nil {
		t.Fatal("expected error for unknown command, got nil")
	}
}

func Test_Run_CallsRegisteredHandler(t *testing.T) {
	// Arrange
	cmds := New()
	called := false
	cmds.Register("ping", func(s *State, cmd Command) error {
		called = true
		return nil
	})

	// Act
	err := cmds.Run(&State{}, Command{Name: "ping"})

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected handler to be called")
	}
}

func Test_Run_PassesArgsToHandler(t *testing.T) {
	// Arrange
	cmds := New()
	var got []string
	cmds.Register("echo", func(s *State, cmd Command) error {
		got = cmd.Args
		return nil
	})

	// Act
	cmds.Run(&State{}, Command{Name: "echo", Args: []string{"a", "b"}})

	// Assert
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("expected args [a b], got %v", got)
	}
}

func Test_Register_OverwritesPreviousHandler(t *testing.T) {
	// Arrange
	cmds := New()
	first, second := false, false
	cmds.Register("cmd", func(s *State, cmd Command) error { first = true; return nil })
	cmds.Register("cmd", func(s *State, cmd Command) error { second = true; return nil })

	// Act
	cmds.Run(&State{}, Command{Name: "cmd"})

	// Assert
	if first {
		t.Fatal("first handler should have been overwritten")
	}
	if !second {
		t.Fatal("second handler should be called")
	}
}
