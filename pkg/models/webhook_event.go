package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type WebhookEvent struct {
	ID                     uuid.UUID
	WebhookEvent           string
	EmailProvider          string
	To                     string
	EmailProviderMessageId string
	Reason                 string
	Event                  string
	SdMessageId            uuid.UUID
	Timestamp              int
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *WebhookEvent) BeforeCreate(tx *gorm.DB) error {

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	base.ID = uuid
	return nil
}
