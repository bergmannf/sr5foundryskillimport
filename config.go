package sr5foundryskillimport

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	StoragePath        string
	SpellDownloadUrl   string
	ChummerDownloadUrl string
	SpellTypes         []string
}

func GetConfig() Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Could not find working directory, using relative path.")
		wd = "./data/"
	}
	return Config{
		StoragePath:        fmt.Sprintf("%s/%s", wd, "data"),
		SpellDownloadUrl:   "http://adragon202.no-ip.org/Shadowrun/index.php/SR5:Spells:%s",
		ChummerDownloadUrl: "https://raw.githubusercontent.com/chummer5a/chummer5a/master/Chummer/data/spells.xml",
		SpellTypes: []string{"Combat", "Detection", "Health", "Illusion", "Manipulation"},
	}
}
