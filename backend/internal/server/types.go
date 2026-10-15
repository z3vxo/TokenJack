package server

import "github.com/z3vxo/TokenJack/internal/db"

// LoginHandler
type LoginReq struct {
	Usernmae string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Access  string `json:"jwt"`
	Refresh string `json:"jwt_refresh"`
}

// campaign new
type CampaingNewReq struct {
	TenantId string         `json:"tenant_id"`
	Name     string         `json:"name"`
	Resource ResourceScopes `json:"resources"`
}

type ResourceScopes struct {
	Name   string     `json:"name"`
	Scopes []db.Scope `json:"scopes"`
}

// campaigns list
type CampaignsListResp struct {
	Total     int           `json:"total"`
	Campaigns []db.Campaign `json:"campaigns"`
}
