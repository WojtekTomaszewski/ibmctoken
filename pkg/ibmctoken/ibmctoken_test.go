package ibmctoken_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/WojtekTomaszewski/ibmctoken/pkg/ibmctoken"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func newMockClient(status int, body string, err error) *mockClient {
	if err != nil {
		return &mockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, err
			},
		}
	}

	return &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		},
	}
}

func TestNewToken(t *testing.T) {
	token := ibmctoken.NewToken("test-token")

	if token.ApiKey != "test-token" {
		t.Errorf("Api key was not set correctly")
	}
}

func TestRequestToken(t *testing.T) {
	// Initialize mock client

	t.Run("StatusCodeOK", func(t *testing.T) {
		token := ibmctoken.NewToken("test-token")

		mockBody := `{"access_token":"access-token","expires_in":3600,"expiration":1577808000}`
		token.Client = newMockClient(http.StatusOK, mockBody, nil)

		err := token.RequestToken()

		if err != nil {
			t.Errorf("Error was not expected: %v", err)
		}

		if token.AccessToken != "access-token" {
			t.Errorf("Access token was not set correctly: %s", token.AccessToken)
		}
	})

	t.Run("StatusCodeNotOK", func(t *testing.T) {
		token := ibmctoken.NewToken("test-token")

		mockBody := `{"errorMessage":"error-message","errorDetails":"error-details"}`
		token.Client = newMockClient(http.StatusBadRequest, mockBody, nil)

		err := token.RequestToken()

		if err == nil {
			t.Errorf("Error was expected")
		}

		if err.Error() != "error-message error-details" {
			t.Errorf("Different error message expected: %v", err)
		}

		if token.AccessToken != "" {
			t.Errorf("Access token was not expected: %+v", token)
		}
	})

	t.Run("StatusCodeOkBrokenBody", func(t *testing.T) {
		token := ibmctoken.NewToken("test-token")

		mockBody := `{"access_token":"access-token","expires_in":3600,"expiration":1577808000`
		token.Client = newMockClient(http.StatusOK, mockBody, nil)

		err := token.RequestToken()

		if err == nil {
			t.Errorf("Error was expected: %s", err)
		}

		if err.Error() != "unexpected EOF" {
			t.Errorf("Different error message expected: %v", err)
		}

		if token.AccessToken != "" {
			t.Errorf("Access token was not expected: %+v", token)
		}
	})

	t.Run("StatusCodeNotOkBrokenBody", func(t *testing.T) {
		token := ibmctoken.NewToken("test-token")

		mockBody := `{"errorMessage":"error-message","errorDetails":"error-details`
		token.Client = newMockClient(http.StatusBadRequest, mockBody, nil)

		err := token.RequestToken()

		if err == nil {
			t.Errorf("Error was expected")
		}

		if err.Error() != "unexpected EOF" {
			t.Errorf("Different error message expected: %v", err)
		}

		if token.AccessToken != "" {
			t.Errorf("Access token was not expected: %+v", token)
		}
	})

	t.Run("DoError", func(t *testing.T) {
		token := ibmctoken.NewToken("test-token")

		mockBody := `{"errorMessage":"error-message","errorDetails":"error-details`
		token.Client = newMockClient(http.StatusBadGateway, mockBody, fmt.Errorf("error"))

		err := token.RequestToken()

		if err == nil {
			t.Errorf("Error was expected")
		}

		if err.Error() != "error" {
			t.Errorf("Different error message expected: %v", err)
		}

		if token.AccessToken != "" {
			t.Errorf("Access token was not expected: %+v", token)
		}
	})

}

func TestExpired(t *testing.T) {
	token := ibmctoken.NewToken("test-token")
	token.Expiration = 1

	if !token.Expired() {
		t.Errorf("Token should be expired")
	}
}
