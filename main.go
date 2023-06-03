package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/femnad/barn/selection"
)

const (
	name    = "barn"
	version = "0.1.0"
)

type commonArgs struct {
	Config string `arg:"-f,--file" default:"~/.config/barn/barn.yml" help:"Config file path"`
}

type bucketsCmd struct {
	commonArgs
}

type outputCmd struct {
	commonArgs
	Id             string `arg:"positional" help:"Selection Id"`
	ShowZeroCounts bool   `arg:"-z,--zero" help:"Show entries with zero counts"`
}

type purgeCmd struct {
	commonArgs
	Ids []string `arg:"positional,required"`
}

type truncateCmd struct {
	commonArgs
	Id    string   `arg:"-i,--id,required" help:"Selection ID"`
	Entry []string `arg:"positional,required" help:"Keys to truncate"`
}

type chooseCmd struct {
	commonArgs
	DontReverse bool   `arg:"-r,--dont-reverse" help:"Don't reverse output when listing choices"`
	Id          string `arg:"-i,--id,required" help:"Selection ID"`
	Selection   string `arg:"positional" help:"choice to mark as selected"`
}

type args struct {
	Buckets  *bucketsCmd  `arg:"subcommand:buckets" help:"List existing buckets"`
	Choose   *chooseCmd   `arg:"subcommand:choose" help:"Make a selection based on given choices and update counts"`
	Output   *outputCmd   `arg:"subcommand:output" help:"Show stored entries for the given selection ID"`
	Purge    *purgeCmd    `arg:"subcommand:purge" help:"Purge given bucket"`
	Truncate *truncateCmd `arg:"subcommand:truncate" help:"Truncate the desired keys for the given bucket"`
}

func (args) Version() string {
	return fmt.Sprintf("%s %s", name, version)
}

func showSelections(config, id string, reverse bool) {
	err := selection.Show(config, id, reverse)
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

func doList(cmd *bucketsCmd) {
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
	err := selection.Purge(cmd.Config, cmd.Ids)
	if err != nil {
		log.Fatalf("error purging buckets: %v", err)
	}
}

func doSelect(cmd *chooseCmd) {
	if cmd.Selection == "" {
		showSelections(cmd.Config, cmd.Id, !cmd.DontReverse)
		return
	}

	markSelection(cmd.Config, cmd.Id, cmd.Selection)
}

func doTruncate(cmd *truncateCmd) {
	err := selection.Truncate(cmd.Config, cmd.Id, cmd.Entry)
	if err != nil {
		log.Fatalf("error truncating keys from bucket %s: %v", cmd.Id, err)
	}
}

func main() {
	var parsed args
	p := arg.MustParse(&parsed)

	if p.Subcommand() == nil {
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	switch {
	case parsed.Buckets != nil:
		doList(parsed.Buckets)
	case parsed.Choose != nil:
		doSelect(parsed.Choose)
	case parsed.Output != nil:
		doOutput(parsed.Output)
	case parsed.Purge != nil:
		doPurge(parsed.Purge)
	case parsed.Truncate != nil:
		doTruncate(parsed.Truncate)
	default:
		p.WriteHelp(os.Stderr)
	}
}
