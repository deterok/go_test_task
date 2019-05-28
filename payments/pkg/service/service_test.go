package service

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

//
// ────────────────────────────────────────────────── I ──────────
//   :::::: U T I L S : :  :   :    :     :        :          :
// ────────────────────────────────────────────────────────────
//

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

func getRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 5 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "redis:6379") },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	pool.Get().Do("FLUSHALL")

	return pool
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

			redis := getRedis()
			defer redis.Close()

			s := &basicPaymentsService{
				lockf: NewLockFactory(redis),
				uowf:  NewUOWPaymentsFactory(db),
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

			redis := getRedis()
			defer redis.Close()

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
				lockf: NewLockFactory(redis),
				uowf:  NewUOWPaymentsFactory(db),
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

			redis := getRedis()
			defer redis.Close()

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
				lockf: NewLockFactory(redis),
				uowf:  NewUOWPaymentsFactory(db),
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

			redis := getRedis()
			defer redis.Close()

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
				lockf: NewLockFactory(redis),
				uowf:  NewUOWPaymentsFactory(db),
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
