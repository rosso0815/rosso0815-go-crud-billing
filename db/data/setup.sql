-- START

DROP TABLE IF EXISTS customer CASCADE;

DROP SEQUENCE IF EXISTS table_cust_id_seq CASCADE;

CREATE SEQUENCE table_cust_id_seq;

ALTER SEQUENCE table_cust_id_seq RESTART WITH 100001;

CREATE TABLE customer (
  customer_id  integer NOT NULL DEFAULT nextval('table_cust_id_seq') PRIMARY key,
  created_at timestamp NOT NULL default timezone('utc', now()),
  modified_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  first_name text not null,
  last_name text not null,
  street text not null,
  town text not null,
  remark text not null,
  UNIQUE (first_name, last_name)
);

CREATE OR REPLACE FUNCTION update_customer_modified_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at := now();
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_customer_modified_at
BEFORE UPDATE ON customer
FOR EACH ROW
EXECUTE FUNCTION update_customer_modified_at();

-- CREATE FUNCTION customer_lastmod() RETURNS trigger AS $$
-- BEGIN
--   NEW.modified_at := NOW();
--   RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER
--   customer_lastmod
-- BEFORE UPDATE ON
--   customer
-- FOR EACH ROW EXECUTE PROCEDURE
--   customer_lastmod();




-- CREATE OR REPLACE FUNCTION trigger_set_timestamp()
-- RETURNS TRIGGER AS $$
-- BEGIN
--   NEW.modified_at = NOW();
--   RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- --create a trigger to execute the function
-- CREATE TRIGGER set_timestamp
-- BEFORE UPDATE ON public.customer
-- FOR EACH ROW
-- EXECUTE PROCEDURE trigger_set_timestamp();

-- EOF
