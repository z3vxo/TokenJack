package core

// LoginHandler
type LoginReq struct {
	Usernmae string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Access  string `json:"jwt"`
	Refresh string `json:"jwt_refresh"`
}
