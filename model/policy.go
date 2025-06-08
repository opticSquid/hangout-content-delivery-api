package model

// PolicyStatement defines a single statement in the custom policy
type PolicyStatement struct {
	Resource  string          `json:"Resource"`
	Condition PolicyCondition `json:"Condition"`
}

// PolicyCondition defines the conditions for the policy
type PolicyCondition struct {
	DateLessThan    map[string]int64  `json:"DateLessThan,omitempty"`
	DateGreaterThan map[string]int64  `json:"DateGreaterThan,omitempty"`
	IPAddress       map[string]string `json:"IpAddress,omitempty"`
}

// Policy defines the structure for a custom CloudFront policy
type Policy struct {
	Statement []PolicyStatement `json:"Statement"`
}
