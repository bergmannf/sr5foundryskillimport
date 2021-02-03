package chummer

type Spell struct {
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

type Spells struct {
	Spells []Spell `xml:"spells>spell"`
}
