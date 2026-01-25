-- name: InvoicesList :many
select
    i.invoice_id
from
    invoice i, customer c
where
    i.customer_id = c.customer_id
and
    (
        c.customer_id::text like '%' || @search::text || '%'
    or
        c.first_name ilike '%' || @search::text || '%'
    or
        c.last_name ilike '%' || @search::text || '%'
    )
ORDER BY i.invoice_id DESC
LIMIT
    @page_size
OFFSET
    @page_count;

-- name: InvoiceGetById :one
select
  *
from
	invoice i
where
	i.invoice_id = @invoice_id;

-- name: InvoiceDeleteById :exec
DELETE FROM invoice WHERE invoice_id = @invoice_id;

-- name: InvoiceByCustomerAndMonthAndYear :one
select invoice_id from invoice
where
	customer_id = @customer_id
AND
	invoice_month = @invoice_month
AND
	invoice_year = @invoice_year;

-- name: InvoiceentryGetByInvoiceidAndDay :one
select
    invoiceentry_id,
    invoice_id,
    work_day,
    work_hours
from
    invoiceentry
where
    invoice_id = @invoice_id
AND
    work_day = @work_day;

-- name: InvoiceInsert :one
INSERT INTO invoice
    (
        customer_id,
        invoice_month,
        invoice_year,
        remark,
        invoice_modified_at,
        ahv,
        alv,
        quellensteuer,
        hours_total,
        salary_hour,
        bill_add,
        bill_sum,
        bill_payed,
        bill_payed_at
    )
VALUES
    (	@customer_id,
        @invoice_month,
        @invoice_year,
        @remark,
        now(),
        @ahv,
        @alv,
        @quellensteuer,
        @totalhours,
        @salary_hour,
        @bill_add,
        @bill_sum,
        false,
        now()
        )
returning invoice_id;

-- name: InvoiceUpdate :execresult
update invoice
set
    customer_id = @customer_id,
    invoice_month = @invoice_month,
    invoice_year = @invoice_year,
    remark	= @remark,
    invoice_modified_at = now(),
    ahv = @ahv,
    alv = @alv,
    quellensteuer = @quellensteuer,
    hours_total = @totalhours,
    salary_hour = @salary_hour,
    bill_add = @bill_add,
    bill_sum = @bill_sum,
    bill_payed = false,
    bill_payed_at = now()
where
   invoice_id = @invoice_id;

-- name: InvoicePayed :execresult
update invoice
set
    bill_payed = true,
    bill_payed_at = now()
where
   invoice_id = @invoice_id;

-- name: InvoiceUnPayed :execresult
update invoice
set
    bill_payed = false,
    bill_payed_at = now()
where
   invoice_id = @invoice_id;

-- name: InvoiceentryInsert :execresult
INSERT INTO invoiceEntry
    (invoice_id, work_day, work_hours)
VALUES
    (@invoice_id, @work_day, @work_hours);

-- name: InvoiceentryUpdate :execresult
UPDATE invoiceEntry
SET
    work_hours = @work_hours
WHERE
    invoice_id = @invoice_id
AND
    work_day = @work_day;

-- EOF
