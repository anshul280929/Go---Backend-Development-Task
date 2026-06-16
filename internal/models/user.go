package models

import "time"

// CreateUserRequest is the payload for POST /users.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is the payload for PUT /users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UserResponse is the JSON response returned by the API.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  *int   `json:"age,omitempty"` // nil when omitted (create/update)
}

// CalculateAge computes the age in years from a date of birth to today.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()

	// If the birthday hasn't occurred yet this year, subtract one.
	if now.YearDay() < dob.YearDay() {
		age--
	}

	return age
}
