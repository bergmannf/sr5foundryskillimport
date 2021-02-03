package foundry

type Empty struct{}

type Description struct {
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

type Action struct {
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

type SpellData struct {
	Description  Description           `json:"description"`
	Action       Action                `json:"action"`
	Drain        int                   `json:"drain"`
	Category     string                `json:"category"`
	Type         string                `json:"type"`
	Range        string                `json:"range"`
	Duration     string                `json:"duration"`
	Combat       CombatSpellData       `json:"combat"`
	Detection    DetectionSpellData    `json:"detection"`
	Illusion     IllusionSpellData     `json:"illusion"`
	Manipulation ManipulationSpellData `json:"manipulation"`
}

type Spell struct {
	Name string           `json:"name"`
	Type string           `json:"type"`
	Data SpellData `json:"data"`
}
