package services

import (
	"context"
	"log"
	"testing"

	"github.com/rosso0815/rosso0815-go-crud-billing/config"
	"github.com/rosso0815/rosso0815-go-crud-billing/db"
)

func Test_PgxInvoiceSave(t *testing.T) {

	t.Log("@@@ Test_InvoiceList")
	log.Println("start")
	cfg := config.New(nil)
	ctx := context.Background()
	store := NewStore(ctx, cfg)
	db.LoadSQLFile(store.Db.Db, "../db/data/001_customer.sql")
	db.LoadSQLFile(store.Db.Db, "../db/data/002_invoice.sql")
	db.LoadSQLFile(store.Db.Db, "../db/data/003_userkv.sql")

	// 	sql := `
	// -- start
	// delete from invoice if exists;
	// delete  from customer;
	// --done
	// 	`
	// conn, err := store.Db.Db.Exec(ctx, string(query))
	// log.Println("conn", conn)
	// log.Println("err", err)
	// invoicesSearch := store.InvoiceListBySearch(ctx, "", 10, 0)
	// for _, i := range invoicesSearch {
	// 	fmt.Println("i", i)
	// }
	// assert.Equal(t, 3, len(invoices), "they should be equal")

	// invoice := store.InvoiceGetById(ctx, 100102)
	// fmt.Println("invoice", invoice.Customer.FirstName, invoice.Customer.LastName)
	// for _, i := range invoice.InvoiceItems {
	// 	fmt.Println("invoiceItem", i)
	// }
	// assert.Equal(t, item.InvoiceId, 100101, "should be equal")

	// invoice := store.InvoiceGet(ctx, 100002, 2, 2025)
	// t.Log("invoice:", invoice.InvoiceId, invoice.InvoiceMonth, invoice.InvoiceYear, invoice.Remark)
	// t.Log("invoice.Customer", invoice.Customer)
	// assert.Equal(t, 1, len(store.InvoiceListBySearch(ctx, "2", 10, 0)), "they should be equal")

	// assert.Equal(t, 1, store.InvoiceListBySearchCount(ctx, "2"), "they should be equal")

	// err = store.InvoiceDelete(ctx, 100102)
	// assert.Equal(t, err, nil)

	// err = store.InvoiceDelete(ctx, 100108)
	// assert.Error(t, err, "ItemDelete: expected 1 row affected")

	// id, _ := store.InvoiceInsert(ctx, Invoice{
	// 	// CustomerId:   1,
	// 	InvoiceMonth: 1,
	// 	InvoiceYear:  1,
	// 	Remark:       "Test",
	// })
	// log.Println("Inserted Customer ID:", id)

	// store.InvoiceUpdate(ctx, Invoice{
	// 	InvoiceId: id,
	// 	Remark:    "TestUpdate",
	// })

	// items := store.InvoiceItemGet(ctx, 100101)
	// assert.Equal(t, 2, len(items), "should be equal")
	// for i, item := range items {
	// 	fmt.Printf("i:%d day:%d duration:%.2f\n", i, item.WorkDay, item.WorkHours)
	// }

	// var invoiceItems []InvoiceItem
	// invoiceItems = append(invoiceItems, InvoiceItem{
	// 	InvoiceId: 100101,
	// 	WorkDay:   1,
	// 	WorkHours: 8.5,
	// })
	// invoiceItems = append(invoiceItems, InvoiceItem{
	// 	InvoiceId: 100101,
	// 	WorkDay:   2,
	// 	WorkHours: 7.5,
	// })
	// err = store.InvoiceItemSet(ctx, 100101, invoiceItems)
	// assert.Equal(t, err, nil)

}

