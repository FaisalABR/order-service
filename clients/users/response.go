package clients

import "github.com/google/uuid"

type UserResponse struct {
	Code    string   `json:"code"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type UserData struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
}
