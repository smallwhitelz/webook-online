package web

type LoginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type SignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginJWTReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EditReq struct {
	// 注意，其它字段，尤其是密码、邮箱和手机，
	// 修改都要通过别的手段
	// 邮箱和手机都要验证
	// 密码更加不用多说了
	Nickname string `json:"nickname"`
	// YYYY-MM-DD
	Birthday    string `json:"birthday"`
	Description string `json:"description"`
}