func Test_PgxInvoice(t *testing.T) {
	t.Log("@@@ Test_InvoiceList")
	// log.Println("start")
	// cfg := config.New(nil)
	// ctx := context.Background()
	// store := NewStore(ctx, cfg)
	// t.Log(store)

	// invoices, err := store.Db.Queries.InvoicesList(ctx, store.Db.Db, db_gen.InvoicesListParams{
	// 	PageCount: 0,
	// 	PageSize:  100,
	// })
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// t.Log("invoice len:", len(invoices))

	// for i, item := range invoices {
	// 	t.Log(i, item)
	// }

	// id := 100111
	// store.Db.Queries.InvoiceGetById(ctx, store.Db.Db, id)
	// t.Log("done")
	// invoice_id, err := store.Db.Queries.InvoiceIdByCustomerAndMonth(ctx, store.Db.Db, db_gen.InvoiceIdByCustomerAndMonthParams{
	// 	CustomerID:   1,
	// 	InvoiceYear:  2000,
	// 	InvoiceMonth: 1,
	// })
	// if err != nil && err != pgx.ErrNoRows {
	// 	log.Panicln(err)
	// }
	// t.Log(invoice_id)
	// conn := db.New(ctx, nil)
	// defer conn.Close()
	// store := NewPgxStore(conn)
	// log.Println("handler", store)

	// err := db.LoadSQLFile(conn, "setup.sql")
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// err = db.LoadSQLFile(conn, "data.sql")
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// invoice := Invoice{
	//   InvoiceMonth: 12,
	//   InvoiceYear:  2025,
	//   Remark:       "12/25",
	//   Customer: Customer{
	//     CustomerID: 100001,
	//   },
	//   InvoiceItems: []InvoiceItem{
	//     {
	//       WorkDayNumber: 1,
	//       WorkHours:     1,
	//     },
	//     {
	//       WorkDayNumber: 2,
	//       WorkHours:     1.5,
	//     },
	//     {
	//       WorkDayNumber: 3,
	//       WorkHours:     1.2,
	//     },
	//     {
	//       WorkDayNumber: 4,
	//       WorkHours:     5,
	//     },
	//     {
	//       WorkDayNumber: 5,
	//       WorkHours:     6.3,
	//     },
	//     {
	//       WorkDayNumber: 6,
	//       WorkHours:     0.4,
	//     },
	//   },
	// }
	//
	// store.InvoiceSave(ctx, invoice)

	// invoicesSearch := store.InvoiceListBySearch(ctx, "", 10, 0)
	// for _, i := range invoicesSearch {
	// 	fmt.Println("i", i)
	// }
	// assert.Equal(t, 3, len(invoices), "they should be equal")

	// invoice := store.InvoiceGetById(ctx, 100102)
	// fmt.Println("invoice", invoice.Customer.FirstName, invoice.Customer.LastName)
	// for _, i := range invoice.InvoiceItems {
	// 	fmt.Println("invoiceItem", i)
	// }
	// assert.Equal(t, item.InvoiceId, 100101, "should be equal")

	// invoice := store.InvoiceGet(ctx, 100002, 2, 2025)
	// t.Log("invoice:", invoice.InvoiceId, invoice.InvoiceMonth, invoice.InvoiceYear, invoice.Remark)
	// t.Log("invoice.Customer", invoice.Customer)
	// assert.Equal(t, 1, len(store.InvoiceListBySearch(ctx, "2", 10, 0)), "they should be equal")

	// assert.Equal(t, 1, store.InvoiceListBySearchCount(ctx, "2"), "they should be equal")

	// err = store.InvoiceDelete(ctx, 100102)
	// assert.Equal(t, err, nil)

	// err = store.InvoiceDelete(ctx, 100108)
	// assert.Error(t, err, "ItemDelete: expected 1 row affected")

	// id, _ := store.InvoiceInsert(ctx, Invoice{
	// 	// CustomerId:   1,
	// 	InvoiceMonth: 1,
	// 	InvoiceYear:  1,
	// 	Remark:       "Test",
	// })
	// log.Println("Inserted Customer ID:", id)

	// store.InvoiceUpdate(ctx, Invoice{
	// 	InvoiceId: id,
	// 	Remark:    "TestUpdate",
	// })

	// items := store.InvoiceItemGet(ctx, 100101)
	// assert.Equal(t, 2, len(items), "should be equal")
	// for i, item := range items {
	// 	fmt.Printf("i:%d day:%d duration:%.2f\n", i, item.WorkDay, item.WorkHours)
	// }

	// var invoiceItems []InvoiceItem
	// invoiceItems = append(invoiceItems, InvoiceItem{
	// 	InvoiceId: 100101,
	// 	WorkDay:   1,
	// 	WorkHours: 8.5,
	// })
	// invoiceItems = append(invoiceItems, InvoiceItem{
	// 	InvoiceId: 100101,
	// 	WorkDay:   2,
	// 	WorkHours: 7.5,
	// })
	// err = store.InvoiceItemSet(ctx, 100101, invoiceItems)
	// assert.Equal(t, err, nil)

}
