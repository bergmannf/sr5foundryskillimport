package internal
import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseSpell(t *testing.T) {
	acid, _ := goquery.NewDocumentFromReader(strings.NewReader(`<h2><span id="Acid_Stream,_Toxic_Wave"></span><span class="mw-headline" id="Acid_Stream.2C_Toxic_Wave">Acid Stream, Toxic
Wave</span></h2>
<table style="width:33%;">
<tbody><tr>
<th colspan="3" style="text-align:left;"><font size="3"><span style="color:#990000">ACID STREAM</span></font>
</th></tr>
<tr>
<th colspan="3" style="text-align:left;">(INDIRECT, ELEMENTAL)
<hr />
</th></tr>
<tr>
<td><b>Type:</b>P</td>
<td><b>Range:</b>LOS</td>
<td><b>Damage:</b>P
</td></tr>
<tr>
<td><b>Duration:</b>I</td>
<td><b>Drain:</b>F-3</td>
<td><b>Source:</b>Core
</td></tr></tbody></table>
<table style="width:33%;">
<tbody><tr>
<th colspan="3" style="text-align:left;"><font size="3"><span style="color:#990000">TOXIC WAVE</span></font>
</th></tr>
<tr>
<th colspan="3" style="text-align:left;">(INDIRECT, ELEMENTAL)
<hr />
</th></tr>
<tr>
<td><b>Type:</b>P</td>
<td><b>Range:</b>LOS(A)</td>
<td><b>Damage:</b>P
</td></tr>
<tr>
<td><b>Duration:</b>I</td>
<td><b>Drain:</b>F-1</td>
<td><b>Source:</b>Core
</td></tr></tbody></table>
<p>These spells create a powerful corrosive that sprays the target, causing terrible burns and eating away organic and metallic material—treat it as <b><a href="/Shadowrun/index.php/SR5:Combat_Rules:Elemental_Damage#Acid_Damage" title="SR5:Combat Rules:Elemental Damage">Acid Damage</a></b>, with appropriate effects on the affected area and any objects therein. The acid quickly evaporates, but the damage it inflicts remains. Acid Stream is a single-target spell, Toxic Wave is an area spell.
</p>`))
		type args struct {
			s *goquery.Selection
		}
		tests := []struct {
			name string
			args args
			want []Spell
		}{
			{
				name: "Acid Stream",
				args: args{ s: acid.Contents() },
				want: []Spell{
					{
						Name: "ACID STREAM",
						Effects: []string{"INDIRECT", "ELEMENTAL"},
						Category: "",
						Type: "P",
						Range: "LOS",
						Damage: "P",
						Duration: "I",
						Drain: "F-3",
						Source: "Core",
						Description: "These spells create a powerful corrosive that sprays the target, causing terrible burns and eating away organic and metallic material—treat it as Acid Damage, with appropriate effects on the affected area and any objects therein. The acid quickly evaporates, but the damage it inflicts remains. Acid Stream is a single-target spell, Toxic Wave is an area spell.",
					},
					{
						Name: "TOXIC WAVE",
						Effects: []string{"INDIRECT", "ELEMENTAL"},
						Category: "",
						Type: "P",
						Range: "LOS(A)",
						Damage: "P",
						Duration: "I",
						Drain: "F-1",
						Source: "Core",
						Description: "These spells create a powerful corrosive that sprays the target, causing terrible burns and eating away organic and metallic material—treat it as Acid Damage, with appropriate effects on the affected area and any objects therein. The acid quickly evaporates, but the damage it inflicts remains. Acid Stream is a single-target spell, Toxic Wave is an area spell.",
					},
				},
			},
		}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSpell(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSpell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSpells(t *testing.T) {
	type args struct {
		spells   string
		category string
	}
	tests := []struct {
		name string
		args args
		want []Spell
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSpells(tt.args.spells, tt.args.category); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSpells() = %v, want %v", got, tt.want)
			}
		})
	}
}
