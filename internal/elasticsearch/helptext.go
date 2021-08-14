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

type indexedHelpText struct {
	Id        string    `json:"id"`
	HelpText  string    `json:"helptext"`
	TaskId    string    `json:"taskid"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexHelpText(ctx context.Context, helptext internal.HelpText) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Index")
	defer span.End()
	body := indexedHelpText{
		Id:        helptext.Id,
		HelpText:  helptext.HelpText,
		TaskId:    helptext.Task_id,
		CreatedAt: helptext.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_HELPTEXT,
		Body:       &buf,
		DocumentID: helptext.Id,
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

func (a *RBAC) DeleteHelpText(ctx context.Context, helptextId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_HELPTEXT,
		DocumentID: helptextId,
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

func (a *RBAC) GetHelpText(ctx context.Context, helptextId string) (internal.HelpText, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_HELPTEXT,
		DocumentID: helptextId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.HelpText{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedHelpText `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.HelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	return internal.HelpText{
		Id:        hits.Source.Id,
		HelpText:  hits.Source.HelpText,
		Task_id:   hits.Source.TaskId,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) HelpTextByTask(ctx context.Context, taskId *string) (internal.HelpTextByTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountHelpTask.ByAccount")
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
		return internal.HelpTextByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE_TASK},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.HelpTextByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.HelpTextByTask{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedHelpText `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.HelpTextByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	task := internal.Tasks{
		Id: *taskId,
	}
	res := make([]internal.HelpText, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].HelpText = hit.Source.HelpText
		res[i].Task_id = hit.Source.TaskId
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.HelpTextByTask{
		Task:     task,
		HelpText: res[0],
	}, nil
}

func (a *RBAC) ListHelpText(ctx context.Context, args internal.ListArgs) (internal.ListHelpText, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_HELPTEXT},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListHelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListHelpText{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedHelpText `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListHelpText{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.HelpText, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].HelpText = hit.Source.HelpText
		res[i].Task_id = hit.Source.TaskId
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListHelpText{
		HelpText: res,
		Total:    hits.Hits.Total.Value,
	}, nil
}
