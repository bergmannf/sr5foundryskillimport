package chummer

import (
	"fmt"
	"log"
	"encoding/xml"
	"os"
	"io/ioutil"
	"net/http"
	"github.com/bergmannf/sr5foundryskillimport"
)

var spells Spells

func Download() {
	config := sr5foundryskillimport.GetConfig()
	filePath := fmt.Sprintf("%s/%s", config.StoragePath, "spells.xml")
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Println("Downloading spells from chummer.")
		res, err := http.Get(config.ChummerDownloadUrl)
		if err != nil {
			log.Println("Could not retrieve URL: ", err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("Could not read response body.", err)
			return
		}
		ioutil.WriteFile(filePath, body, 0644)
	}
}

func Load() Spells {
	Download()
	if len(spells.Spells) > 0 {
		return spells
	}
	xmlFile, _ := os.Open("./data/spells.xml")
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &spells)
	return spells
}
