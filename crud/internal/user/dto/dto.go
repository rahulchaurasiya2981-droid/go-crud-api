package dto

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   *int   `json:"age"`
}

func (request *CreateUserRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)

	if request.Name == "" {
		return fmt.Errorf("name is required")
	}
	if request.Email == "" {
		return fmt.Errorf("email is required")
	}
	parsedEmail, err := mail.ParseAddress(request.Email)
	if err != nil || parsedEmail.Address != request.Email {
		return fmt.Errorf("email must be a valid email address")
	}
	if request.Age == nil {
		return fmt.Errorf("age is required")
	}
	if *request.Age < 0 {
		return fmt.Errorf("age must be zero or greater")
	}

	return nil
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
