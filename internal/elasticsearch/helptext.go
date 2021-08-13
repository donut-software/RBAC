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
