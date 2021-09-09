package queries

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	models "github.com/snapdocs/go-playground/pkg/models"
)

type Queries struct {
	db *gorm.DB
}

func New(adapter *gorm.DB) *Queries {
	return &Queries{db: adapter}
}

func (store *Queries) GetFirstWebhookEvent() (*models.WebhookEvent, error) {
	var webhookEvent models.WebhookEvent
	if err := store.db.Table("webhook_events").First(&webhookEvent).Error; err != nil {
		return nil, err
	}

	return &webhookEvent, nil
}

func (store *Queries) CreateWebhookEvent(webhookEvent *models.WebhookEvent) (*models.WebhookEvent, error) {
	result := store.db.Create(webhookEvent)
	if result.Error != nil {
		return nil, result.Error
	}

	return webhookEvent, nil
}

func (store *Queries) GetWebhookByID(id uuid.UUID) (*models.WebhookEvent, error) {
	var webhookEvent models.WebhookEvent
	if err := store.db.Table("webhook_events").Where("id=?", id).First(&webhookEvent).Error; err != nil {
		return nil, err
	}

	return &webhookEvent, nil
}
