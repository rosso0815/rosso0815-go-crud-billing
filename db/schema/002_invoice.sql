-- +goose Up

CREATE SEQUENCE table_invoice_id_seq;

ALTER SEQUENCE table_invoice_id_seq RESTART WITH 100101;

CREATE TABLE Invoice (

    invoice_id int NOT NULL DEFAULT nextval('table_invoice_id_seq') PRIMARY KEY,

    created_at timestamp NOT NULL default timezone ('utc', now()),
    modified_at timestamp not null default now(),

    -- customer_id int not null,
    customer_id integer references customer (customer_id),

    invoice_month int not null,
    invoice_year int not null,
    remark text not null,
    invoice_modified_at timestamp not null default now(),

    ahv float not null default 0.0,
    alv float not null default 0.0,
    quellensteuer float not null default 0.0,
    hours_total float not null default 0.0,
    salary_hour float not null default 0.0,

    bill_add float not null default 0.0,
    bill_sum float not null default 0.0,

    bill_payed boolean not null default false,
    bill_payed_at timestamp not null default now()
);

CREATE UNIQUE INDEX invoice_customer_month ON Invoice(customer_id, invoice_month, invoice_year);

-- +goose StatementBegin
CREATE FUNCTION sync_lastmod() RETURNS trigger AS $$
BEGIN
  NEW.modified_at := NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER
  sync_lastmod
BEFORE UPDATE ON
  invoice
FOR EACH ROW EXECUTE PROCEDURE
  sync_lastmod();

CREATE SEQUENCE table_invoiceentry_id_seq;

ALTER SEQUENCE table_invoiceentry_id_seq RESTART WITH 101101;

CREATE TABLE InvoiceEntry (
    invoiceentry_id integer NOT NULL DEFAULT nextval('table_invoiceentry_id_seq') PRIMARY KEY,
    invoice_id integer references invoice (invoice_id) on delete cascade,
    work_day int not null default 0,
    work_hours real not null default 0.0
);

CREATE UNIQUE INDEX ie_id_day ON InvoiceEntry(invoice_id, work_day);

-- +goose Down
DROP SEQUENCE IF EXISTS table_invoice_id_seq CASCADE;

DROP TABLE IF EXISTS invoice CASCADE;

DROP SEQUENCE IF EXISTS table_invoiceentry_id_seq CASCADE;

DROP TABLE IF EXISTS invoiceentry CASCADE;

DROP INDEX IF EXISTS ie_id_day;

DROP FUNCTION IF EXISTS sync_lastmod;

--- EOF
