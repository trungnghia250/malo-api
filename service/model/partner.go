package model

import (
	"mime/multipart"
	"time"
)

type PartnerConfig struct {
	ID            string    `json:"_id" bson:"_id"`
	StoreURL      string    `json:"store_url,omitempty" bson:"store_url,omitempty"`
	APIKey        string    `json:"api_key,omitempty" bson:"api_key,omitempty"`
	ModeSyncOrder string    `json:"mode_sync_order,omitempty" bson:"mode_sync_order,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt    time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy    string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	StoreName     string    `json:"store_name,omitempty" bson:"store_name,omitempty"`
	Logo          string    `json:"logo,omitempty" bson:"logo,omitempty"`

	From             string `json:"from,omitempty" bson:"from,omitempty"`
	MessageServiceID string `json:"message_service_id,omitempty" bson:"message_service_id,omitempty"`
	AccountID        string `json:"account_id,omitempty" bson:"account_id,omitempty"`
	AuthenToken      string `json:"authen_token,omitempty" bson:"authen_token,omitempty"`

	SendgridKey string `json:"sendgrid_key,omitempty" bson:"sendgrid_key,omitempty"`
	Address     string `json:"address,omitempty" bson:"address,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
}

type UploadRequest struct {
	File *multipart.FileHeader `form:"file"`
}
