package dataloader_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nipe0324/gqlgen-todos-example/dataloader"
	"github.com/nipe0324/gqlgen-todos-example/db"
	"github.com/stretchr/testify/require"
)

func TestUserLoader(t *testing.T) {
	conf := dataloader.UserLoaderConfig{
		// Fetch はデータを取得する関数を実装する
		Fetch: func(keys []uint64) ([]*db.User, []error) {
			t.Logf("fetch users by keys: %+v", keys)

			users := make([]*db.User, len(keys))
			errors := make([]error, len(keys))

			// 偶数のユーザIDは成功、奇数のユーザIDはエラーにする
			// keys のインデックス値に該当のユーザやエラーを設定する
			for i, key := range keys {
				if key%2 == 0 {
					users[i] = &db.User{ID: key}
				} else {
					errors[i] = fmt.Errorf("not found user by id : %d", key)
				}
			}

			return users, errors
		},

		// Wait はデータを処理するまで待つ時間。時間が経過すればFetchが呼ばれる
		// 時間が長いほど処理が遅くなる。時間が短いほど複数回データ取得処理が実行されるトレードオフ
		Wait: 10 * time.Millisecond,

		// MaxBatch は１バッチあたりの最大のキーの数。0を指定した場合は制限なし
		MaxBatch: 10,
	}

	loader := dataloader.NewUserLoader(conf)

	// NOTE: 以下の行が追記
	t.Run("並列でデータを取得する", func(t *testing.T) {
		t.Run("ユーザの取得に成功", func(t *testing.T) {
			t.Parallel()
			u, err := loader.Load(2) // 奇数のため成功
			require.NoError(t, err)
			require.Equal(t, u.ID, uint64(2))
		})

		t.Run("ユーザの取得に失敗", func(t *testing.T) {
			t.Parallel()
			u, err := loader.Load(1) // 偶数のため失敗
			require.Error(t, err)
			require.Nil(t, u)
		})

		t.Run("複数のユーザを取得", func(t *testing.T) {
			t.Parallel()
			u, err := loader.LoadAll([]uint64{4, 3})
			require.Equal(t, u[0].ID, uint64(4))
			require.NoError(t, err[0])
			require.Nil(t, u[1])
			require.Error(t, err[1])
		})

		t.Run("WaitTime経過後にアクセス", func(t *testing.T) {
			t.Parallel()
			time.Sleep(20 * time.Millisecond)
			u, err := loader.Load(6)
			require.NoError(t, err)
			require.Equal(t, u.ID, uint64(6))
		})
	})
}
