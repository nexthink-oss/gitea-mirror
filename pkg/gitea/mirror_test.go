package gitea

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

type fakeSource struct{ token string }

func (s *fakeSource) GetType() string  { return "fake" }
func (s *fakeSource) GetToken() string { return s.token }
func (s *fakeSource) GetCloneURL(context.Context, *config.Repository) (string, error) {
	return "https://source.example/owner/repo.git", nil
}

type editRequest struct {
	path string
	auth string
	body map[string]any
}

func newTestController(t *testing.T, version string) (*Controller, *[]editRequest) {
	t.Helper()

	var edits []editRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/version":
			_ = json.NewEncoder(w).Encode(map[string]string{"version": version})
		case r.Method == http.MethodPatch:
			raw, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			body := map[string]any{}
			assert.NoError(t, json.Unmarshal(raw, &body))
			edits = append(edits, editRequest{path: r.URL.Path, auth: r.Header.Get("Authorization"), body: body})
			_ = json.NewEncoder(w).Encode(map[string]any{"name": "repo", "mirror": true})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c, err := NewController(&config.Target{Url: srv.URL, Token: "target-token"})
	require.NoError(t, err)

	return c, &edits
}

func testRepo(publicSource bool) *config.Repository {
	interval := 8 * time.Hour
	publicTarget := false
	return &config.Repository{
		Owner:        "owner",
		Name:         "repo",
		Interval:     &interval,
		PublicSource: &publicSource,
		PublicTarget: &publicTarget,
	}
}

func assertNoMirrorAuth(t *testing.T, body map[string]any) {
	t.Helper()
	assert.NotContains(t, body, "mirror_username")
	assert.NotContains(t, body, "mirror_password")
	assert.NotContains(t, body, "mirror_token")
}

func TestUpdateMirrorRotatesToken(t *testing.T) {
	c, edits := newTestController(t, "1.27.3")

	repo, err := c.UpdateMirror(t.Context(), &fakeSource{token: "new-token"}, testRepo(false))
	require.NoError(t, err)
	require.NotNil(t, repo)

	require.Len(t, *edits, 1)
	edit := (*edits)[0]
	assert.Equal(t, "/api/v1/repos/owner/repo", edit.path)
	assert.Equal(t, "token target-token", edit.auth)
	assert.Equal(t, "oauth2", edit.body["mirror_username"])
	assert.Equal(t, "new-token", edit.body["mirror_token"])
	assert.Equal(t, true, edit.body["private"])
	assert.Equal(t, "8h0m0s", edit.body["mirror_interval"])
}

func TestUpdateMirrorOldServer(t *testing.T) {
	c, edits := newTestController(t, "1.26.0")

	_, err := c.UpdateMirror(t.Context(), &fakeSource{token: "new-token"}, testRepo(false))
	var notUpdated *TokenNotUpdated
	require.ErrorAs(t, err, &notUpdated)

	require.Len(t, *edits, 1)
	assertNoMirrorAuth(t, (*edits)[0].body)
	assert.Equal(t, "8h0m0s", (*edits)[0].body["mirror_interval"])
}

func TestUpdateMirrorLeavesCredentials(t *testing.T) {
	tests := map[string]struct {
		source       *fakeSource
		publicSource bool
	}{
		"public source":    {source: &fakeSource{token: "new-token"}, publicSource: true},
		"skip credentials": {source: nil, publicSource: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c, edits := newTestController(t, "1.27.3")

			var err error
			if tt.source == nil {
				_, err = c.UpdateMirror(t.Context(), nil, testRepo(tt.publicSource))
			} else {
				_, err = c.UpdateMirror(t.Context(), tt.source, testRepo(tt.publicSource))
			}
			require.NoError(t, err)

			require.Len(t, *edits, 1)
			assertNoMirrorAuth(t, (*edits)[0].body)
			assert.Equal(t, true, (*edits)[0].body["private"])
		})
	}
}
