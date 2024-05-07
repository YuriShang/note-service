package user_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"note_service/app/internal/apperror"
	"note_service/app/pkg/logging"
	"note_service/app/pkg/rest"
	"time"
)

var _ UserClient = &client{}

type client struct {
	base     rest.BaseClient
	Resource string
}

func NewClient(baseURL string, resource string, logger logging.Logger) UserClient {
	c := client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
	return &c
}

type UserClient interface {
	GetUserByToken(ctx context.Context, token Token) (User, error)
}

func (c *client) GetUserByToken(ctx context.Context, t Token) (u User, err error) {
	bearerToken := fmt.Sprintf("%s %s", t.TokenType, t.AccessToken)
	c.base.Logger.Debug("add access_token to filter options")
	filters := []rest.FilterOptions{
		{
			//empty
		},
	}

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, filters)
	if err != nil {
		return u, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create new request due to error: %w", err)
	}
	req.Header.Set("Authorization", bearerToken)

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return u, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		if err = json.NewDecoder(response.Body()).Decode(&u); err != nil {
			return u, fmt.Errorf("failed to decode body due to error %w", err)
		}

		_, err1 := time.Parse("2006-01-02T15:04:05.999999", u.RegisterTime)
		_, err2 := time.Parse("2006-01-02T15:04:05.999999", u.PasswordSetTime)
		if err1 != nil || err2 != nil {
			return u, fmt.Errorf("failed to parse time field due to error %w", err)
		}
		return u, nil
	} else if response.StatusCode() == 401 {
		return u, fmt.Errorf("Unauthorized")
	}
	return u, apperror.APIError(response.Error.Message, response.Error.DeveloperMessage)
}
