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

type indexedMenu struct {
	Id        string    `json:"id"`
	Name      string    `json:"menu"`
	TaskId    string    `json:"taskid"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexMenu(ctx context.Context, menu internal.Menu) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Index")
	defer span.End()
	body := indexedMenu{
		Id:        menu.Id,
		Name:      menu.Name,
		TaskId:    menu.Task_id,
		CreatedAt: menu.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_MENU,
		Body:       &buf,
		DocumentID: menu.Id,
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

func (a *RBAC) DeleteMenu(ctx context.Context, menuId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_MENU,
		DocumentID: menuId,
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

func (a *RBAC) GetMenu(ctx context.Context, menuId string) (internal.Menu, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_MENU,
		DocumentID: menuId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Menu{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedMenu `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Menu{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	return internal.Menu{
		Id:        hits.Source.Id,
		Name:      hits.Source.Name,
		Task_id:   hits.Source.TaskId,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) MenuByTask(ctx context.Context, taskId *string) (internal.MenuByTask, error) {
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
		return internal.MenuByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE_TASK},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.MenuByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.MenuByTask{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedMenu `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.MenuByTask{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	task := internal.Tasks{
		Id: *taskId,
	}
	res := make([]internal.Menu, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].Name = hit.Source.Name
		res[i].Task_id = hit.Source.TaskId
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.MenuByTask{
		Task: task,
		Menu: res,
	}, nil
}
