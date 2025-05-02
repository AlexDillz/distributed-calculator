package storage

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"-"`
}

type Expression struct {
	ID         int     `json:"-"`
	UserID     int     `json:"-"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	Error      string  `json:"-"`
}
