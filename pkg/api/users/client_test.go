package users

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seternate/go-lanty/pkg/api/models/user"
	"github.com/seternate/go-lanty/pkg/api/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClientInterface struct {
	buildURLErr           error
	buildURLResult        string
	newRequestErr         error
	newRequestResult      *http.Request
	doErr                 error
	doResult              *http.Response
	restRequestErr        error
	restRequestResult     *http.Response
	httpClient            *http.Client
	headers               map[string]string
	restRequestCallCount  int
	lastRESTRequestMethod string
	lastRESTRequestPath   paths.Path
	lastRESTRequestParams paths.Params
	lastRESTRequestHeaders map[string]string
	lastRESTRequestBody   []byte
}

func newMockClientInterface() *mockClientInterface {
	return &mockClientInterface{
		headers: make(map[string]string),
	}
}

func (m *mockClientInterface) BuildURL(path paths.Path, params paths.Params) (string, error) {
	if m.buildURLErr != nil {
		return "", m.buildURLErr
	}
	return m.buildURLResult, nil
}

func (m *mockClientInterface) NewRESTRequestWithContext(ctx context.Context, method string, url string, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Request, error) {
	if body != nil {
		m.lastRESTRequestBody, _ = io.ReadAll(body)
	}
	if m.newRequestErr != nil {
		return nil, m.newRequestErr
	}
	return m.newRequestResult, nil
}

func (m *mockClientInterface) Do(req *http.Request) (*http.Response, error) {
	if m.doErr != nil {
		return nil, m.doErr
	}
	return m.doResult, nil
}

func (m *mockClientInterface) RESTRequest(ctx context.Context, method string, path paths.Path, pathParams paths.Params, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Response, error) {
	m.restRequestCallCount++
	m.lastRESTRequestMethod = method
	m.lastRESTRequestPath = path
	m.lastRESTRequestParams = pathParams
	m.lastRESTRequestHeaders = headers
	if body != nil {
		m.lastRESTRequestBody, _ = io.ReadAll(body)
	}
	if m.restRequestErr != nil {
		return nil, m.restRequestErr
	}
	return m.restRequestResult, nil
}

func (m *mockClientInterface) HTTPClient() *http.Client {
	return m.httpClient
}

func (m *mockClientInterface) Headers() map[string]string {
	return m.headers
}

func TestNew(t *testing.T) {
	mockClient := newMockClientInterface()
	client := New(mockClient)
	
	assert.NotNil(t, client)
	assert.Equal(t, mockClient, client.apiClient)
}

func TestClient_GetUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedUsers := []*usermodel.User{
			{IPv4Address: "192.168.1.1", Username: "user1"},
			{IPv4Address: "192.168.1.2", Username: "user2"},
		}
		body, _ := json.Marshal(expectedUsers)
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(body)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(body))
		
		client := New(mockClient)
		users, err := client.GetUsers(context.Background())
		
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "192.168.1.1", users[0].IPv4Address)
		assert.Equal(t, "user1", users[0].Username)
		assert.Equal(t, "192.168.1.2", users[1].IPv4Address)
		assert.Equal(t, "user2", users[1].Username)
		assert.Equal(t, paths.GetUsers.Method, mockClient.lastRESTRequestMethod)
		assert.Equal(t, paths.GetUsers, mockClient.lastRESTRequestPath)
	})

	t.Run("API request error", func(t *testing.T) {
		mockClient := newMockClientInterface()
		mockClient.restRequestErr = assert.AnError
		
		client := New(mockClient)
		users, err := client.GetUsers(context.Background())
		
		require.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "failed API request")
	})

	t.Run("JSON decode error", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.WriteString("invalid json")
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(resp.Body)
		
		client := New(mockClient)
		users, err := client.GetUsers(context.Background())
		
		require.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "failed to decode response")
	})

	t.Run("empty list", func(t *testing.T) {
		body, _ := json.Marshal([]*usermodel.User{})
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(body)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(resp.Body)
		
		client := New(mockClient)
		users, err := client.GetUsers(context.Background())
		
		require.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 0)
	})
}

func TestClient_PutUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		reqBody := &usermodel.UpsertUserRequest{Username: "testuser"}
		expectedUser := &usermodel.User{IPv4Address: "192.168.1.1", Username: "testuser"}
		
		body, _ := json.Marshal(expectedUser)
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(body)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(resp.Body)
		
		client := New(mockClient)
		result, err := client.PutUser(context.Background(), "192.168.1.1", reqBody)
		
		require.NoError(t, err)
		assert.Equal(t, "192.168.1.1", result.IPv4Address)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, paths.PutUser.Method, mockClient.lastRESTRequestMethod)
		assert.Equal(t, paths.PutUser, mockClient.lastRESTRequestPath)
		assert.Equal(t, paths.UserParams{IPv4Address: "192.168.1.1"}, mockClient.lastRESTRequestParams)
		assert.Equal(t, "application/json", mockClient.lastRESTRequestHeaders["Content-Type"])
	})

	t.Run("JSON marshal error", func(t *testing.T) {
		// This test checks that API request errors are properly handled
		// The request struct is valid, so marshaling won't fail, but the API call will
		mockClient := newMockClientInterface()
		mockClient.restRequestErr = assert.AnError
		
		client := New(mockClient)
		reqBody := &usermodel.UpsertUserRequest{Username: "test"}
		result, err := client.PutUser(context.Background(), "192.168.1.1", reqBody)
		
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed API request")
	})

	t.Run("API request error", func(t *testing.T) {
		mockClient := newMockClientInterface()
		mockClient.restRequestErr = assert.AnError
		
		client := New(mockClient)
		reqBody := &usermodel.UpsertUserRequest{Username: "test"}
		result, err := client.PutUser(context.Background(), "192.168.1.1", reqBody)
		
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed API request")
	})

	t.Run("JSON decode error", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.WriteString("invalid json")
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(resp.Body)
		
		client := New(mockClient)
		reqBody := &usermodel.UpsertUserRequest{Username: "test"}
		result, err := client.PutUser(context.Background(), "192.168.1.1", reqBody)
		
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode response")
	})
}

func TestClient_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(resp.Body)
		
		client := New(mockClient)
		err := client.DeleteUser(context.Background(), "192.168.1.1")
		
		require.NoError(t, err)
		assert.Equal(t, paths.DeleteUser.Method, mockClient.lastRESTRequestMethod)
		assert.Equal(t, paths.DeleteUser, mockClient.lastRESTRequestPath)
		assert.Equal(t, paths.UserParams{IPv4Address: "192.168.1.1"}, mockClient.lastRESTRequestParams)
	})

	t.Run("API request error", func(t *testing.T) {
		mockClient := newMockClientInterface()
		mockClient.restRequestErr = assert.AnError
		
		client := New(mockClient)
		err := client.DeleteUser(context.Background(), "192.168.1.1")
		
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed API request")
	})
}
