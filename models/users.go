package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}
