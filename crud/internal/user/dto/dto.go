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

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

func (request *UpdateUserRequest) Validate() error {
	if request.Name == nil && request.Email == nil && request.Age == nil {
		return fmt.Errorf("at least one field must be provided")
	}
	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			return fmt.Errorf("name cannot be empty")
		}
		request.Name = &name
	}
	if request.Email != nil {
		email := strings.TrimSpace(*request.Email)
		if email == "" {
			return fmt.Errorf("email cannot be empty")
		}
		parsedEmail, err := mail.ParseAddress(email)
		if err != nil || parsedEmail.Address != email {
			return fmt.Errorf("email must be a valid email address")
		}
		request.Email = &email
	}
	if request.Age != nil && *request.Age < 0 {
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
