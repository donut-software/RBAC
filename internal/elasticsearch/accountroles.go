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
