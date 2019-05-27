package service

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/jinzhu/gorm"
)

//
// ──────────────────────────────────────────────────────────────────────────────────────────── I ──────────
//   :::::: D U M M Y   C L A S S E S   D E C L A R A T I O N S : :  :   :    :     :        :          :
// ──────────────────────────────────────────────────────────────────────────────────────────────────────
//

// ─── DUMMY LOCK ─────────────────────────────────────────────────────────────────

type dummyLock struct{}

func (dummyLock) Lock(key string) error {
	return nil
}

func (dummyLock) Unlock(key string) error {
	return nil
}

// ─── INMEMORY ACCOUNTS REPOSITORY ──────────────────────────────────────────────────

type inmemAccountsRepository struct {
	accsByID map[int64]*Account
}

func newInmemAccountsRepository() AccountsRepository {
	return &inmemAccountsRepository{make(map[int64]*Account)}
}

func (r *inmemAccountsRepository) randID() int64 {
	for {
		newID := rand.Int63n(40000)
		if _, ok := r.accsByID[newID]; !ok {
			return newID
		}
	}
}

func (r *inmemAccountsRepository) Create(ctx context.Context, a *Account) (*Account, error) {
	a.ID = r.randID()
	r.accsByID[a.ID] = a
	return a, nil
}

func (r *inmemAccountsRepository) Update(ctx context.Context, a *Account) (*Account, error) {
	r.accsByID[a.ID] = a
	return a, nil
}

func (r *inmemAccountsRepository) Delete(ctx context.Context, id int64) error {
	delete(r.accsByID, id)
	return nil
}

func (r *inmemAccountsRepository) Get(ctx context.Context, id int64) (*Account, error) {
	if a, ok := r.accsByID[id]; ok {
		return a, nil
	}

	return nil, ErrAccountNotFound
}

func (r *inmemAccountsRepository) GetAll(ctx context.Context) ([]*Account, error) {
	result := make([]*Account, len(r.accsByID))
	i := 0
	for _, a := range r.accsByID {
		result[i] = a
		i++
	}

	return result, nil
}

// ─── INMEMORY OPERATIONS REPOSITORY ─────────────────────────────────────────────

type inmemOperationsRepository struct {
	opsByID map[int64]*Operation
}

func newInmemOperationsRepository() OperationsRepository {
	return &inmemOperationsRepository{
		opsByID: make(map[int64]*Operation),
	}
}

func (r *inmemOperationsRepository) randID() int64 {
	for {
		newID := rand.Int63n(40000)
		if _, ok := r.opsByID[newID]; !ok {
			return newID
		}
	}
}

func (r *inmemOperationsRepository) Create(ctx context.Context, o *Operation) (*Operation, error) {
	o.ID = uint(r.randID())
	r.opsByID[int64(o.ID)] = o
	return o, nil
}

func (r *inmemOperationsRepository) Get(ctx context.Context, id int64) (*Operation, error) {
	if a, ok := r.opsByID[id]; ok {
		return a, nil
	}

	return nil, ErrAccountNotFound
}

func (r *inmemOperationsRepository) GetByAccID(ctx context.Context, id int64) ([]*Operation, error) {
	result := []*Operation{}

	for _, o := range r.opsByID {
		for _, accID := range o.Participants {
			if accID == id {
				result = append(result, o)
			}
		}
	}

	return result, nil
}

func (r *inmemOperationsRepository) GetAll(ctx context.Context) ([]*Operation, error) {
	result := make([]*Operation, len(r.opsByID))
	i := 0
	for _, a := range r.opsByID {
		result[i] = a
		i++
	}

	return result, nil
}

// ─── DUMMY UNIT OF WORK ──────────────────────────────────────────────────────

type dummyUOW struct {
	accs AccountsRepository
	ops  OperationsRepository
}

func newDummyUOW() UOWPayments {
	return &dummyUOW{
		accs: newInmemAccountsRepository(),
		ops:  newInmemOperationsRepository(),
	}
}

func (dummyUOW) Save() error {
	return nil
}

func (dummyUOW) Revert() error {
	return nil
}

func (u *dummyUOW) Accounts() AccountsRepository {
	return u.accs
}

func (u *dummyUOW) Operations() OperationsRepository {
	return u.ops
}

type dummyUOWFactory struct {
}

func newDummyUOWFactory() UOWPaymentsFactory {
	return &dummyUOWFactory{}
}

func (dummyUOWFactory) Make() (UOWPayments, error) {
	return newDummyUOW(), nil
}

// ─── UTILS ──────────────────────────────────────────────────────────────────────
func getDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=postgres dbname=testdb sslmode=disable user=postgres")
	if err != nil {
		panic(err)
	}

	if err := InitModels(db); err != nil {
		panic(err)
	}
	db.Exec("DELETE FROM accounts;")
	db.Exec("DELETE FROM operations;")
	db.Exec("DELETE FROM transactions;")

	return db
}

//
// ────────────────────────────────────────────────── II ──────────
//   :::::: T E S T S : :  :   :    :     :        :          :
// ────────────────────────────────────────────────────────────
//

// ─── CREATE ACCOUNTS ────────────────────────────────────────────────────────────

