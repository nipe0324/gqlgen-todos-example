package dataloader

import (
	"context"
	"fmt"
	"time"

	"github.com/nipe0324/gqlgen-todos-example/db"
)

// Context に DataLoader を保持するキー
const CtxKey = "dataloader"

type DataLoader struct {
	UserByID *UserLoader
}

func NewDataLoader(conn *db.SQLHandler) *DataLoader {
	return &DataLoader{
		UserByID: newUserByIDLoader(conn),
	}
}

func newUserByIDLoader(conn *db.SQLHandler) *UserLoader {
	conf := UserLoaderConfig{
		// Fetch はデータを取得するメソッドを実装する
		Fetch: func(keys []uint64) ([]*db.User, []error) {
			// DBからユーザ情報を取得
			var dbUsers []*db.User
			if err := conn.Find(&dbUsers, keys).Error; err != nil {
				return nil, []error{}
			}

			// ユーザIDとユーザ情報をマップにする
			userByID := map[uint64]*db.User{}

			for _, dbUser := range dbUsers {
				userByID[dbUser.ID] = dbUser
			}

			// ユーザとエラーのリストを作成する
			// keys のインデックス値に該当のユーザやエラーを設定する
			users := make([]*db.User, len(keys))
			errors := make([]error, len(keys))

			for i, key := range keys {
				user := userByID[key]
				if user != nil {
					users[i] = user
				} else {
					errors[i] = fmt.Errorf("not found user by id : %d", key)
				}
			}

			return users, errors
		},

		// Wait はバッチでデータ取得するまでに待つ時間
		// 時間が長いほど処理が遅くなる。時間が短いほどバッチでデータを取得できなくなるトレードオフ
		Wait: 10 * time.Millisecond,

		// MaxBatch は１バッチあたりの最大のキーの数。0を指定したは制限なし
		MaxBatch: 100,
	}

	return NewUserLoader(conf)
}

// Context から DataLoader を取得する
func For(ctx context.Context) *DataLoader {
	return ctx.Value(CtxKey).(*DataLoader)
}
