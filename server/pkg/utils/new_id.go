package utils

import (
	"log"

	"github.com/google/uuid"
)

func NewId() string {
	id,err := uuid.NewV7()
	if err != nil {
		log.Printf("Lỗi khi sinh UUIDv7: %v\n", err)
	}
	return id.String()
}