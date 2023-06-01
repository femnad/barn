package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/femnad/barn/selection"
)

type commonArgs struct {
	Config string `arg:"-f,--file" default:"~/.config/barn/barn.yml" help:"Config file path"`
}

type selectCmd struct {
	commonArgs
	Id        string `arg:"-i,--id,required"`
	Selection string `arg:"positional"`
}

type outputCmd struct {
	commonArgs
	Id             string `arg:"positional"`
	ShowZeroCounts bool   `arg:"-z,--zero" help:"Show entries with zero counts"`
}

type args struct {
	Select *selectCmd `arg:"subcommand:select" help:"Select based on given choices and update counts"`
	Output *outputCmd `arg:"subcommand:output" help:"Show stored entries for the given selection ID"`
}

func showSelections(config, id string) {
	err := selection.Show(config, id)
	if err != nil {
		log.Fatalf("error getting selections for id %s: %v", id, err)
	}
}

func markSelection(config, id, choice string) {
	err := selection.Mark(config, id, choice)
	if err != nil {
		log.Fatalf("error marking selection as %s for id %s: %v", choice, id, err)
	}
}

func doSelect(cmd *selectCmd) {
	if cmd.Selection == "" {
		showSelections(cmd.Config, cmd.Id)
		return
	}

	markSelection(cmd.Config, cmd.Id, cmd.Selection)
}

func doOutput(cmd *outputCmd) {
	err := selection.Iterate(cmd.Config, cmd.Id, cmd.ShowZeroCounts)
	if err != nil {
		log.Fatalf("error iterating over bucket %s: %v", cmd.Id, err)
	}
}

func main() {
	var parsed args
	p := arg.MustParse(&parsed)

	if p.Subcommand() == nil {
		p.Fail("missing subcommand, see -h output for usage")
	}

	switch {
	case parsed.Select != nil:
		doSelect(parsed.Select)
	case parsed.Output != nil:
		doOutput(parsed.Output)
	default:
		p.WriteHelp(os.Stderr)
	}
}
