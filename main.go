package main

import (
	"log"

	"github.com/alexflint/go-arg"

	"github.com/femnad/barn/selection"
)

type args struct {
	Id        string `arg:"-i,--id,required"`
	Config    string `arg:"-f,--file" default:"~/.config/barn/barn.yml" help:"Config file path"`
	Selection string `arg:"positional"`
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

func main() {
	var parsed args
	arg.MustParse(&parsed)

	if parsed.Selection == "" {
		showSelections(parsed.Config, parsed.Id)
		return
	}

	markSelection(parsed.Config, parsed.Id, parsed.Selection)
}
