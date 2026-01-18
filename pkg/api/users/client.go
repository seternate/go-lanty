package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/seternate/go-lanty/pkg/api/internal"
	"github.com/seternate/go-lanty/pkg/api/models/user"
	"github.com/seternate/go-lanty/pkg/api/paths"
)

type Client struct {
	apiClient internal.ClientInterface
}

func New(apiClient internal.ClientInterface) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

func (c *Client) GetUsers(ctx context.Context) ([]*usermodel.User, error) {
	resp, err := c.apiClient.RESTRequest(ctx, paths.GetUsers.Method, paths.GetUsers, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	var users []*usermodel.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return users, nil
}

func (c *Client) PutUser(ctx context.Context, ipv4Address string, req *usermodel.UpsertUserRequest) (*usermodel.User, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.apiClient.RESTRequest(ctx, paths.PutUser.Method, paths.PutUser, paths.UserParams{IPv4Address: ipv4Address}, map[string]string{"Content-Type": "application/json"}, nil, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	var userResp usermodel.User
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &userResp, nil
}

func (c *Client) DeleteUser(ctx context.Context, ipv4Address string) error {
	resp, err := c.apiClient.RESTRequest(ctx, paths.DeleteUser.Method, paths.DeleteUser, paths.UserParams{IPv4Address: ipv4Address}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
