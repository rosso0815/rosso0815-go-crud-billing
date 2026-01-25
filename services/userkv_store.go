package services

import (
	"context"

	db_gen "github.com/rosso0815/rosso0815-go-crud-billing/db/generated"
)

type UserKV struct {
	UserkvId string `json:"userkvid"`
	UserName string `json:"userid"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

func FromDbUserkv(db_userkv db_gen.Userkv) UserKV {
	userkv := UserKV{}
	userkv.Key = db_userkv.Key
	userkv.Value = db_userkv.Value
	return userkv

}

func (m *Store) UserkvList(ctx context.Context, user_id string) ([]UserKV, error) {
	var userkv []UserKV
	db_kvs, err := m.Db.Queries.UserkvList(ctx, m.Db.Db, user_id)
	if err != nil {
		return nil, err
	}
	for _, i := range db_kvs {
		userkv = append(userkv, FromDbUserkv(i))
	}
	return userkv, nil
}

// --- EOF
