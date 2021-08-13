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
