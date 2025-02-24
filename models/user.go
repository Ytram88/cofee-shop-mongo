package models

type User struct {
	UserID   string `json:"user_id" bson:"user_id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type RegisterUserPayload struct {
	UserID   string `json:"user_id" bson:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password"`
}

type UserLoginPayload struct {
	Email    string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
