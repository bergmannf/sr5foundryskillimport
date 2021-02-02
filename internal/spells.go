package internal

import (
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
)


const basePath = "./data/%s"
const baseUrl = "http://adragon202.no-ip.org/Shadowrun/index.php/SR5:Spells:%s"
func getSpellTypes() []string {
	return []string{"Combat", "Detection", "Health", "Illusion", "Manipulation"}
}

type Spell struct {
	Name string;
	Effects []string;
	Category string;
	Type string;
	Range string;
	Damage string;
	Duration string;
	Drain string;
	Description string;
	Source string;
}

func DownloadAllSpells() {
	for _, category := range getSpellTypes() {
		DownloadSpells(category)
	}
}

// Download the HTML contain the spells of the given category
func DownloadSpells(category string) {
	fullUrl := fmt.Sprintf(baseUrl, category)
	res, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Could not retrieve URL: ", category)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Could not read response body.", category)
		return
	}
	ioutil.WriteFile(fmt.Sprintf(basePath, category), body, 0644)
}

// Parse a single section of HTML into spells.
// Each section can contain 1+ spells with the same description.
func ParseSpell(s *goquery.Selection) []Spell {
	spells := []Spell{}
	description := strings.TrimSpace(s.Find("p").Text())
	s.Find("table").Each(func(i int, table *goquery.Selection){
		var name, spellType, spellRange, damage, duration, drain, source string
		var effects []string
		sel := table.Find("th")
		for i := range sel.Nodes {
			if i == 0 {
				name = strings.TrimSpace(sel.Eq(i).Text())
			} else {
				tmp := sel.Eq(i).Text()
				tmp = strings.Replace(tmp, "(", "", 1)
				tmp = strings.Replace(tmp, ")", "", 1)
				tmp = strings.TrimSpace(tmp)
				split := strings.Split(tmp, ",")
				for _, e := range split {
					e = strings.TrimSpace(e)
					effects = append(effects, e)
				}
			}
		}
		table.Find("td").Each(func(i int, data *goquery.Selection) {
			text := data.Text()
			if strings.Contains(text, "Type:") {
				spellType = strings.TrimSpace(strings.Replace(text, "Type:", "", 1))
			}
			if strings.Contains(text, "Range:") {
				spellRange = strings.TrimSpace(strings.Replace(text, "Range:", "", 1))
			}
			if strings.Contains(text, "Damage:") {
				damage = strings.TrimSpace(strings.Replace(text, "Damage:", "", 1))
			}
			if strings.Contains(text, "Duration:") {
				duration = strings.TrimSpace(strings.Replace(text, "Duration:", "", 1))
			}
			if strings.Contains(text, "Drain:") {
				drain = strings.TrimSpace(strings.Replace(text, "Drain:", "", 1))
			}
			if strings.Contains(text, "Source:") {
				source = strings.TrimSpace(strings.Replace(text, "Source:", "", 1))
			}
		})
		spell := Spell {
			Name: name,
			Effects: effects,
			Type: spellType,
			Range: spellRange,
			Damage: damage,
			Description: description,
			Duration: duration,
			Drain: drain,
			Source: source,
		}
		spells = append(spells, spell)
	})
	return spells
}

// Parse all spells in the given HTML
func ParseSpells(spellHtml string, category string) []Spell {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(spellHtml))
	var spells []Spell
	if err != nil {
		fmt.Println("Problem loading ", category, err)
	}
	document.Find(".spell").Each(func (i int, s *goquery.Selection){
		xs := ParseSpell(s)
		spells = append(spells, xs...)
	})
	return spells
}

// Add a div element around every spell group to make it queryable in goquery
func TransformSpells(spells string) string {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(spells))
	if err != nil {
		fmt.Println("Problem loading ", err)
	}
	// Something like this should work: https://github.com/PuerkitoBio/goquery/issues/338
	document.Find(".mw-parser-output > h2").Each(func (i int, s *goquery.Selection){
		s.AddSelection(s.NextUntil("h2")).WrapAllHtml("<div class='spell'></div>")
	})
	html, _ := goquery.OuterHtml(document.Selection)
	return html
}

func LoadSpells() []Spell {
	var spells []Spell
	for _, category := range getSpellTypes() {
		spellText, err := ioutil.ReadFile(fmt.Sprintf(basePath, category))
		if err != nil {
			fmt.Println("Could not read file ", category, ": ", err)
		}
		transformedText := TransformSpells(string(spellText))
		spells = append(spells, ParseSpells(transformedText, category)...)
	}
	return spells
}

func SaveSpells(spells []Spell) {
	for _, category := range getSpellTypes() {
		path := fmt.Sprintf(basePath, category + ".json")
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			fmt.Println("Could not open file: ", path)
		}
		for _, sp := range spells {
			spj, err := json.Marshal(sp)
			if err != nil {
				fmt.Println("Cloud not marshal spell to json: ", err)
			}
			f.Write(spj)
			f.WriteString("\n")
		}
	}
}
