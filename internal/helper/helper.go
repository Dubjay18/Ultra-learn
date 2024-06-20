package helper

import "github.com/google/uuid"

func GenerateUserId() string {
	return uuid.New().String()
}
