package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	db_gen "github.com/rosso0815/rosso0815-go-crud-billing/db/generated"
)

type InvoiceItem struct {
	Weekday       time.Time
	WorkDayNumber int     `json:"workdaynumber,string"`
	WorkHours     float64 `json:"workhours,string"`
}
type Invoice struct {
	InvoiceId  int `json:"invoiceid,string"`
	CreatedAt  time.Time
	ModifiedAt time.Time

	InvoiceMonth      int       `json:"invoicemonth,string"`
	InvoiceYear       int       `json:"invoiceyear,string"`
	InvoiceModifiedAt time.Time `json:"invoice_modified_at"`
	Remark            string    `json:"remark"`
	Customer          Customer  //`json:"customer"`
	InvoiceItems      []InvoiceItem
	AHV               float64   `json:"ahv,string"`
	ALV               float64   `json:"alv,string"`
	Quellensteuer     float64   `json:"quellensteuer,string"`
	TotalHours        float64   `json:"totalhours,string"`
	SalaryHour        float64   `json:"salaryhour,string"`
	BillingAdd        float64   `json:"billingadd,string"`
	BillingSum        float64   `json:"billingsum,string"`
	BillPayed         bool      `json:"bill_payed"`
	BillPayedAt       time.Time `json:"bill_payed_at"`
}

func toStoreInvoice(db_i db_gen.Invoice) Invoice {
	i := Invoice{}
	i.InvoiceId = db_i.InvoiceID
	i.ModifiedAt = db_i.ModifiedAt
	i.Customer.CustomerID = db_i.CustomerID
	i.InvoiceYear = db_i.InvoiceYear
	i.InvoiceMonth = db_i.InvoiceMonth
	i.InvoiceModifiedAt = db_i.InvoiceModifiedAt
	i.Remark = db_i.Remark
	i.AHV = db_i.Ahv
	i.ALV = db_i.Alv
	i.Quellensteuer = db_i.Quellensteuer
	i.TotalHours = db_i.HoursTotal
	i.SalaryHour = db_i.SalaryHour
	i.BillingAdd = db_i.BillAdd
	i.BillingSum = db_i.BillSum
	i.BillPayed = db_i.BillPayed
	i.BillPayedAt = db_i.BillPayedAt
	return i
}

func (inv *Invoice) Calculate() {
	inv.TotalHours = 0.0
	for _, i := range inv.InvoiceItems {
		inv.TotalHours = inv.TotalHours + i.WorkHours
	}
	inv.BillingSum = inv.TotalHours * inv.SalaryHour
	inv.BillingSum = inv.BillingSum - (inv.TotalHours * inv.SalaryHour * inv.AHV / 100)
	inv.BillingSum = inv.BillingSum - (inv.TotalHours * inv.SalaryHour * inv.ALV / 100)
	inv.BillingSum = inv.BillingSum - (inv.TotalHours * inv.SalaryHour * inv.Quellensteuer / 100)
	inv.BillingSum = inv.BillingSum + inv.BillingAdd
	inv.BillingSum = math.Round(inv.BillingSum*20) / 20
}

func (i *Invoice) GetAhvTotal() string {
	return fmt.Sprintf("%.2f", i.TotalHours*i.SalaryHour*i.AHV/100)
}

func (i *Invoice) GetAlvTotal() string {
	return fmt.Sprintf("%.2f", i.TotalHours*i.SalaryHour*i.ALV/100)
}

func (i *Invoice) GetQuellensteuerTotal() string {
	return fmt.Sprintf("%.2f", i.TotalHours*i.SalaryHour*i.Quellensteuer/100)
}

func (i *Invoice) GetTotalHours() string {
	return fmt.Sprintf("%.2f", i.TotalHours)
}

// ---

func (m *Store) InvoiceListBySearch(ctx context.Context, search string, page_size int, page_count int) ([]Invoice, error) {
	var invoices []Invoice
	ids, err := m.Db.Queries.InvoicesList(ctx, m.Db.Db, db_gen.InvoicesListParams{
		Search:    search,
		PageCount: int64(page_count),
		PageSize:  int64(page_size),
	})
	if err != nil {
		log.Println("Error invoice search", err)
		return nil, err
	}
	for _, i := range ids {
		invoice, err := m.Db.Queries.InvoiceGetById(ctx, m.Db.Db, i)
		if err != nil {
			log.Println("Error invoice search", err)
			return invoices, err
		}
		s_invoice := toStoreInvoice(invoice)
		customer, err := m.Db.Queries.CustomerGetById(ctx, m.Db.Db, invoice.CustomerID)
		if err != nil {
			log.Println("Error invoice search", err)
			return invoices, err
		}
		s_invoice.Customer = toStoreCustomer(customer)
		invoices = append(invoices, s_invoice)
	}
	return invoices, nil
}