func Test_basicPaymentsService_CreateAccount(t *testing.T) {
	type args struct {
		ctx      context.Context
		name     string
		currency string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple create",
			args: args{
				name:     "test",
				currency: "USD",
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			db := getDB()
			defer db.Close()
			s := &basicPaymentsService{
				lock: &dummyLock{},
				uow:  NewUOWPaymentsFactory(db),
			}

			got, err := s.CreateAccount(tt.args.ctx, tt.args.name, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("basicPaymentsService.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotEqual(t, 0, got.ID, "basicPaymentsService.CreateAccount() error: ID is 0")
			assert.Equal(t, tt.args.currency, got.Currency, "basicPaymentsService.CreateAccount() error: currencies aren't equal")
			assert.Equal(t, tt.args.name, got.Name, "basicPaymentsService.CreateAccount() error: names aren't equal")

		})
	}
}

// ─── GET ACCOUNTS ───────────────────────────────────────────────────────────────

func Test_basicPaymentsService_GetAccount(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple getting",
			args: args{
				id: 1,
			},
		},
		{
			name: "another simple getting",
			args: args{
				id: 2,
			},
		},
		{
			name: "account doesn't exist",
			args: args{
				id: 500,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Init fixtures

			db := getDB()
			defer db.Close()

			err := db.Save(&Account{
				ID:       1,
				Name:     "test",
				Currency: "USD",
			}).Error

			assert.NoError(t, err)

			err = db.Save(&Account{
				ID:       2,
				Name:     "test",
				Currency: "USD",
			}).Error

			assert.NoError(t, err)

			// End of initing fixtures

			s := &basicPaymentsService{
				lock: &dummyLock{},
				uow:  NewUOWPaymentsFactory(db),
			}

			got, err := s.GetAccount(tt.args.ctx, tt.args.id)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("basicPaymentsService.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.NotEqual(t, 0, got.ID, "basicPaymentsService.CreateAccount() error: ID is 0")
			assert.Equal(t, "USD", got.Currency, "basicPaymentsService.CreateAccount() error: currencies aren't equal")
			assert.Equal(t, "test", got.Name, "basicPaymentsService.CreateAccount() error: names aren't equal")
		})
	}
}

func Test_basicPaymentsService_GetAccounts(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    []*Account
		wantErr bool
	}{
		{
			name: "simple getting",
			args: args{},
			want: []*Account{
				{ID: 1, Name: "test1", Currency: "USD", Amount: decimal.Zero},
				{ID: 2, Name: "test2", Currency: "BTC", Amount: decimal.Zero},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db := getDB()
			defer db.Close()

			// Init fixtures
			for _, a := range tt.want {
				err := db.Save(&Account{
					ID:       a.ID,
					Name:     a.Name,
					Currency: a.Currency,
				}).Error

				assert.NoError(t, err)
			}

			s := &basicPaymentsService{
				lock: &dummyLock{},
				uow:  NewUOWPaymentsFactory(db),
			}

			got, err := s.GetAccounts(tt.args.ctx)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("basicPaymentsService.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			for _, a := range got {
				a.CreatedAt = time.Time{}
				a.UpdatedAt = time.Time{}
				assert.True(t, a.Amount.Equal(decimal.Zero))
				a.Amount = decimal.Zero
			}

			assert.ElementsMatch(t, got, tt.want)

		})
	}
}

// ─── MAKE OPERATIONS ────────────────────────────────────────────────────────────

func Test_basicPaymentsService_MakeTransfer(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    map[int64]*Account
		action  func(PaymentsService) (*Operation, error)
		wantErr bool
	}{
		{
			name: "simple transfer",
			args: args{},
			want: map[int64]*Account{
				1: {ID: 1, Name: "test1", Currency: "USD", Amount: decimal.RequireFromString("20")},
				2: {ID: 2, Name: "test2", Currency: "USD", Amount: decimal.RequireFromString("10")},
			},
			action: func(s PaymentsService) (*Operation, error) {
				return s.MakeTransfer(nil, 2, 1, "USD", decimal.RequireFromString("5"))
			},
		},
		{
			name: "full transfer",
			args: args{},
			want: map[int64]*Account{
				1: {ID: 1, Name: "test1", Currency: "USD", Amount: decimal.RequireFromString("30")},
				2: {ID: 2, Name: "test2", Currency: "USD", Amount: decimal.RequireFromString("0")},
			},
			action: func(s PaymentsService) (*Operation, error) {
				return s.MakeTransfer(nil, 2, 1, "USD", decimal.RequireFromString("15"))
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db := getDB()
			defer db.Close()

			// Init fixtures
			for _, a := range tt.want {
				err := db.Save(&Account{
					ID:       a.ID,
					Name:     a.Name,
					Currency: a.Currency,
					Amount:   decimal.RequireFromString("15"),
				}).Error

				assert.NoError(t, err)
			}

			s := &basicPaymentsService{
				lock: &dummyLock{},
				uow:  NewUOWPaymentsFactory(db),
			}

			_, err := tt.action(s)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			for _, wantAcc := range tt.want {
				a, err := s.GetAccount(tt.args.ctx, wantAcc.ID)
				if !assert.NoError(t, err) {
					t.FailNow()
				}
				assert.NotEmpty(t, a.CreatedAt)
				assert.NotEmpty(t, a.UpdatedAt)

				want := tt.want[a.ID]
				if !assert.NotNil(t, want) {
					t.FailNow()
				}

				assert.Equal(t, want.Currency, a.Currency)
				assert.Equal(t, want.Name, a.Name)
				assert.True(t, a.Amount.Equal(want.Amount), "Got: %s; Want: %s", a.Amount, want.Amount)
			}
		})
	}
}
