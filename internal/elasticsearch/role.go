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

type indexedRole struct {
	Id        string    `json:"id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexRole(ctx context.Context, role internal.Roles) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Index")
	defer span.End()
	body := indexedRole{
		Id:        role.Id,
		Role:      role.Role,
		CreatedAt: role.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_ROLE,
		Body:       &buf,
		DocumentID: role.Id,
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

func (a *RBAC) DeleteRole(ctx context.Context, roleId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_ROLE,
		DocumentID: roleId,
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

func (a *RBAC) GetRole(ctx context.Context, roleId string) (internal.Roles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_ROLE,
		DocumentID: roleId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Roles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Roles{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedRole `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Roles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	return internal.Roles{
		Id:        hits.Source.Id,
		Role:      hits.Source.Role,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

func (a *RBAC) ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_ROLE},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListRole{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedRole `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.Roles, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].Role = hit.Source.Role
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListRole{
		Roles: res,
		Total: hits.Hits.Total.Value,
	}, nil
}
