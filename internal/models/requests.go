package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Request struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Description string
	Certificate string
	Status      string
	Type        string
	CreatedAt   time.Time
}