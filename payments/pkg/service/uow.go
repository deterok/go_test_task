package service

import (
	"github.com/jinzhu/gorm"
)

type UOWPayments interface {
	Save() error
	Revert() error

	Accounts() AccountsRepository
	Operations() OperationsRepository
}

type UOWPaymentsFactory interface {
	Make() (UOWPayments, error)
}

type uowPayments struct {
	db     *gorm.DB
	accRep AccountsRepository
	opRep  OperationsRepository
}

func NewUOWPayments(db *gorm.DB, accRep AccountsRepository, opRep OperationsRepository) UOWPayments {
	return &uowPayments{
		db:     db,
		accRep: accRep,
		opRep:  opRep,
	}
}

func (u *uowPayments) Save() error {
	return u.db.Commit().Error
}

func (u *uowPayments) Revert() error {
	return u.db.Rollback().Error
}

func (u *uowPayments) Accounts() AccountsRepository {
	return u.accRep
}

func (u *uowPayments) Operations() OperationsRepository {
	return u.opRep
}

type uowPaymentsFactory struct {
	db *gorm.DB
}

func NewUOWPaymentsFactory(db *gorm.DB) UOWPaymentsFactory {
	return &uowPaymentsFactory{db}
}

func (f *uowPaymentsFactory) Make() (UOWPayments, error) {
	tx := f.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	return NewUOWPayments(
		tx,
		NewAccountsRepository(tx),
		NewOperationsRepository(tx),
	), nil
}
