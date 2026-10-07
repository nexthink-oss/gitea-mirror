package gitea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitea.dev/sdk"
)

// TODO: drop this shim once EditRepoOption gains mirror credentials upstream
// (https://gitea.com/isometry/gitea-go-sdk/src/branch/feat/edit-repo-mirror-auth).
type editMirrorOption struct {
	gitea.EditRepoOption
	MirrorUsername *string `json:"mirror_username,omitempty"`
	MirrorToken    *string `json:"mirror_token,omitempty"`
}

func (c *Controller) editMirror(ctx context.Context, owner, name string, opt editMirrorOption) (*gitea.Repository, error) {
	body, err := json.Marshal(opt)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/api/v1/repos/%s/%s",
		strings.TrimSuffix(c.forge.GetUrl(), "/"), url.PathEscape(owner), url.PathEscape(name))

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.forge.GetToken())
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		var apiErr struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		return nil, fmt.Errorf("%s: %s", resp.Status, apiErr.Message)
	}

	repo := new(gitea.Repository)
	if err := json.NewDecoder(resp.Body).Decode(repo); err != nil {
		return nil, err
	}

	return repo, nil
}
