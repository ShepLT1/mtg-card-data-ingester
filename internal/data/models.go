package data

import "time"

type Set struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `json:"name"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"`
	ReleasedAt  time.Time `json:"released_at"`
}

type Card struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `json:"name"`
	SetID        	int64  		`gorm:"uniqueIndex:idx_card_identity" json:"set"`
	CollectorNum 	string 		`gorm:"uniqueIndex:idx_card_identity" json:"collector_number"`
	PromoType    	string 		`gorm:"uniqueIndex:idx_card_identity" json:"promo_type"`
	Finish       	string 		`gorm:"uniqueIndex:idx_card_identity" json:"finish"`
	AltStyle			string		`json:"alt_style"`
	ImageURI     	string 		`json:"image_uri"`
}

type Listing struct {
	ID        	int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CardID    	int64     `gorm:"not null" json:"card_id"` // foreign key to Card
	Price     	float64   `json:"price"`
	CreatedAt 	time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedDate time.Time `gorm:"type:date;not null" json:"created_date"`
	// Optional preload field, not stored as a column
	Card      	*Card     `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE;" json:"card,omitempty"`
}