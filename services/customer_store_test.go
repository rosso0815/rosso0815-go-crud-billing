package services

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rosso0815/rosso0815-go-crud-billing/config"
	"github.com/rosso0815/rosso0815-go-crud-billing/db"
	db_gen "github.com/rosso0815/rosso0815-go-crud-billing/db/generated"
	"gotest.tools/assert"
)

func newCustomerTestStore(t *testing.T) *Store {
	t.Helper()

	dbURI := os.Getenv("GOAPP_DB_URI")
	if dbURI == "" {
		dbURI = os.Getenv("DB_URI")
	}
	if dbURI == "" {
		t.Skip("GOAPP_DB_URI or DB_URI must be set for integration tests")
	}

	t.Setenv("DB_URI", dbURI)
	t.Setenv("WEB_LISTENER", ":0")

	ctx := context.Background()
	store := NewStore(ctx, config.New(nil))
	t.Cleanup(store.Close)

	_, err := store.Db.Db.Exec(ctx, "delete from invoiceentry; delete from invoice; delete from customer;")
	assert.NilError(t, err)

	err = db.LoadSQLFile(store.Db.Db, "../db/data/001_customer.sql")
	assert.NilError(t, err)

	return store
}

func Test_toStoreCustomer(t *testing.T) {
	now := time.Now()
	in := db_gen.Customer{
		CustomerID:    42,
		CreatedAt:     now,
		ModifiedAt:    now,
		FirstName:     "Ada",
		LastName:      "Lovelace",
		Street:        "Main Street",
		Town:          "Bern",
		Remark:        "VIP",
		Email:         "ada@example.com",
		Phone:         "0041",
		Ahv:           5.3,
		Alv:           1.1,
		Quellensteuer: 5.0,
		SalaryHour:    88,
	}

	out := toStoreCustomer(in)

	assert.Equal(t, out.CustomerID, in.CustomerID)
	assert.Equal(t, out.FirstName, in.FirstName)
	assert.Equal(t, out.LastName, in.LastName)
	assert.Equal(t, out.Street, in.Street)
	assert.Equal(t, out.Town, in.Town)
	assert.Equal(t, out.Remark, in.Remark)
	assert.Equal(t, out.Email, in.Email)
	assert.Equal(t, out.Phone, in.Phone)
	assert.Equal(t, out.Ahv, in.Ahv)
	assert.Equal(t, out.Alv, in.Alv)
	assert.Equal(t, out.Quellensteuer, in.Quellensteuer)
	assert.Equal(t, out.SalaryHour, in.SalaryHour)
}

func Test_CustomerListAndSearch(t *testing.T) {
	store := newCustomerTestStore(t)
	ctx := context.Background()

	list, err := store.CustomerList(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 7, len(list))

	filtered, err := store.CustomerListBySearch(ctx, "last_name2", 10, 0)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "customer_02", filtered[0].FirstName)

	count, err := store.CustomerListBySearchCount(ctx, "customer_1")
	assert.NilError(t, err)
	assert.Equal(t, 4, count)
}

func Test_CustomerCRUD(t *testing.T) {
	store := newCustomerTestStore(t)
	ctx := context.Background()

	id, err := store.CustomerCreate(ctx, Customer{
		FirstName:     "create_first",
		LastName:      "create_last",
		Street:        "Street 1",
		Town:          "Town 1",
		Remark:        "Remark 1",
		Email:         "create@example.com",
		Phone:         "111",
		Ahv:           5.3,
		Alv:           1.1,
		Quellensteuer: 5.0,
		SalaryHour:    45,
	})
	assert.NilError(t, err)
	assert.Assert(t, id > 0)

	created, err := store.CustomerGetById(ctx, id)
	assert.NilError(t, err)
	assert.Equal(t, "create_first", created.FirstName)
	assert.Equal(t, "create_last", created.LastName)

	err = store.CustomerUpdate(ctx, Customer{
		CustomerID:    id,
		FirstName:     "update_first",
		LastName:      "update_last",
		Street:        "Street 2",
		Town:          "Town 2",
		Remark:        "Remark 2",
		Email:         "update@example.com",
		Phone:         "222",
		Ahv:           4.5,
		Alv:           0.9,
		Quellensteuer: 4.2,
		SalaryHour:    99,
	})
	assert.NilError(t, err)

	updated, err := store.CustomerGetById(ctx, id)
	assert.NilError(t, err)
	assert.Equal(t, "update_first", updated.FirstName)
	assert.Equal(t, "update_last", updated.LastName)
	assert.Equal(t, "update@example.com", updated.Email)
	assert.Equal(t, 99.0, updated.SalaryHour)

	err = store.CustomerDelete(ctx, id)
	assert.NilError(t, err)

	_, err = store.CustomerGetById(ctx, id)
	assert.Assert(t, err != nil)
	assert.Assert(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_CustomerCreateValidation(t *testing.T) {
	store := newCustomerTestStore(t)
	ctx := context.Background()

	id, err := store.CustomerCreate(ctx, Customer{
		FirstName: "ab",
		LastName:  "cd",
	})

	assert.Equal(t, 0, id)
	assert.Assert(t, err != nil)
	assert.Error(t, err, "firstname or lastname are too short")
}
