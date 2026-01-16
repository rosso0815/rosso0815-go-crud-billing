package services

import (
	"context"
	"errors"
	"log"
	"time"

	db_gen "github.com/rosso0815/rosso0815-go-crud-billing/db/generated"
)

type Customer struct {
	CustomerID    int `json:"customerid,string"`
	CreatedAt     time.Time
	ModifiedAt    time.Time
	FirstName     string  `json:"firstname"`
	LastName      string  `json:"lastname"`
	Street        string  `json:"street"`
	Town          string  `json:"town"`
	Remark        string  `json:"remark"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Ahv           float64 `json:"ahv,string"`
	Alv           float64 `json:"alv,string"`
	Quellensteuer float64 `json:"quellensteuer,string"`
	SalaryHour    float64 `json:"salaryhour,string"`
}

func toStoreCustomer(d db_gen.Customer) Customer {
	c := Customer{}
	c.CustomerID = d.CustomerID
	c.FirstName = d.FirstName
	c.LastName = d.LastName
	c.Street = d.Street
	c.Town = d.Town
	c.Remark = d.Remark
	c.Email = d.Email
	c.Phone = d.Phone
	c.Ahv = d.Ahv
	c.Alv = d.Alv
	c.Quellensteuer = d.Quellensteuer
	c.SalaryHour = d.SalaryHour
	return c
}

func (m *Store) CustomerList(ctx context.Context) ([]Customer, error) {
	i, err := m.Db.Queries.CustomersList(ctx, m.Db.Db)
	if err != nil {
		return nil, err
	}
	var customers []Customer
	for _, id := range i {
		customer, err := m.Db.Queries.CustomerGetById(ctx, m.Db.Db, id)
		if err != nil {
			return nil, err
		}
		customers = append(customers, toStoreCustomer(customer))
	}
	return customers, nil
}

func (m *Store) CustomerListBySearch(ctx context.Context, search string, page_size int, page_count int) ([]Customer, error) {
	i, err := m.Db.Queries.CustomersListBySearch(ctx, m.Db.Db,
		db_gen.CustomersListBySearchParams{
			Search:    search,
			PageCount: int64(page_count),
			PageSize:  int64(page_size)})
	if err != nil {
		return nil, err
	}
	var customers []Customer
	for _, id := range i {
		customer, err := m.Db.Queries.CustomerGetById(ctx, m.Db.Db, id)
		if err != nil {
			return nil, err
		}
		customers = append(customers, toStoreCustomer(customer))
	}
	return customers, nil
}

func (m *Store) CustomerListBySearchCount(ctx context.Context, search string) (int, error) {
	count, err := m.Db.Queries.CustomersListCount(ctx, m.Db.Db, search)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return int(count), nil
}

func (m *Store) CustomerGetById(ctx context.Context, customerId int) (Customer, error) {
	s_customer, err := m.Db.Queries.CustomerGetById(ctx, m.Db.Db, customerId)
	if err != nil {
		log.Printf("Error CustomerGetById %d %s", customerId, err)
		return Customer{}, err
	}
	return toStoreCustomer(s_customer), nil
}

func (m *Store) CustomerDelete(ctx context.Context, customerId int) error {
	err := m.Db.Queries.CustomerDelete(ctx, m.Db.Db, customerId)
	if err != nil {
		log.Printf("Error CustomerGetById %d %s", customerId, err)
		return err
	}
	return nil
}

func (m *Store) CustomerCreate(ctx context.Context, customer Customer) (itemId int, err error) {
	s_customer := db_gen.CustomerCreateParams{
		Firstname:     customer.FirstName,
		Lastname:      customer.LastName,
		Street:        customer.Street,
		Town:          customer.Town,
		Remark:        customer.Remark,
		Email:         customer.Email,
		Phone:         customer.Phone,
		Ahv:           customer.Ahv,
		Alv:           customer.Alv,
		Quellensteuer: customer.Quellensteuer,
		Salaryhour:    customer.SalaryHour,
	}

	if len(s_customer.Firstname) < 3 || len(s_customer.Lastname) < 3 {
		return 0, errors.New("firstname or lastname are too short")
	}

	id, err := m.Db.Queries.CustomerCreate(ctx, m.Db.Db, s_customer)
	if err != nil {
		log.Println("Error:", err)
		return 0, err
	}
	return id, nil
}

func (m *Store) CustomerUpdate(ctx context.Context, customer Customer) error {
	s_customer := db_gen.CustomerUpdateParams{
		Firstname:     customer.FirstName,
		Lastname:      customer.LastName,
		Street:        customer.Street,
		Town:          customer.Town,
		Remark:        customer.Remark,
		Email:         customer.Email,
		Phone:         customer.Phone,
		Customerid:    customer.CustomerID,
		Ahv:           customer.Ahv,
		Alv:           customer.Alv,
		Quellensteuer: customer.Quellensteuer,
		Salaryhour:    customer.SalaryHour,
	}
	err := m.Db.Queries.CustomerUpdate(ctx, m.Db.Db, s_customer)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	return nil
}

// --- EOF
