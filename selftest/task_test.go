package selftest

import (
	"fmt"
	"testing"
	"time"

	"github.com/deliveroo/assert-go"
	"github.com/icrowley/fake"
)

func TestCreateGetUpdateDeleteTask(t *testing.T) {
	withAccount(t, func(api *API) {
		var (
			id          interface{}
			description = fake.Sentence()
		)
		t.Run("create", func(t *testing.T) {
			resp := api.Post(t, "/tasks", m{
				"description": description,
			})
			resp.AssertStatusCode(t, 200)
			resp.JSONPathEqual(t, "description", description)
			resp.JSONPathEqual(t, "completed", nil)
			id = resp.JSONPath(t, "id")
		})
		t.Run("get", func(t *testing.T) {
			resp := api.Get(t, "/tasks/"+fmt.Sprint(id))
			resp.AssertStatusCode(t, 200)
			resp.JSONPathEqual(t, "description", description)
			resp.JSONPathEqual(t, "completed", nil)
		})
		t.Run("update", func(t *testing.T) {
			suffix := fake.Word()
			when := time.Now().UTC().Format(time.RFC3339)
			resp := api.Put(t, "/tasks/"+fmt.Sprint(id), m{
				"description": description + suffix,
				"completed":   when,
			})
			resp.AssertStatusCode(t, 200)
			resp.JSONPathEqual(t, "description", description+suffix)
			resp.JSONPathEqual(t, "completed", when)
			assert.NotNil(t, resp.JSONPath(t, "completed"))
		})
		t.Run("delete", func(t *testing.T) {
			resp := api.Delete(t, "/tasks/"+fmt.Sprint(id), nil)
			resp.AssertStatusCode(t, 200)
			api.Get(t, "/tasks/"+fmt.Sprint(id)).AssertStatusCode(t, 404)
		})
	})
}

func TestTaskValidation(t *testing.T) {
	withAccount(t, func(api *API) {
		resp := api.Post(t, "/tasks", m{
			"description": "",
		})
		resp.AssertStatusCode(t, 400)
		assert.Equal(t, resp.ErrorMessage(t), "description is required")
	})
}

func TestGetAllTasks(t *testing.T) {
	withAccount(t, func(api *API) {
		expected := make(map[string]bool)
		t.Run("create 10 tasks", func(t *testing.T) {
			for i := 1; i <= 10; i++ {
				description := fmt.Sprintf("task %d", i)
				expected[description] = true
				resp := api.Post(t, "/tasks", m{
					"description": description,
				})
				resp.AssertStatusCode(t, 200)
			}
		})
		t.Run("get all", func(t *testing.T) {
			resp := api.Get(t, "/tasks")
			resp.AssertStatusCode(t, 200)
			type task struct {
				Description string `json:"description"`
			}
			var tasks []task
			resp.BindBody(t, &tasks)
			assert.Equal(t, len(tasks), 10)
			for _, tt := range tasks {
				_, ok := expected[tt.Description]
				assert.True(t, ok)
				delete(expected, tt.Description)
			}
			assert.Equal(t, len(expected), 0)
		})
	})
}

func TestGetIncompleteTasks(t *testing.T) {
	withAccount(t, func(api *API) {
		t.Run("create 10 tasks", func(t *testing.T) {
			for i := 1; i <= 10; i++ {
				var completed *string
				if i%2 == 0 {
					now := time.Now().UTC().Format(time.RFC3339)
					completed = &now
				}
				description := fmt.Sprintf("task %d", i)
				resp := api.Post(t, "/tasks", m{
					"description": description,
					"completed":   completed,
				})
				resp.AssertStatusCode(t, 200)
			}
		})
		t.Run("get all incomplete", func(t *testing.T) {
			resp := api.Get(t, "/tasks?filter=incomplete")
			resp.AssertStatusCode(t, 200)
			type task struct {
				Completed *string `json:"completed"`
			}
			var tasks []task
			resp.BindBody(t, &tasks)
			assert.Equal(t, len(tasks), 5)
			for _, tt := range tasks {
				assert.Nil(t, tt.Completed)
			}
		})
	})
}
