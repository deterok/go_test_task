package service

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func InitModels(db *gorm.DB) error {
	err := db.AutoMigrate(
		Account{},
		Operation{},
		Transaction{},
	).Error

	if err != nil {
		return errors.Wrap(err, "migration failed")
	}

	return nil
}
