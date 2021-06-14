package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/nipe0324/gqlgen-todos-example/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	conn *db.SQLHandler
}

func NewResolver(conn *db.SQLHandler) *Resolver {
	return &Resolver{
		conn: conn,
	}
}

// ログインしているユーザを取得
func (r *Resolver) getCurrentUser() (*db.User, error) {
	var user *db.User

	// TODO: ユーザ認証がないのでUserIDは固定で1
	err := r.conn.Where(&db.User{ID: 1}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
