package kvrepo

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationKV struct { // TODO добавить TTL
	ApplicationID uuid.UUID `gorm:"not null;primaryKey" json:"id"`
	Key           string    `gorm:"not null;primaryKey" json:"key"`

	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updatedAt"`

	Value any `gorm:"not null;serializer:json" json:"value"`
}
