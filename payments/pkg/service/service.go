package service

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	ErrDifferentCurrencies = errors.New("accounts currencies must be same")
	ErrBalanceTooLow       = errors.New("balance too low")
)

// PaymentsService describes the interface of the system of account management and cash transactions over this acounts
type PaymentsService interface {
	CreateAccount(ctx context.Context, name, currency string) (*Account, error)
	GetAccount(ctx context.Context, id int64) (*Account, error)
	GetAccounts(ctx context.Context) ([]*Account, error)
	GetAccountOperations(ctx context.Context, accID int64) ([]*Operation, error)
	MakeDeposit(ctx context.Context, to int64, currency string, amount decimal.Decimal) (*Operation, error)
	MakeTransfer(ctx context.Context, from, to int64, currency string, amount decimal.Decimal) (*Operation, error)
}

// ─── INTERFACE REALIZATION ──────────────────────────────────────────────────────

type basicPaymentsService struct {
	lock LockRepository
	uow  UOWPaymentsFactory
}

// NewBasicPaymentsService returns a naive implementation of PaymentsService.
func NewBasicPaymentsService(lock LockRepository, uow UOWPaymentsFactory) PaymentsService {
	return &basicPaymentsService{
		lock: lock,
		uow:  uow,
	}
}

// CreateAccount creates new account
func (s *basicPaymentsService) CreateAccount(ctx context.Context, name, currency string) (*Account, error) {
	uow, err := s.uow.Make()
	if err != nil {
		return nil, errors.Wrap(err, "uow context createing failed")
	}
	defer uow.Save()

	a := &Account{
		Name:     name,
		Currency: currency,
	}

	a, err = uow.Accounts().Create(ctx, a)
	if err != nil {
		uow.Revert()
	}

	return a, nil
}

// GetAccount returns prepared account by id
func (s *basicPaymentsService) GetAccount(ctx context.Context, id int64) (*Account, error) {
	uow, err := s.uow.Make()
	if err != nil {
		return nil, errors.Wrap(err, "uow context createing failed")
	}
	defer uow.Save()

	a, err := uow.Accounts().Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "account (%d) getting failed", id)
	}

	return a, nil
}

// GetAccounts returns all accounts in the system
func (s *basicPaymentsService) GetAccounts(ctx context.Context) ([]*Account, error) {
	uow, err := s.uow.Make()
	if err != nil {
		return nil, errors.Wrap(err, "uow context createing failed")
	}
	defer uow.Save()

	a, err := uow.Accounts().GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "accounts getting failed")
	}

	return a, nil
}

// GetAccountOperations returns operations list of the account
func (s *basicPaymentsService) GetAccountOperations(ctx context.Context, accID int64) ([]*Operation, error) {
	uow, err := s.uow.Make()
	if err != nil {
		return nil, errors.Wrap(err, "uow context createing failed")
	}
	defer uow.Save()

	a, err := uow.Operations().GetByAccID(ctx, accID)
	if err != nil {
		return nil, errors.Wrap(err, "accounts getting failed")
	}

	return a, nil
}

// MakeDeposit creates new deposit operation for the account
func (s *basicPaymentsService) MakeDeposit(ctx context.Context, to int64, currency string, amount decimal.Decimal) (*Operation, error) {
	key := s.getMutexKey(to)
	if err := s.lock.Lock(key); err != nil {
		return nil, errors.Wrapf(err, "mutex (%s) locking failed", key)
	}
	defer s.lock.Unlock(key)

	uow, err := s.uow.Make()
	if err != nil {
		return nil, errors.Wrap(err, "uow context createing failed")
	}
	defer uow.Save()

	a, err := uow.Accounts().Get(ctx, to)
	if err != nil {
		return nil, errors.Wrapf(err, "account (%d) getting failed", to)
	}

	a.Amount = a.Amount.Add(amount)

	if a.Currency != currency {
		return nil, ErrDifferentCurrencies
	}

	if _, err := uow.Accounts().Update(ctx, a); err != nil {
		uow.Revert()
		return nil, errors.Wrapf(err, "account (%d) update failed", a.ID)
	}

	o := &Operation{
		Type: OperationTypeDeposit,
		Transactions: []Transaction{
			{
				From:     WorldAccountID,
				To:       to,
				Currency: currency,
				Amount:   amount,
			},
		},
		Participants: []int64{WorldAccountID, to},
	}

	if _, err := uow.Operations().Create(ctx, o); err != nil {
		return nil, errors.Wrap(err, "operation createing failed")
	}

	return o, nil
}

// MakeTransfer creates new transfer operation for the pair of accounts
func (s *basicPaymentsService) MakeTransfer(ctx context.Context, from int64, to int64, currency string, amount decimal.Decimal) (*Operation, error) {
	uow, err := s.uow.Make()
	if err != nil {
		return nil, errors.Wrap(err, "uow context createing failed")
	}
	defer uow.Save()

	a1, err := uow.Accounts().Get(ctx, from)
	if err != nil {
		return nil, errors.Wrapf(err, "account (%d) getting failed", from)
	}

	a2, err := uow.Accounts().Get(ctx, to)
	if err != nil {
		return nil, errors.Wrapf(err, "account (%d) getting failed", to)
	}

	if a1.Currency != currency || a1.Currency != a2.Currency {
		return nil, ErrDifferentCurrencies
	}

	if a1.Amount.LessThan(amount) {
		return nil, ErrBalanceTooLow
	}

	a1.Amount = a1.Amount.Sub(amount)
	a2.Amount = a2.Amount.Add(amount)

	if _, err := uow.Accounts().Update(ctx, a1); err != nil {
		uow.Revert()
		return nil, errors.Wrapf(err, "account (%d) update failed", a1.ID)
	}

	if _, err := uow.Accounts().Update(ctx, a2); err != nil {
		uow.Revert()
		return nil, errors.Wrapf(err, "account (%d) update failed", a2.ID)
	}

	o := &Operation{
		Transactions: []Transaction{
			{
				From:     from,
				To:       to,
				Currency: currency,
				Amount:   amount,
			},
		},
		Participants: []int64{from, to},
	}
	if _, err := uow.Operations().Create(ctx, o); err != nil {
		return nil, errors.Wrap(err, "operation createing failed")
	}

	return o, nil
}

// ─── HELPER METHODS ─────────────────────────────────────────────────────────────
func (s *basicPaymentsService) getMutexKey(ids ...int64) string {
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.FormatInt(id, 10)
	}

	return strings.Join(strIDs, ":")
}

// New returns a PaymentsService with all of the expected middleware wired in.
func New(lock LockRepository, uowFactory UOWPaymentsFactory, middleware []Middleware) PaymentsService {
	var svc PaymentsService = NewBasicPaymentsService(lock, uowFactory)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
