package service

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var ErrAccountNotFound = errors.Errorf("account not found")

const WorldAccountID = -1

type Account struct {
	ID       int64           `gorm:"primary_key" json:"id"`
	Name     string          `json:"name"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `sql:"type:decimal(20,8);" json:"amount"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

type AccountsRepository interface {
	Create(ctx context.Context, a *Account) (*Account, error)
	Update(ctx context.Context, a *Account) (*Account, error)
	Delete(ctx context.Context, id int64) error

	Get(ctx context.Context, id int64) (*Account, error)
	GetAll(ctx context.Context) ([]*Account, error)
}

// ─── IMPLEMENTATION ─────────────────────────────────────────────────────────────

type accountsRepository struct {
	db *gorm.DB
}

func NewAccountsRepository(db *gorm.DB) AccountsRepository {
	return &accountsRepository{db}
}

func (r *accountsRepository) Init() error {
	if err := r.db.AutoMigrate(&Account{}).Error; err != nil {
		return errors.Wrap(err, "migration failed")
	}
	return nil
}

func (r *accountsRepository) Create(ctx context.Context, a *Account) (*Account, error) {
	if err := r.db.Create(&a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (r *accountsRepository) Get(ctx context.Context, id int64) (*Account, error) {
	acc := Account{}
	if err := r.db.Find(&acc, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	return &acc, nil
}

func (r *accountsRepository) GetAll(ctx context.Context) ([]*Account, error) {
	a := []*Account{}

	if err := r.db.Find(&a).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	return a, nil
}

func (r *accountsRepository) Update(ctx context.Context, a *Account) (*Account, error) {
	if err := r.db.Save(a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (r *accountsRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.Delete(id).Error; err != nil {
		return err
	}

	return nil
}
