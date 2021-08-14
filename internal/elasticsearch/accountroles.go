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

type indexedAccountRoles struct {
	Id              string    `json:"id"`
	AccountUsername string    `json:"account"`
	RoleId          string    `json:"role"`
	CreatedAt       time.Time `json:"createdat"`
}

func (a *RBAC) IndexAccountRole(ctx context.Context, accRole internal.AccountRoles) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Index")
	defer span.End()
	body := indexedAccountRoles{
		Id:              accRole.Id,
		AccountUsername: accRole.Account.UserName,
		RoleId:          accRole.Role.Id,
		CreatedAt:       accRole.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_ACCOUNT_ROLE,
		Body:       &buf,
		DocumentID: accRole.Id,
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

func (a *RBAC) DeleteAccountRole(ctx context.Context, accRoleId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_ACCOUNT_ROLE,
		DocumentID: accRoleId,
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

func (a *RBAC) GetAccountRole(ctx context.Context, accRoleId string) (internal.AccountRoles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_ACCOUNT_ROLE,
		DocumentID: accRoleId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.AccountRoles{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedAccountRoles `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.AccountRoles{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	account := internal.Account{
		UserName: hits.Source.AccountUsername,
	}
	role := internal.Roles{
		Id: hits.Source.RoleId,
	}
	return internal.AccountRoles{
		Id:        hits.Source.Id,
		Account:   account,
		Role:      role,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) AccountRoleByAccount(ctx context.Context, username *string) (internal.AccountRoleByAccountResult, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.ByAccount")
	defer span.End()

	should := make([]interface{}, 0, 4)
	if username != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"account": username,
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
		return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ACCOUNT_ROLE},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.AccountRoleByAccountResult{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedAccountRoles `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.AccountRoleByAccountResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	account := internal.Account{
		UserName: *username,
	}
	res := make([]internal.Roles, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {

		res[i].Id = hit.Source.RoleId
	}

	return internal.AccountRoleByAccountResult{
		Account: account,
		Roles:   res,
	}, nil
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) AccountRoleByRole(ctx context.Context, roleId *string) (internal.AccountRoleByRoleResult, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.ByRole")
	defer span.End()

	should := make([]interface{}, 0, 4)

	if roleId != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"role": roleId,
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
		return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ACCOUNT_ROLE},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.AccountRoleByRoleResult{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedAccountRoles `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.AccountRoleByRoleResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	role := internal.Roles{
		Id: *roleId,
	}
	res := make([]internal.Account, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {

		res[i].UserName = hit.Source.AccountUsername
	}

	return internal.AccountRoleByRoleResult{
		Role:    role,
		Account: res,
	}, nil
}

func (a *RBAC) ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_ACCOUNT_ROLE},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListAccountRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListAccountRole{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedAccountRoles `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListAccountRole{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.AccountRoles, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		account := internal.Account{
			UserName: hit.Source.AccountUsername,
		}
		role := internal.Roles{
			Id: hit.Source.RoleId,
		}
		res[i].Account = account
		res[i].Role = role
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListAccountRole{
		AccountRoles: res,
		Total:        hits.Hits.Total.Value,
	}, nil
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) AccountRoleByRoleReturnId(ctx context.Context, roleId string) ([]string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.ByRole")
	defer span.End()

	should := make([]interface{}, 0, 4)

	should = append(should, map[string]interface{}{
		"match": map[string]interface{}{
			"role": roleId,
		},
	})

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
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ACCOUNT_ROLE},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return nil, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedAccountRoles `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	var res []string

	for _, hit := range hits.Hits.Hits {

		res = append(res, hit.Source.Id)
	}

	return res, nil
}
