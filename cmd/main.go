package main

import (
	"flag"

	"github.com/bergmannf/sr5foundryskillimport/internal/chummer"
	"github.com/bergmannf/sr5foundryskillimport/internal/spells"
)

func main() {
	download := flag.Bool("download", false, "Download the spells")
	flag.Parse()
	if *download {
		spells.DownloadAllSpells()
	}
	sps := spells.Load()
	spells.Save(sps)
	cs := chummer.Load()
	for _, sp := range sps {
		sp.ToFoundry(cs)
	}
}
