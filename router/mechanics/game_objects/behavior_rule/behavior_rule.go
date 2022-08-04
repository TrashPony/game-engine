package behavior_rule

type BehaviorRules struct {
	Rules []*BehaviorRule `json:"rules"`
	Meta  *Meta           `json:"meta"`
	Key   string          `json:"key"`
}

type BehaviorRule struct {
	Action   string        `json:"action"`
	Meta     *Meta         `json:"meta"`
	PassRule *BehaviorRule `json:"access_rule"`
	StopRule *BehaviorRule `json:"stop_rule"`
}

type Meta struct {
	ID     int    `json:"ID"`
	Type   string `json:"type"`
	BaseID int    `json:"base_id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Radius int    `json:"radius"`
	Role   string `json:"role"`
}
