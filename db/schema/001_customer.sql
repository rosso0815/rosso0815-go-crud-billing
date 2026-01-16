-- +goose Up

CREATE SEQUENCE table_cust_id_seq;

ALTER SEQUENCE table_cust_id_seq RESTART WITH 100001;

CREATE TABLE customer (

  customer_id  integer NOT NULL DEFAULT nextval('table_cust_id_seq') PRIMARY key,

  created_at timestamp NOT NULL default timezone('utc', now()),
  modified_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),

  first_name text not null,
  last_name text not null,
  street text not null default '',
  town text not null default '',
  remark text not null default '',
  email text not null default '',
  phone text not null default '',

  ahv float not null default 0.0,
  alv float not null default 0.0,
  quellensteuer float not null default 0.0,
  salary_hour float not null default 0.0

);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_customer_modified_at()
RETURNS TRIGGER
AS $$
BEGIN
    NEW.modified_at := now();
    RETURN NEW;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_update_customer_modified_at
BEFORE UPDATE ON customer
FOR EACH ROW
EXECUTE FUNCTION update_customer_modified_at();

-- +goose Down

DROP TABLE IF EXISTS customer CASCADE;

DROP SEQUENCE IF EXISTS table_cust_id_seq CASCADE;

DROP FUNCTION update_customer_modified_at;

--- EOF
