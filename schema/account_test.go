package schema

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/stretchr/testify/assert"
)

func TestApp_Init(t *testing.T) {
	initDB()
}

func TestApp_app(t *testing.T) {
	t.Skip()
	query := `
		query app($id: ID!) {
			app(id: $id) {
				id
				name
				icon
			}
		}
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"app": map[string]interface{}{
				"id":   "app-1",
				"name": "应用管理器",
				"icon": "icon",
			},
		},
		Errors: nil,
	}

	params := graphql.Params{
		Schema:        NewSchema(),
		RequestString: query,
		VariableValues: map[string]interface{}{
			"id": "app-1",
		},
	}
	result := graphql.Do(params)

	assert.Equal(t, []gqlerrors.FormattedError(nil), result.Errors)
	assert.Equal(t, expected, result)
}

func TestApp_ListApp(t *testing.T) {
	t.Skip()
	query := `
		query apps($status: String) {
			apps(status: $status) {
				totalCount
				edges{
					id
					name
					icon
				}
			}
		}
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"apps": map[string]interface{}{
				"totalCount": 1,
				"edges": []interface{}{
					map[string]interface{}{
						"id":   "app-1",
						"name": "应用管理器",
						"icon": "icon",
					},
				},
			},
		},
		Errors: nil,
	}

	params := graphql.Params{
		Schema:        NewSchema(),
		RequestString: query,
		VariableValues: map[string]interface{}{
			"status": "",
		},
	}
	result := graphql.Do(params)

	assert.Equal(t, []gqlerrors.FormattedError(nil), result.Errors)
	assert.Equal(t, expected, result)
}

func TestApp_Add(t *testing.T) {
	t.Skip()
	query := `
		mutation createApp($name: String!, $icon: String, $remark: String) {
			createApp(name: $name, icon: $icon, remark: $remark) {
				id
				name
				icon
			}
		}
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"createApp": map[string]interface{}{
				"id":   "app-2",
				"name": "测试应用",
				"icon": "",
			},
		},
		Errors: nil,
	}

	// ctx := context.WithValue(context.Background(), "session", &servermodel.Session{
	// 	ID:   uint(1),
	// 	Name: "test",
	// 	Role: "test",
	// })

	params := graphql.Params{
		// Context:       ctx,
		Schema:        NewSchema(),
		RequestString: query,
		VariableValues: map[string]interface{}{
			"name": "测试应用",
		},
	}
	result := graphql.Do(params)

	assert.Equal(t, []gqlerrors.FormattedError(nil), result.Errors)
	assert.Equal(t, expected, result)
}
