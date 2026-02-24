package casbin

type AllocDomainRoleResourcesReq struct {
	Identifiers []string `json:"identifiers" binding:""`
}

type AllocRolesInDomainReq struct {
	Domain string   `json:"domain" binding:"required"`
	Roles  []string `json:"roles" binding:""`
}
