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

type indexedProfile struct {
	ProfileID         string    `json:"profileId"`
	FirstName         string    `json:"firstname"`
	LastName          string    `json:"lastname"`
	ProfilePicture    string    `json:"profile_picture"`
	ProfileBackground string    `json:"profile_background"`
	Email             string    `json:"email"`
	Mobile            string    `json:"mobile"`
	CreatedAt         time.Time `json:"createdat"`
}

func (a *RBAC) IndexProfile(ctx context.Context, profile internal.Profile) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Index")
	defer span.End()
	body := indexedProfile{
		ProfileID:         profile.Id,
		FirstName:         profile.First_Name,
		LastName:          profile.Last_Name,
		ProfileBackground: profile.Profile_Background,
		ProfilePicture:    profile.Profile_Picture,
		Email:             profile.Email,
		Mobile:            profile.Mobile,
		CreatedAt:         profile.CreatedAt,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}
	req := esv7api.IndexRequest{
		Index:      INDEX_PROFILE,
		Body:       &buf,
		DocumentID: profile.Id,
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

func (a *RBAC) DeleteProfile(ctx context.Context, profileId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Delete")
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      INDEX_PROFILE,
		DocumentID: profileId,
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

func (a *RBAC) GetProfile(ctx context.Context, profileId string) (internal.Profile, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Get")
	defer span.End()
	req := esv7api.GetRequest{
		Index:      INDEX_PROFILE,
		DocumentID: profileId,
	}
	resp, err := req.Do(ctx, a.client)
	if err != nil {
		return internal.Profile{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "GetRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		fmt.Println(resp.String())
		return internal.Profile{}, internal.NewErrorf(internal.ErrorCodeUnknown, "GetRequest.Do %s", resp.StatusCode)
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//fmt.Println(string(body))

	var hits struct {
		Source indexedProfile `json:"_source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return internal.Profile{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}
	return internal.Profile{
		Id:                 hits.Source.ProfileID,
		Profile_Picture:    hits.Source.ProfilePicture,
		Profile_Background: hits.Source.ProfileBackground,
		First_Name:         hits.Source.FirstName,
		Last_Name:          hits.Source.LastName,
		Mobile:             hits.Source.Mobile,
		Email:              hits.Source.Email,
		CreatedAt:          hits.Source.CreatedAt,
	}, err
}
