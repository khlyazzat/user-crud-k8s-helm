package utils

import (
	"github.com/google/uuid"
)

func UUIDMustParse(input string) uuid.UUID {
	parsedUUID, err := uuid.Parse(input)
	if err != nil {
		return uuid.Nil
	}
	return parsedUUID
}
