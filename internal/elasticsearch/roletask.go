package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"rbac/internal"
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
