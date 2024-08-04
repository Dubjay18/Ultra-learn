package models

type User struct {
	ID        string `json:"id" db:"id"`
	Avatar    string `json:"avatar" db:"avatar"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Role      int    `json:"role" db:"role"`
	Provider  string `json:"provider,omitempty" db:"provider,omitempty"`
}
