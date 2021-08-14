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

type indexedNavigation struct {
	Id        string    `json:"id"`
	Name      string    `json:"navigation"`
	TaskId    string    `json:"taskid"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexNavigation(ctx context.Context, navigation internal.Navigation) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Index")
	defer span.End()
	body := indexedNavigation{
		Id:        navigation.Id,
		Name:      navigation.Name,
		TaskId:    navigation.Task_id,
		CreatedAt: navigation.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_NAVIGATION,
		Body:       &buf,
		DocumentID: navigation.Id,
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

func (a *RBAC) DeleteNavigation(ctx context.Context, navigationId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_NAVIGATION,
		DocumentID: navigationId,
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

func (a *RBAC) GetNavigation(ctx context.Context, navigationId string) (internal.Navigation, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_NAVIGATION,
		DocumentID: navigationId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Navigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Navigation{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedNavigation `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Navigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	return internal.Navigation{
		Id:        hits.Source.Id,
		Name:      hits.Source.Name,
		Task_id:   hits.Source.TaskId,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) NavigationByTask(ctx context.Context, taskId *string) (internal.NavigationByTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Meny.ByTask")
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
		return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE_TASK},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.NavigationByTask{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedNavigation `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.NavigationByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	task := internal.Tasks{
		Id: *taskId,
	}
	res := make([]internal.Navigation, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].Name = hit.Source.Name
		res[i].Task_id = hit.Source.TaskId
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.NavigationByTask{
		Task:       task,
		Navigation: res,
	}, nil
}

func (a *RBAC) ListNavigation(ctx context.Context, args internal.ListArgs) (internal.ListNavigation, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_NAVIGATION},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListNavigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListNavigation{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedNavigation `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListNavigation{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.Navigation, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].Name = hit.Source.Name
		res[i].Task_id = hit.Source.TaskId
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListNavigation{
		Navigation: res,
		Total:      hits.Hits.Total.Value,
	}, nil
}
