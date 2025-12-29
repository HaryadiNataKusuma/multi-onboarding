package model

// AddonRuleDetail represents an addon rule detail entity
type AddonRuleDetail struct {
	ID               int64   `json:"id"`
	AddonRuleID      int64   `json:"addon_rule_id"`
	StartCondition   string  `json:"start_condition"`
	EndCondition     string  `json:"end_condition"`
	ValueType        string  `json:"value_type"`
	Value            float64 `json:"value"`
	DurationRuleType string  `json:"duration_rule_type"`
	MinAdult         int     `json:"min_adult"`
	MaxAdult         int     `json:"max_adult"`
	MaxAge           int     `json:"max_age"`
	CreatedBy        int64   `json:"created_by"`
	CreatedAt        string  `json:"created_at"`
	UpdatedBy        *int64  `json:"updated_by,omitempty"`
	UpdatedAt        *string `json:"updated_at,omitempty"`
}

// AddonRuleDetailDraft represents a draft addon rule detail
type AddonRuleDetailDraft struct {
	ProductCode      string  `json:"product_code"`
	AddonCode        string  `json:"addon_code"`
	AddonName        string  `json:"addon_name"`        // For template only, not saved to DB
	ProductName      string  `json:"product_name"`     // For template only, not saved to DB
	AddonRuleID      string  `json:"addon_rule_id"`   // For template only, not saved to DB (will be looked up)
	StartCondition   string  `json:"start_condition"`
	EndCondition     string  `json:"end_condition"`
	ValueType        string  `json:"value_type"`
	Value            float64 `json:"value"`
	DurationRuleType string  `json:"duration_rule_type"`
	MinAdult         int     `json:"min_adult"`
	MaxAdult         int     `json:"max_adult"`
	MaxAge           int     `json:"max_age"`
}

