package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

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
	var user db.User
	err := r.conn.Where(db.User{ID: obj.UserID}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return toGqlUser(&user), nil
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
