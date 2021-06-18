package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nipe0324/gqlgen-todos-example/dataloader"
	"github.com/nipe0324/gqlgen-todos-example/db"
	"github.com/nipe0324/gqlgen-todos-example/graph/generated"
	"github.com/nipe0324/gqlgen-todos-example/graph/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	currentUser, err := r.getCurrentUser()
	if err != nil {
		return nil, err
	}

	todo := db.Todo{
		Text:   input.Text,
		UserID: currentUser.ID,
	}

	err = r.conn.Create(&todo).Error
	if err != nil {
		return nil, err
	}

	return toGqlTodo(&todo), nil
}

func (r *mutationResolver) UploadTodoCsv(ctx context.Context, file graphql.Upload) ([]*model.Todo, error) {
	currentUser, err := r.getCurrentUser()
	if err != nil {
		return nil, err
	}

	var todos []*db.Todo

	// CSVファイルの読み込み
	reader := csv.NewReader(file.File)

	// ヘッダー行の読み込み（ヘッダーは無視するため）
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// ヘッダー以降の行の読み込み
	for {
		line, err := reader.Read()
		if err == io.EOF {
			// ファイルの読み込みが完了
			break
		} else if err != nil {
			// ファイル読み込みに失敗
			return nil, err
		}

		todo := &db.Todo{Text: line[0], UserID: currentUser.ID}
		todos = append(todos, todo)
	}

	// バルクインサート
	err = r.conn.Create(&todos).Error
	if err != nil {
		return nil, err
	}

	return toGqlTodos(todos), nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	currentUser, err := r.getCurrentUser()
	if err != nil {
		return nil, err
	}

	var todos []*db.Todo
	err = r.conn.Where(db.Todo{UserID: currentUser.ID}).Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return toGqlTodos(todos), nil
}

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	// データローダーからユーザ情報を取得
	user, err := dataloader.For(ctx).UserByID.Load(obj.UserID)
	if err != nil {
		return nil, err
	}

	return toGqlUser(user), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Pos(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
