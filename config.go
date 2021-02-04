package sr5foundryskillimport

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	ChummerDownloadUrl      string
	CritterPowerDownloadUrl string
	StoragePath             string
	SpellDownloadUrl        string
	SpellTypes              []string
}

func GetConfig() Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Could not find working directory, using relative path.")
		wd = "./data/"
	}
	return Config{
		ChummerDownloadUrl: "https://raw.githubusercontent.com/chummer5a/chummer5a/master/Chummer/data/spells.xml",
		CritterPowerDownloadUrl: "http://adragon202.no-ip.org/Shadowrun/index.php/SR5:Critter_Powers",
		StoragePath:        fmt.Sprintf("%s/%s", wd, "data"),
		SpellDownloadUrl:   "http://adragon202.no-ip.org/Shadowrun/index.php/SR5:Spells:%s",
		SpellTypes:         []string{"Combat", "Detection", "Health", "Illusion", "Manipulation"},
	}
}
