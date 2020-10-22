package cli

import (
	"flag"
	"strings"
	"testing"
)

func TestFishCompletion(t *testing.T) {
	// Given
	app := testApp()
	app.Flags = append(app.Flags, &PathFlag{
		Name:      "logfile",
		TakesFile: true,
	})

	// When
	res, err := app.ToFishCompletion()

	// Then
	expect(t, err, nil)
	expectFileContent(t, "testdata/expected-fish-full.fish", res)
}

func TestFishAddFileFlag(t *testing.T) {
	flags := []Flag{
		&GenericFlag{
			Name:      "jen-and-eric",
			TakesFile: false,
		},
		&PathFlag{
			Name:      "logfile",
			TakesFile: false,
		},
		&BoolFlag{
			Name: "dry-run",
		},
		&IntFlag{
			Name: "depth",
		},
	}

	for _, flag := range flags {
		completion := &strings.Builder{}
		fishAddFileFlag(flag, completion)
		expect(t, completion.String(), " -f")
	}

	flags = []Flag{
		&GenericFlag{
			Name:      "jen-and-eric",
			TakesFile: true,
		},
		&PathFlag{
			Name:      "logfile",
			TakesFile: true,
		},
	}

	for _, flag := range flags {
		completion := &strings.Builder{}
		fishAddFileFlag(flag, completion)
		t.Log(flag.Names())
		expect(t, completion.String(), "")
	}
}

type bareIntFlag int

func (bareIntFlag) Apply(*flag.FlagSet) error { return nil }
func (bareIntFlag) Names() []string           { return []string{} }
func (bareIntFlag) IsSet() bool               { return false }
func (bareIntFlag) String() string            { return "bareIntFlag" }

type bareBoolFlag bool

func (*bareBoolFlag) Apply(*flag.FlagSet) error { return nil }
func (*bareBoolFlag) Names() []string           { return []string{} }
func (*bareBoolFlag) IsSet() bool               { return false }
func (*bareBoolFlag) String() string            { return "bareBoolFlag" }

func TestFishAddFileFlagUnexpectedInput(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("test panicked: %v", err)
		}
	}()
	bareInt := bareIntFlag(5)
	bareBool := bareBoolFlag(false)
	flags := []Flag{
		bareInt,
		&bareBool,
		(*bareBoolFlag)(nil),
	}
	for _, flag := range flags {
		completion := &strings.Builder{}
		fishAddFileFlag(flag, completion)
		expect(t, completion.String(), " -f")
	}
}
