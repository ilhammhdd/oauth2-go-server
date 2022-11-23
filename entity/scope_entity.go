package entity

type Scope struct {
	Scope          string  `json:"scope"`
	ParentScopesID *uint64 `json:"parent_scopes_id"`
	Permission     uint8   `json:"permission"`
	ClientsID      *uint64 `json:"clients_id"`
	TableTemplateCols
}

type ScopeWithRel struct {
	Scope          *Scope          `json:"scope"`
	ParentScope    *Scope          `json:"parent_scope"`
	ChildrenScopes []*ScopeWithRel `json:"children_scopes"`
	AuthzCodes     []*AuthzCode    `json:"authz_codes"`
	Clients        *Client         `json:"clients"`
}
