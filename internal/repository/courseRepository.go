package repository

type Course struct {
	ID          string   `json:"id"db:"id"`
	Name        string   `json:"name"db:"name"`
	Description string   `json:"description"db:"description"`
	Price       float64  `json:"price"db:"price"`
	students    []string `json:"students"db:"students"`
	creator     string   `json:"creator"db:"creator"`
}
