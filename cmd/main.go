package main

import (
	"flag"

	"github.com/bergmannf/sr5foundryskillimport/internal/chummer"
	"github.com/bergmannf/sr5foundryskillimport/internal/powers"
	"github.com/bergmannf/sr5foundryskillimport/internal/spells"
)

func main() {
	download := flag.Bool("download", false, "Download the spells")
	flag.Parse()
	if *download {
		spells.DownloadAll()
		powers.Download()
	}
	sps := spells.Load()
	spells.Save(sps)
	ps := powers.Load()
	powers.Save(ps)
	cs := chummer.Load()
	for _, sp := range sps {
		sp.ToFoundry(cs)
	}
}
