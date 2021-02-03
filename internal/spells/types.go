package spells

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bergmannf/sr5foundryskillimport/internal/chummer"
	"github.com/bergmannf/sr5foundryskillimport/internal/foundry"
)

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

func (s Spell) ToFoundry() foundry.Spell {
	cs, err := FindChummerSpell(s.Name)
	if err != nil {
		fmt.Println(err)
	}
	drain, _ := strconv.Atoi(strings.Replace(s.Drain, "F", "", 1))
	damage, err := strconv.Atoi(cs.Damage)
	var spellType string
	if cs.Type == "P" {
		spellType = "physical"
	} else {
		spellType = "mana"
	}
	if err != nil {
		damage = 0
	}
	return foundry.Spell{
		Name: cs.Name,
		Type: "spell",
		Data: foundry.SpellData{
			Description: foundry.Description{
				Value:  s.Description,
				Chat:   "",
				Source: cs.Source + " " + strconv.Itoa(cs.Page),
			},
			Action: foundry.Action{
				Type:           cs.Type,
				Category:       cs.Category,
				Attribute:      "magic",
				Attribute2:     "",
				Skill:          "spellcasting",
				Spec:           false,
				Mod:            0,
				ModDescription: "",
				Limit:          foundry.LimitData{},
				Extended:       false,
				Damage: foundry.DamageData{
					Type:    foundry.DamageType{},
					Element: foundry.DamageElement{},
					Value:   damage,
					Base:    0,
					Ap:      foundry.Ap{},
				},
				Opposed:     foundry.OpposedData{},
				AltMod:      0,
				DicePoolMod: foundry.Empty{},
			},
			Drain:        drain,
			Category:     strings.ToLower(cs.Category),
			Type:         spellType,
			Range:        strings.ToLower(cs.Range),
			Duration:     strings.ToLower(cs.Duration),
			Combat:       foundry.CombatSpellData{},
			Detection:    foundry.DetectionSpellData{},
			Illusion:     foundry.IllusionSpellData{},
			Manipulation: foundry.ManipulationSpellData{},
		},
	}
}

func FindChummerSpell(spellname string) (chummer.Spell, error) {
	if spellname == "DISRUPT [FOCUS]" {
		spellname = "DISRUPT [OBJECT]"
	} else if spellname == "CAMOUFLAGE CHECK" {
		spellname = "CAMOUFLAGE"
	} else if spellname == "(CRITTER) FORM" {
		spellname = "[CRITTER] FORM"
	}
	for _, cs := range chummer.Load().Spells {
		if strings.ToLower(spellname) == strings.ToLower(cs.Name) {
			return cs, nil
		}
	}
	return chummer.Spell{}, errors.New(fmt.Sprintf("could not find spell for %s", spellname))
}
