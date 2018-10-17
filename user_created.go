package event

type UserCreatedPayload struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	HashPassword string `json:"hash_password"`
}
