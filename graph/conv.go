package graph

import (
	"fmt"

	"github.com/nipe0324/gqlgen-todos-example/db"
	"github.com/nipe0324/gqlgen-todos-example/graph/model"
)

// DBのTodoのリストをGraphQLのTodoのリストに変換
func toGqlTodos(t []*db.Todo) []*model.Todo {
	todos := make([]*model.Todo, len(t))

	for i, item := range t {
		todos[i] = toGqlTodo(item)
	}

	return todos
}

// DBのTodoをGraphQLのTodoに変換
func toGqlTodo(t *db.Todo) *model.Todo {
	if t == nil {
		return nil
	}

	return &model.Todo{
		ID:   fmt.Sprintf("%d", t.ID),
		Text: t.Text,
		Done: t.Done,
		User: toGqlUser(&t.User),
	}
}

// DBのUserをGraphQLのUserに変換
func toGqlUser(u *db.User) *model.User {
	if u == nil {
		return nil
	}

	return &model.User{
		ID:   fmt.Sprintf("%d", u.ID),
		Name: u.Name,
	}
}
