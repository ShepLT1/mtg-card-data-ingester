package data

import (
	"context"

	"gorm.io/gorm/clause"
)

// SaveSetUpsert inserts a set if Code doesn't exist, or updates Name/ReleasedAt if it does.
func SaveSetUpsert(ctx context.Context, s *Set) error {
	return DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}}, // unique key
		UpdateAll: true,                             // update all fields if conflict
	}).Create(s).Error
}

// SaveCardUpsert inserts a card if it doesn't exist, or updates Name/ImageURI if it does.
func SaveCardUpsert(ctx context.Context, c *Card) error {
	return DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "set_id"},
			{Name: "collector_num"},
			{Name: "promo_type"},
			{Name: "finish"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"name", "image_uri"}), // fields to update on conflict
	}).Create(c).Error
}

func SaveListing(ctx context.Context, l *Listing) error {
	return DB.WithContext(ctx).Create(l).Error
}
