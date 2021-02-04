package powers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/bergmannf/sr5foundryskillimport"
	"os"
	"encoding/json"
)

func Download() {
	log.Println("Download critter powers from wiki.")
	config := sr5foundryskillimport.GetConfig()
	res, err := http.Get(config.CritterPowerDownloadUrl)
	if err != nil {
		log.Println("Could not download powers: ", err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read response body")
		return
	}
	ioutil.WriteFile(fmt.Sprintf("%s/%s", config.StoragePath, "powers"), body, 0644)
}

func Load() []CritterPower {
	var powers []CritterPower
	config := sr5foundryskillimport.GetConfig()
	powerText, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", config.StoragePath, "powers"))
	if err != nil {
		fmt.Println("Could not read file ", err)
	}
	transformedText := transform(string(powerText))
	powers = ParsePowers(transformedText)
	return powers
}

func transform(powers string) string {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(powers))
	if err != nil {
		fmt.Println("Problem loading ", err)
	}
	// Something like this should work: https://github.com/PuerkitoBio/goquery/issues/338ll
	document.Find(".mw-parser-output > h2").Each(func(i int, s *goquery.Selection) {
		s.AddSelection(s.NextUntil("h2")).WrapAllHtml("<div class='power'></div>")
	})
	html, _ := goquery.OuterHtml(document.Selection)
	return html
}

func ParsePowers(powersHtml string) []CritterPower {
	var powers []CritterPower
	document, err := goquery.NewDocumentFromReader(strings.NewReader(powersHtml))
	if err != nil {
		log.Println("Could not construct reader: ", err)
		return powers
	}
	document.Find(".power").Each(func(i int, s *goquery.Selection) {
		xs := Parse(s)
		powers = append(powers, xs)
	})
	return powers
}

func Parse(s *goquery.Selection) CritterPower {
	var power CritterPower
	description := strings.TrimSpace(s.Find("p").Text())
	var name, action, powerType, powerRange, duration, source string
	s.Find("h2").Each(func(i int, data *goquery.Selection) {
		name = strings.TrimSpace(data.Text())
	})
	s.Find("table").Each(func(i int, table *goquery.Selection) {
		table.Find("th").Each(func(i int, data *goquery.Selection) {
			text := data.Text()
			if strings.Contains(text, "Type:") {
				powerType = strings.TrimSpace(data.NextUntil("th").Text())
			}
			if strings.Contains(text, "Range:") {
				powerRange = strings.TrimSpace(data.NextUntil("th").Text())
			}
			if strings.Contains(text, "Duration:") {
				duration = strings.TrimSpace(data.NextUntil("th").Text())
			}
			if strings.Contains(text, "Source:") {
				source = strings.TrimSpace(data.NextUntil("th").Text())
			}
			if strings.Contains(text, "Action:") {
				action = strings.TrimSpace(data.NextUntil("th").Text())
			}
		})
		power = CritterPower{
			Name:        name,
			Type:        powerType,
			Range:       powerRange,
			Source:      source,
			Action:      action,
			Description: description,
			Duration:    duration,
		}
	})
	return power
}

func Save(powers []CritterPower) {
	config := sr5foundryskillimport.GetConfig()
	path := fmt.Sprintf("%s/%s.%s", config.StoragePath, "powers", "json")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("Could not open file: ", path)
	}
	for _, sp := range powers {
		spj, err := json.Marshal(sp)
		if err != nil {
			fmt.Println("Cloud not marshal powers to json: ", err)
		}
		f.Write(spj)
		f.WriteString("\n")
	}
}
