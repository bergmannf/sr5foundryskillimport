package internal

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const basePath = "./data/%s"
const baseUrl = "http://adragon202.no-ip.org/Shadowrun/index.php/SR5:Spells:%s"

var chummerSpells ChummerSpells

func getSpellTypes() []string {
	return []string{"Combat", "Detection", "Health", "Illusion", "Manipulation"}
}

type Spell struct {
	Name        string
	Effects     []string
	Category    string
	Type        string
	Range       string
	Damage      string
	Duration    string
	Drain       string
	Description string
	Source      string
}

type Empty struct{}

type FoundryDescription struct {
	Value  string `json:"value"`
	Chat   string `json:"chat"`
	Source string `json:"source"`
}

type LimitData struct {
	Value     int    `json:"value"`
	Base      int    `json:"base"`
	Attribute string `json:"attribute"`
	Mod       Empty  `json:"mod"`
}

type DamageType struct {
	Value string `json:"value"`
	Base  string `json:"base"`
}

type DamageElement struct {
	Value string `json:"value"`
	Base  string `json:"base"`
}

type Ap struct {
	Value int   `json:"value"`
	Base  int   `json:"base"`
	Mod   Empty `json:"mod"`
}

type DamageData struct {
	Type      DamageType    `json:"type"`
	Element   DamageElement `json:"element"`
	Value     int           `json:"value"`
	Base      int           `json:"base"`
	Ap        Ap            `json:"ap"`
	Attribute string        `json:"attribute"`
	Mod       Empty         `json:"mod"`
}

type OpposedData struct {
	Type        string `json:"type"`
	Attribute   string `json:"attribute"`
	Attribute2  string `json:"attribute2"`
	Skill       string `json:"skill"`
	Mod         int    `json:"mod"`
	Description string `json:"description"`
}

type FoundryAction struct {
	Type           string      `json:"type"`
	Category       string      `json:"category"`
	Attribute      string      `json:"attribute"`
	Attribute2     string      `json:"attribute2"`
	Skill          string      `json:"skill"`
	Spec           bool        `json:"spec"`
	Mod            int         `json:"mod"`
	ModDescription string      `json:"mod_description"`
	Limit          LimitData   `json:"limit"`
	Extended       bool        `json:"extended"`
	Damage         DamageData  `json:"damage"`
	Opposed        OpposedData `json:"opposed"`
	AltMod         int         `json:"alt_mod"`
	DicePoolMod    Empty       `json:"dice_pool_mod"`
}

type CombatSpellData struct {
	Type string `json:"type"`
}

type DetectionSpellData struct {
	Passive  bool   `json:"passive"`
	Type     string `json:"type"`
	Extended bool   `json:"extended"`
}

type IllusionSpellData struct {
	Type  string `json:"type"`
	Sense string `json:"sense"`
}

type ManipulationSpellData struct {
	Damaging      bool `json:"damaging"`
	Mental        bool `json:"mental"`
	Environmental bool `json:"environmental"`
	Physical      bool `json:"physical"`
}

type FoundrySpellData struct {
	Description  FoundryDescription    `json:"description"`
	Action       FoundryAction         `json:"action"`
	Drain        uint                  `json:"drain"`
	Category     string                `json:"category"`
	Type         string                `json:"type"`
	Range        string                `json:"range"`
	Duration     string                `json:"duration"`
	Combat       CombatSpellData       `json:"combat"`
	Detection    DetectionSpellData    `json:"detection"`
	Illusion     IllusionSpellData     `json:"illusion"`
	Manipulation ManipulationSpellData `json:"manipulation"`
}

type FoundrySpell struct {
	Name string           `json:"name"`
	Type string           `json:"type"`
	Data FoundrySpellData `json:"data"`
}

type ChummerSpell struct {
	Id         string `xml:"id"`
	Name       string `xml:"name"`
	Page       int    `xml:"page"`
	Source     string `xml:"source"`
	Category   string `xml:"category"`
	Damage     string `xml:"damage"`
	Descriptor string `xml:"descriptor"`
	Duration   string `xml:"duration"`
	Dv         string `xml:"dv"`
	Range      string `xml:"range"`
	Type       string `xml:"type"`
}

type ChummerSpells struct {
	Spells  []ChummerSpell `xml:"spells>spell"`
}

func (s Spell) ToFoundry() FoundrySpell {
	_, err := FindChummerSpell(s.Name)
	if err != nil {
		fmt.Println(err)
	}
	return FoundrySpell{}
}

func FindChummerSpell(spellname string) (ChummerSpell, error) {
	if spellname == "DISRUPT [FOCUS]" {
		spellname = "DISRUPT [OBJECT]"
	} else if spellname == "CAMOUFLAGE CHECK" {
		spellname = "CAMOUFLAGE"
	} else if spellname == "(CRITTER) FORM" {
		spellname = "[CRITTER] FORM"
	}
	for _, cs := range LoadChummer().Spells {
		if strings.ToLower(spellname) == strings.ToLower(cs.Name) {
			return cs, nil
		}
	}
	return ChummerSpell{}, errors.New(fmt.Sprintf("could not find spell for %s", spellname))
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
func ParseSpell(s *goquery.Selection, c string) []Spell {
	spells := []Spell{}
	description := strings.TrimSpace(s.Find("p").Text())
	s.Find("table").Each(func(i int, table *goquery.Selection) {
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
		spell := Spell{
			Name:        name,
			Effects:     effects,
			Type:        spellType,
			Range:       spellRange,
			Damage:      damage,
			Description: description,
			Duration:    duration,
			Drain:       drain,
			Source:      source,
			Category:    c,
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
	document.Find(".spell").Each(func(i int, s *goquery.Selection) {
		xs := ParseSpell(s, category)
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
	// Something like this should work: https://github.com/PuerkitoBio/goquery/issues/338ll
	document.Find(".mw-parser-output > h2").Each(func(i int, s *goquery.Selection) {
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

func DownloadChummer() {
	url := "https://raw.githubusercontent.com/chummer5a/chummer5a/master/Chummer/data/spells.xml"
	f, err := os.Stat("./data/spells.xml")
	if os.IsNotExist(err) || !f.IsDir() {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("Could not retrieve URL: ", err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Could not read response body.", err)
			return
		}
		ioutil.WriteFile("./data/spells.xml", body, 0644)
	}
}

func LoadChummer() ChummerSpells {
	DownloadChummer()
	if len(chummerSpells.Spells) > 0 {
		return chummerSpells
	}
	xmlFile, _ := os.Open("./data/spells.xml")
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &chummerSpells)
	return chummerSpells
}

func SaveSpells(spells []Spell) {
	for _, category := range getSpellTypes() {
		path := fmt.Sprintf(basePath, category+".json")
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
