package behavior_rule

type BehaviorRules struct {
	Rules []*BehaviorRule `json:"rules"`
	Meta  *Meta           `json:"meta"`
}

type BehaviorRule struct {
	Action   string        `json:"action"`
	Meta     *Meta         `json:"meta"`
	PassRule *BehaviorRule `json:"access_rule"`
	StopRule *BehaviorRule `json:"stop_rule"`
	Exit     bool          `json:"exit"`
}

type Meta struct {
}
