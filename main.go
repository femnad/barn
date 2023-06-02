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

type listCmd struct {
	commonArgs
}

type outputCmd struct {
	commonArgs
	Id             string `arg:"positional"`
	ShowZeroCounts bool   `arg:"-z,--zero" help:"Show entries with zero counts"`
}

type purgeCmd struct {
	commonArgs
	Id string `arg:"positional,required"`
}

type selectCmd struct {
	commonArgs
	Id        string `arg:"-i,--id,required"`
	Selection string `arg:"positional"`
}

type args struct {
	List   *listCmd   `arg:"subcommand:list" help:"List existing buckets"`
	Output *outputCmd `arg:"subcommand:output" help:"Show stored entries for the given selection ID"`
	Purge  *purgeCmd  `arg:"subcommand:purge" help:"Purge given bucket"`
	Select *selectCmd `arg:"subcommand:select" help:"Select based on given choices and update counts"`
}

func showSelections(config, id string) {
	err := selection.Show(config, id)
	if err != nil {
		log.Fatalf("error getting selections for id %s: %v", id, err)
	}
}

func markSelection(config, id, choice string) {
	exitCode, err := selection.Mark(config, id, choice)
	if err != nil {
		log.Fatalf("error marking selection as %s for id %s: %v", choice, id, err)
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func doList(cmd *listCmd) {
	err := selection.ListBuckets(cmd.Config)
	if err != nil {
		log.Fatalf("error listing buckets: %v", err)
	}
}

func doOutput(cmd *outputCmd) {
	err := selection.Iterate(cmd.Config, cmd.Id, cmd.ShowZeroCounts)
	if err != nil {
		log.Fatalf("error iterating over bucket %s: %v", cmd.Id, err)
	}
}

func doPurge(cmd *purgeCmd) {
	err := selection.Purge(cmd.Config, cmd.Id)
	if err != nil {
		log.Fatalf("error purging bucket %s: %v", cmd.Id, err)
	}
}

func doSelect(cmd *selectCmd) {
	if cmd.Selection == "" {
		showSelections(cmd.Config, cmd.Id)
		return
	}

	markSelection(cmd.Config, cmd.Id, cmd.Selection)
}

func main() {
	var parsed args
	p := arg.MustParse(&parsed)

	if p.Subcommand() == nil {
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	switch {
	case parsed.List != nil:
		doList(parsed.List)
	case parsed.Purge != nil:
		doPurge(parsed.Purge)
	case parsed.Output != nil:
		doOutput(parsed.Output)
	case parsed.Select != nil:
		doSelect(parsed.Select)
	default:
		p.WriteHelp(os.Stderr)
	}
}
