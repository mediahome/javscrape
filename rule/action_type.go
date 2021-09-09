package rule

//ActionType action type
type ActionType string

const (
	ActionTypeNone        ActionType = "none"
	ActionTypeAction                 = "action"
	ActionTypeActionGroup            = "action_group"
)
