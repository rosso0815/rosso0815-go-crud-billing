-- +goose Up

CREATE SEQUENCE userkv_id_seq;

CREATE TABLE userkv (
  userkv_id integer NOT NULL DEFAULT nextval('userkv_id_seq') PRIMARY KEY,
  user_id text NOT NULL,
  key text NOT NULL,
  value text not null,
  created_at timestamp NOT NULL default timezone('utc', now()),
  modified_at timestamp NOT NULL default timezone('utc', now())
);

CREATE UNIQUE INDEX userkv_user_key_uindex ON userkv (user_id, key);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_userkv_updated_at()
RETURNS TRIGGER
AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_update_userkv_updated_at
BEFORE UPDATE ON userkv
FOR EACH ROW
EXECUTE FUNCTION update_userkv_updated_at();

-- +goose Down

DROP TABLE IF EXISTS userkv CASCADE;
DROP SEQUENCE IF EXISTS userkv_id_seq CASCADE;
DROP FUNCTION IF EXISTS update_userkv_updated_at;

--- EOF
