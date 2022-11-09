package main

import (
	"fmt"
	"os"

	"github.com/busser/tftree/internal/format"
	"github.com/busser/tftree/internal/terraform"
	"github.com/busser/tftree/internal/tftree"
)

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(format.Error(err))
		os.Exit(1)
	}
}

func run() error {
	var workdir string
	switch {
	case len(os.Args) > 2:
		return fmt.Errorf("too many arguments, must be either 0 or 1")
	case len(os.Args) == 2:
		workdir = os.Args[1]
	default:
		workdir = "."
	}

	tf := terraform.NewRunner(workdir)

	logln("Running \"terraform init\"...")
	err := tf.Init()
	if err != nil {
		return err
	}

	logln("Running \"terraform plan\"...")
	plan, err := tf.Plan()
	if err != nil {
		return err
	}

	root, err := tftree.New(plan, workdir)
	if err != nil {
		return err
	}

	fmt.Print(format.Module(root))

	return nil
}

func logln(msg string) {
	fmt.Fprint(os.Stderr, format.Info(msg))
}
