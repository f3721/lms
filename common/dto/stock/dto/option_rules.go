package dto

type RuleBool bool

func (r *RuleBool) True() {
	*r = true
}

func (r *RuleBool) False() {
	*r = false
}

type OptionRulesCommon struct {
	Update RuleBool `json:"update"`
	Delete RuleBool `json:"delete"`
	Detail RuleBool `json:"detail"`
}

type OptionRulesTransfer struct {
	Audit RuleBool `json:"Audit"`
	OptionRulesCommon
}

type OptionRulesEntry struct {
	Confirm RuleBool `json:"confirm"`
	Print   RuleBool `json:"print"`
	Detail  RuleBool `json:"detail"`
	Audit   RuleBool `json:"Audit"`
}

type OptionRulesOutbound struct {
	Confirm RuleBool `json:"confirm"`
	Print   RuleBool `json:"print"`
	Detail  RuleBool `json:"detail"`
}

type OptionRulesControl struct {
	Audit RuleBool `json:"Audit"`
	OptionRulesCommon
}
