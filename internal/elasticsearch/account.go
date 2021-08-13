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

type indexedAccount struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	ProfileId string    `json:"profileId"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"createdat"`
}

func (a *RBAC) IndexAccount(ctx context.Context, account internal.Account) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Index")
	defer span.End()
	body := indexedAccount{
		ID:        account.Id,
		Username:  account.UserName,
		ProfileId: account.Profile.Id,
		IsBlocked: account.IsBlocked,
		CreatedAt: account.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_ACCOUNT,
		Body:       &buf,
		DocumentID: account.UserName,
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

func (a *RBAC) DeleteAccount(ctx context.Context, username string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_ACCOUNT,
		DocumentID: username,
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

func (a *RBAC) GetAccount(ctx context.Context, username string) (internal.Account, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_ACCOUNT,
		DocumentID: username,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Account{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedAccount `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	profile := internal.Profile{
		Id: hits.Source.ProfileId,
	}
	return internal.Account{
		Id:        hits.Source.ID,
		UserName:  hits.Source.Username,
		Profile:   profile,
		IsBlocked: hits.Source.IsBlocked,
		CreatedAt: hits.Source.CreatedAt,
	}, err
}

// Search returns tasks matching a query.
// XXX: Pagination will be implemented in future episodes
func (a *RBAC) GetAccountById(ctx context.Context, id *string) (internal.Account, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Search")
	defer span.End()

	should := make([]interface{}, 0, 4)

	if id != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"id": *id,
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

	query["sort"] = []interface{}{
		map[string]interface{}{"_doc": "asc"},
	}
	fmt.Println(query)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.SearchRequest{
		Index: []string{INDEX_ACCOUNT},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Account{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedAccount `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Account{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.Account, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		profile := internal.Profile{
			Id: hit.Source.ProfileId,
		}
		res[i].Id = hit.Source.ID
		res[i].UserName = hit.Source.Username
		res[i].Profile = profile
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return res[0], nil
}
func (a *RBAC) ListAccount(ctx context.Context, args internal.ListAccountArgs) (internal.ListAccount, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.List")
	defer span.End()

	req := esv7api.SearchRequest{
		Index: []string{INDEX_ACCOUNT},
		Body:  strings.NewReader(`{"query":{"match_all": {}}}`),
		From:  args.From,
		Size:  args.Size,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.ListAccount{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.ListAccount{}, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedAccount `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.ListAccount{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]internal.Account, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.ID
		res[i].UserName = hit.Source.Username
		res[i].Profile.Id = hit.Source.ProfileId
		res[i].IsBlocked = hit.Source.IsBlocked
		res[i].CreatedAt = hit.Source.CreatedAt
	}

	return internal.ListAccount{
		Accounts: res,
		Total:    hits.Hits.Total.Value,
	}, nil
}
