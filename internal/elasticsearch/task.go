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

type indexedTask struct {
	Id        string    `json:"id"`
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexTask(ctx context.Context, task internal.Tasks) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Index")
	defer span.End()
	body := indexedTask{
		Id:        task.Id,
		Task:      task.Task,
		CreatedAt: task.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_TASK,
		Body:       &buf,
		DocumentID: task.Id,
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

func (a *RBAC) DeleteTask(ctx context.Context, taskId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_TASK,
		DocumentID: taskId,
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

func (a *RBAC) GetTask(ctx context.Context, taskId string) (internal.Tasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_TASK,
		DocumentID: taskId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Tasks{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedTask `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Tasks{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	return internal.Tasks{
		Id:        hits.Source.Id,
		Task:      hits.Source.Task,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

func (a *RBAC) ListTask(ctx context.Context, args internal.ListArgs) (internal.ListTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_TASK},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListTask{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedTask `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.Tasks, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].Task = hit.Source.Task
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListTask{
		Task:  res,
		Total: hits.Hits.Total.Value,
	}, nil
}
