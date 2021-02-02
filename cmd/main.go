package main

import (
	"flag"
	"github.com/bergmannf/sr5foundryskillimport/internal"
)

func main() {
	download := flag.Bool("download", false, "Download the spells")
	flag.Parse()
	if *download {
		internal.DownloadAllSpells()
	}
	spells := internal.LoadSpells()
	internal.SaveSpells(spells)
}
