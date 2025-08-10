package delete_test

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/storage"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		expectedCode int
		respError    string
		mockError    error
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "NotFound Error",
			alias:        "nonexistent_alias",
			expectedCode: http.StatusNotFound,
			respError:    "URL not found",
			mockError:    storage.ErrURLNotFound,
		},
		{
			name:         "DeleteURL Error",
			alias:        "test_alias",
			expectedCode: http.StatusInternalServerError,
			respError:    "failed to delete url",
			mockError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.alias != "" {
				urlDeleterMock.
					On("DeleteURL", mock.Anything, tc.alias).
					Return(tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			r.Delete("/{alias}", delete.New(log, urlDeleterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			url := ts.URL + "/" + tc.alias

			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedCode, resp.StatusCode, string(body))

			if resp.StatusCode == http.StatusNoContent {
				require.Empty(t, body)
				return
			}

			var gotResp save.Response
			err = json.Unmarshal(body, &gotResp)
			require.NoError(t, err)
			require.Equal(t, tc.respError, gotResp.Error)
		})
	}
}
