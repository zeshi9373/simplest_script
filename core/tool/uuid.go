package tool

import (
	"github.com/google/uuid"
)

func Uuid() string {
	u, _ := uuid.NewRandom()

	return u.String()
}
