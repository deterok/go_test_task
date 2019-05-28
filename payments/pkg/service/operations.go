package service

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type OperationType int

const (
	OperationTypeDeposit OperationType = iota
	OperationTypeTransfer
)

func (t OperationType) String() string {
	switch t {
	case OperationTypeDeposit:
		return "Deposit"
	case OperationTypeTransfer:
		return "Transfer"
	}
	return ""
}

var ErrOperationNotFound = errors.Errorf("operation not found")

// Operation is a transactions grouping object
type Operation struct {
	gorm.Model
	Participants pq.Int64Array `gorm:"type:integer[]"`
	Type         OperationType
	Transactions []Transaction
}

// Transaction is an atomic unit account changes
type Transaction struct {
	gorm.Model
	OperationID uint
	From        int64
	To          int64

	Currency string
	Amount   decimal.Decimal `sql:"type:decimal(20,8);"`
}

type OperationsRepository interface {
	Create(ctx context.Context, o *Operation) (*Operation, error)

	Get(ctx context.Context, id int64) (*Operation, error)
	GetByAccID(ctx context.Context, id int64) ([]*Operation, error)
	GetAll(ctx context.Context) ([]*Operation, error)
}

func InitModel(db gorm.DB) error {
	if err := db.AutoMigrate(&Account{}).Error; err != nil {
		return errors.Wrap(err, "migration failed")
	}
	return nil
}

type operationsRepository struct {
	db *gorm.DB
}

func NewOperationsRepository(db *gorm.DB) OperationsRepository {
	return &operationsRepository{db}
}

func (r *operationsRepository) Init() error {
	if err := r.db.AutoMigrate(&Account{}).Error; err != nil {
		return errors.Wrap(err, "migration failed")
	}
	return nil
}

func (r *operationsRepository) Create(ctx context.Context, o *Operation) (*Operation, error) {
	if err := r.db.Create(&o).Error; err != nil {
		return nil, err
	}

	return o, nil
}

func (r *operationsRepository) Get(ctx context.Context, id int64) (*Operation, error) {
	op := Operation{}
	if err := r.db.Set("gorm:auto_preload", true).Find(&op, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOperationNotFound
		}

		return nil, err
	}

	return &op, nil
}

func (r *operationsRepository) GetByAccID(ctx context.Context, id int64) ([]*Operation, error) {
	o := []*Operation{}

	req := r.db.Set("gorm:auto_preload", true).Where("(?) = ANY(participants)", id)

	if err := req.Find(&o).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrOperationNotFound
		}

		return nil, err
	}

	return o, nil
}

func (r *operationsRepository) GetAll(ctx context.Context) ([]*Operation, error) {
	ops := []*Operation{}
	if err := r.db.Set("gorm:auto_preload", true).Find(&ops).Error; err != nil {
		return nil, err
	}
	return ops, nil
}
