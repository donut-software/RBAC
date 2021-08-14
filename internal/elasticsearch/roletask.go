package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"rbac/internal"
	"strings"
	"time"

	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"
	"go.opentelemetry.io/otel/trace"
)

type indexedRoleTask struct {
	Id        string    `json:"id"`
	TaskId    string    `json:"taskid"`
	RoleId    string    `json:"roleid"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexRoleTask(ctx context.Context, roletask internal.RoleTasks) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Index")
	defer span.End()
	body := indexedRoleTask{
		Id:        roletask.Id,
		TaskId:    roletask.Task.Id,
		RoleId:    roletask.Role.Id,
		CreatedAt: roletask.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_ROLE_TASK,
		Body:       &buf,
		DocumentID: roletask.Id,
		Refresh:    "true",
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "IndexRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "IndexRequest.Do %s", resp.StatusCode)
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func (a *RBAC) DeleteRoleTask(ctx context.Context, roletaskId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_ROLE_TASK,
		DocumentID: roletaskId,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "DeleteRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "DeleteRequest.Do %s", resp.StatusCode)
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func (a *RBAC) GetRoleTask(ctx context.Context, roletaskId string) (internal.RoleTasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_ROLE_TASK,
		DocumentID: roletaskId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.RoleTasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.RoleTasks{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedRoleTask `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.RoleTasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	task := internal.Tasks{
		Id: hits.Source.TaskId,
	}
	role := internal.Roles{
		Id: hits.Source.RoleId,
	}
	return internal.RoleTasks{
		Id:        hits.Source.Id,
		Task:      task,
		Role:      role,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) RoleTaskByRole(ctx context.Context, roleId *string) (internal.RoleTaskByRole, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.ByAccount")
	defer span.End()

	should := make([]interface{}, 0, 4)

	if roleId != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"roleid": roleId,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	fmt.Println(query)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE_TASK},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.RoleTaskByRole{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedRoleTask `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.RoleTaskByRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	role := internal.Roles{
		Id: *roleId,
	}
	res := make([]internal.Tasks, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {

		res[i].Id = hit.Source.TaskId
	}

	return internal.RoleTaskByRole{
		Role:  role,
		Tasks: res,
	}, nil
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) RoleTaskByTask(ctx context.Context, taskId *string) (internal.RoleTaskByTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.ByAccount")
	defer span.End()

	should := make([]interface{}, 0, 4)

	if taskId != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"taskid": taskId,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	fmt.Println(query)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE_TASK},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.RoleTaskByTask{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedRoleTask `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.RoleTaskByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	task := internal.Tasks{
		Id: *taskId,
	}
	res := make([]internal.Roles, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {

		res[i].Id = hit.Source.TaskId
	}

	return internal.RoleTaskByTask{
		Task:  task,
		Roles: res,
	}, nil
}

func (a *RBAC) ListRoleTask(ctx context.Context, args internal.ListArgs) (internal.ListRoleTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE_TASK},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListRoleTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListRoleTask{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedRoleTask `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListRoleTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.RoleTasks, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		task := internal.Tasks{
			Id: hit.Source.TaskId,
		}
		role := internal.Roles{
			Id: hit.Source.RoleId,
		}
		res[i].Id = hit.Source.Id
		res[i].Task = task
		res[i].Role = role
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListRoleTask{
		RoleTasks: res,
		Total:     hits.Hits.Total.Value,
	}, nil
}
