package conditions

import (
	"github.com/ashbeelghouri/jsonschematics/utils"
)

type Conditions struct {
	ConditionFns map[string]Condition
	Logger       utils.Logger
}

type Condition func(map[string]interface{}, map[string]interface{}) bool

func (c *Conditions) RegisterCondition(name string, fn Condition) {
	c.Logger.DEBUG("registering condition:", name)
	if c.ConditionFns == nil {
		c.ConditionFns = make(map[string]Condition)
	}
	c.ConditionFns[name] = fn
}

func (c *Conditions) BasicConditions() {
	c.RegisterCondition("FieldIsProvided", FieldIsProvided)
}
