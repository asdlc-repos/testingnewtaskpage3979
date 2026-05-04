package models

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ManagerId string `json:"managerId,omitempty"`
}

type LeaveBalance struct {
	UserId      string  `json:"userId"`
	Entitlement float64 `json:"entitlement"`
	Used        float64 `json:"used"`
	Remaining   float64 `json:"remaining"`
}
