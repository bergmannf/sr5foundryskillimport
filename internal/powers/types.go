package powers

import (
	"github.com/bergmannf/sr5foundryskillimport/internal/foundry"
)

type CritterPower struct {
	Name string
	Type string
	Range string
	Source string
	Action string
	Description string
	Duration string
}

func (c *CritterPower) ToFoundry() foundry.CritterPower {
	return foundry.CritterPower{}
}
