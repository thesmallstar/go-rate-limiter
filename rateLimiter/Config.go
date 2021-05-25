package rateLimiter

//name user is used to keep it less specific to BVS and more general
type ConfigFile struct {
	Users []User `json:"users"`
}

type User struct {
	UserID string      `json:"userID"`
	Limits []URL_Limit `json:"limits"`
}

type URL_Limit struct {
	Route   string `json:"route"`
	TypeReq string `json:"type"`
	Rate    int    `json:"rate"`
	Limit   int    `json:"limit"`
}