func (m *Store) InvoiceGetById(ctx context.Context, invoiceId int) (Invoice, error) {
	db_invoice, err := m.Db.Queries.InvoiceGetById(ctx, m.Db.Db, invoiceId)
	if err != nil {
		log.Println("ERROR GetInvoiceById", err)
		return Invoice{}, err
	}
	invoice := toStoreInvoice(db_invoice)
	customer, err := m.Db.Queries.CustomerGetById(ctx, m.Db.Db, db_invoice.CustomerID)
	if err != nil {
		log.Println("Error invoice search", err)
		return invoice, err
	}
	invoice.Customer = toStoreCustomer(customer)
	for i := range time.Date(invoice.InvoiceYear, time.Month(invoice.InvoiceMonth)+1, 0, 0, 0, 0, 0, time.UTC).Day() {
		ie, err := m.Db.Queries.InvoiceentryGetByInvoiceidAndDay(ctx, m.Db.Db, db_gen.InvoiceentryGetByInvoiceidAndDayParams{
			InvoiceID: db_invoice.InvoiceID,
			WorkDay:   i + 1})
		if err != nil && err.Error() != pgx.ErrNoRows.Error() {
			log.Println("Error invoice search", err)
			return invoice, err
		}
		s_invoiceentry := InvoiceItem{
			WorkDayNumber: i + 1,
			WorkHours:     ie.WorkHours,
			Weekday: time.Date(invoice.InvoiceYear,
				time.Month(invoice.InvoiceMonth), i+1, 20, 34, 58, 651387237, time.UTC),
		}
		invoice.InvoiceItems = append(invoice.InvoiceItems, s_invoiceentry)
	}
	return invoice, nil
}

func (m *Store) InvoiceDeleteById(ctx context.Context, invoiceId int) error {
	err := m.Db.Queries.InvoiceDeleteById(ctx, m.Db.Db, invoiceId)
	if err != nil {
		return err
	}
	return nil
}

// creates or updates a invoice based on customer with month with year
func (m *Store) InvoiceGetByCustomerAndMonth(ctx context.Context, customerId int, invoice_month int, invoice_year int) (Invoice, error) {
	var db_invoice db_gen.Invoice
	var invoice Invoice = Invoice{
		InvoiceMonth: invoice_month,
		InvoiceYear:  invoice_year,
	}
	var err error
	db_invoice.InvoiceID, err = m.Db.Queries.InvoiceByCustomerAndMonthAndYear(ctx, m.Db.Db, db_gen.InvoiceByCustomerAndMonthAndYearParams{
		CustomerID:   customerId,
		InvoiceMonth: invoice_month,
		InvoiceYear:  invoice_year,
	})
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		log.Fatal("InvoiceGet failed:", err)
		return Invoice{}, err
	}
	db_customer, err := m.Db.Queries.CustomerGetById(ctx, m.Db.Db, customerId)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		log.Fatal("InvoiceGet failed:", err)
		return Invoice{}, err
	}

	if db_invoice.InvoiceID != 0 {
		db_invoice, err = m.Db.Queries.InvoiceGetById(ctx, m.Db.Db, db_invoice.InvoiceID)
		if err != nil {
			log.Fatal("InvoiceGet failed:", err)
			return Invoice{}, err
		}
		invoice = toStoreInvoice(db_invoice)
	} else {
		invoice.AHV = db_customer.Ahv
		invoice.ALV = db_customer.Alv
		invoice.Quellensteuer = db_customer.Quellensteuer
		invoice.SalaryHour = db_customer.SalaryHour
	}

	invoice.Customer = toStoreCustomer(db_customer)
	for i := range time.Date(invoice.InvoiceYear, time.Month(invoice.InvoiceMonth)+1, 0, 0, 0, 0, 0, time.UTC).Day() {
		ie, err := m.Db.Queries.InvoiceentryGetByInvoiceidAndDay(ctx, m.Db.Db, db_gen.InvoiceentryGetByInvoiceidAndDayParams{
			InvoiceID: db_invoice.InvoiceID,
			WorkDay:   i + 1})
		if err != nil && err.Error() != pgx.ErrNoRows.Error() {
			log.Println("Error invoice search", err)
			return invoice, err
		}
		s_invoiceentry := InvoiceItem{
			WorkDayNumber: i + 1,
			WorkHours:     ie.WorkHours,
			Weekday: time.Date(invoice.InvoiceYear,
				time.Month(invoice.InvoiceMonth), i+1, 20, 34, 58, 651387237, time.UTC),
		}
		invoice.InvoiceItems = append(invoice.InvoiceItems, s_invoiceentry)
	}
	return invoice, nil
}

