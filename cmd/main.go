package main

import (
	"flag"
	"github.com/bergmannf/sr5foundryskillimport/internal/spells"
)

func main() {
	download := flag.Bool("download", false, "Download the spells")
	flag.Parse()
	if *download {
		DownloadAllSpells()
	}
	spells := LoadSpells()
	SaveSpells(spells)
}