func (m *Store) InvoiceSave(ctx context.Context, invoice Invoice) error {
	invoice.Calculate()
	b_invoice, _ := m.Db.Queries.InvoiceByCustomerAndMonthAndYear(ctx, m.Db.Db, db_gen.InvoiceByCustomerAndMonthAndYearParams{
		CustomerID:   invoice.Customer.CustomerID,
		InvoiceYear:  invoice.InvoiceYear,
		InvoiceMonth: invoice.InvoiceMonth,
	})
	if b_invoice == 0 {
		id, err := m.Db.Queries.InvoiceInsert(ctx, m.Db.Db, db_gen.InvoiceInsertParams{
			CustomerID:    invoice.Customer.CustomerID,
			InvoiceMonth:  invoice.InvoiceMonth,
			InvoiceYear:   invoice.InvoiceYear,
			Remark:        invoice.Remark,
			Ahv:           invoice.AHV,
			Alv:           invoice.ALV,
			Quellensteuer: invoice.Quellensteuer,
			Totalhours:    invoice.TotalHours,
			SalaryHour:    invoice.SalaryHour,
			BillAdd:       invoice.BillingAdd,
			BillSum:       invoice.BillingSum,
		})
		invoice.InvoiceId = id
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err := m.Db.Queries.InvoiceUpdate(ctx, m.Db.Db, db_gen.InvoiceUpdateParams{
			InvoiceID:     b_invoice,
			CustomerID:    invoice.Customer.CustomerID,
			InvoiceMonth:  invoice.InvoiceMonth,
			InvoiceYear:   invoice.InvoiceYear,
			Remark:        invoice.Remark,
			Ahv:           invoice.AHV,
			Alv:           invoice.ALV,
			Quellensteuer: invoice.Quellensteuer,
			Totalhours:    invoice.TotalHours,
			SalaryHour:    invoice.SalaryHour,
			BillAdd:       invoice.BillingAdd,
			BillSum:       invoice.BillingSum,
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}
	for _, i := range invoice.InvoiceItems {
		ie, err := m.Db.Queries.InvoiceentryGetByInvoiceidAndDay(ctx, m.Db.Db, db_gen.InvoiceentryGetByInvoiceidAndDayParams{
			InvoiceID: invoice.InvoiceId,
			WorkDay:   i.WorkDayNumber,
		})
		if err != nil && err.Error() != pgx.ErrNoRows.Error() {
			log.Println(err)
			return err
		}
		if ie.InvoiceID == 0 {
			_, err := m.Db.Queries.InvoiceentryInsert(ctx, m.Db.Db, db_gen.InvoiceentryInsertParams{
				InvoiceID: invoice.InvoiceId,
				WorkDay:   i.WorkDayNumber,
				WorkHours: i.WorkHours,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		} else {
			_, err := m.Db.Queries.InvoiceentryUpdate(ctx, m.Db.Db, db_gen.InvoiceentryUpdateParams{
				InvoiceID: invoice.InvoiceId,
				WorkDay:   i.WorkDayNumber,
				WorkHours: i.WorkHours,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (m *Store) InvoicePayed(ctx context.Context, invoice_id int) error {
	_, err := m.Db.Queries.InvoicePayed(ctx, m.Db.Db, invoice_id)
	return err
}

func (m *Store) InvoiceNotPayed(ctx context.Context, invoice_id int) error {
	_, err := m.Db.Queries.InvoiceUnPayed(ctx, m.Db.Db, invoice_id)
	return err
}

// --- EOF
