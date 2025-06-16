package config_test

import (
	"fmt"
	"iter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nexthink-oss/gitea-mirror/pkg/config"
)

func TestParseRepositorySpec(t *testing.T) {
	defaultOwner := "default-owner"
	tests := []struct {
		name          string
		spec          string
		defaultOwner  string
		expectedRepo  string
		expectedError bool
	}{
		{
			name:         "simple repo name with default owner",
			spec:         "repo-name",
			defaultOwner: defaultOwner,
			expectedRepo: "default-owner/repo-name",
		},
		{
			name:         "repo with explicit owner",
			spec:         "explicit-owner/repo-name",
			defaultOwner: defaultOwner,
			expectedRepo: "explicit-owner/repo-name",
		},
		{
			name:         "empty spec",
			spec:         "",
			defaultOwner: defaultOwner,
			expectedRepo: "default-owner/",
		},
		{
			name:         "just slash",
			spec:         "/",
			defaultOwner: defaultOwner,
			expectedRepo: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Defaults: config.Defaults{
					Owner: tt.defaultOwner,
				},
			}

			repo, err := cfg.ParseRepositorySpec(tt.spec)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRepo, repo)
			}
		})
	}
}

func TestRepositorySetEmptyOrContains(t *testing.T) {
	tests := []struct {
		name     string
		set      config.RepositorySet
		repo     config.Repository
		expected bool
	}{
		{
			name:     "empty set should return true",
			set:      config.RepositorySet{},
			repo:     config.Repository{Owner: "owner", Name: "repo"},
			expected: true,
		},
		{
			name:     "nil set should return true",
			set:      nil,
			repo:     config.Repository{Owner: "owner", Name: "repo"},
			expected: true,
		},
		{
			name: "set contains repo",
			set: config.RepositorySet{
				"owner/repo": struct{}{},
			},
			repo:     config.Repository{Owner: "owner", Name: "repo"},
			expected: true,
		},
		{
			name: "set does not contain repo",
			set: config.RepositorySet{
				"owner/other-repo": struct{}{},
			},
			repo:     config.Repository{Owner: "owner", Name: "repo"},
			expected: false,
		},
		{
			name: "set with multiple repos",
			set: config.RepositorySet{
				"owner1/repo1": struct{}{},
				"owner2/repo2": struct{}{},
			},
			repo:     config.Repository{Owner: "owner1", Name: "repo1"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.set.EmptyOrContains(tt.repo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRepositorySetFromArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		defaultOwner string
		expected     config.RepositorySet
	}{
		{
			name:         "empty args",
			args:         []string{},
			defaultOwner: "default-owner",
			expected:     config.RepositorySet{},
		},
		{
			name:         "single arg without owner",
			args:         []string{"repo1"},
			defaultOwner: "default-owner",
			expected: config.RepositorySet{
				"default-owner/repo1": struct{}{},
			},
		},
		{
			name:         "multiple args with and without owner",
			args:         []string{"repo1", "owner2/repo2"},
			defaultOwner: "default-owner",
			expected: config.RepositorySet{
				"default-owner/repo1": struct{}{},
				"owner2/repo2":        struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Defaults: config.Defaults{
					Owner: tt.defaultOwner,
				},
			}

			result := cfg.RepositorySetFromArgs(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilteredRepositories(t *testing.T) {
	// Helper function to collect repositories from an iterator
	collectRepos := func(seq iter.Seq[config.Repository]) []config.Repository {
		var result []config.Repository
		for repo := range seq {
			result = append(result, repo)
		}
		return result
	}

	tests := []struct {
		name         string
		repositories []config.Repository
		args         []string
		defaultOwner string
		expected     []config.Repository
	}{
		{
			name: "no args returns all repos",
			repositories: []config.Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
			},
			args:         []string{},
			defaultOwner: "default-owner",
			expected: []config.Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
			},
		},
		{
			name: "filter by specific repo",
			repositories: []config.Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
				{Owner: "owner3", Name: "repo3"},
			},
			args:         []string{"owner2/repo2"},
			defaultOwner: "default-owner",
			expected: []config.Repository{
				{Owner: "owner2", Name: "repo2"},
			},
		},
		{
			name: "filter by multiple repos",
			repositories: []config.Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
				{Owner: "owner3", Name: "repo3"},
			},
			args:         []string{"owner1/repo1", "owner3/repo3"},
			defaultOwner: "default-owner",
			expected: []config.Repository{
				{Owner: "owner1", Name: "repo1"},
				{Owner: "owner3", Name: "repo3"},
			},
		},
		{
			name: "default owner in args",
			repositories: []config.Repository{
				{Owner: "default-owner", Name: "repo1"},
				{Owner: "owner2", Name: "repo2"},
			},
			args:         []string{"repo1"},
			defaultOwner: "default-owner",
			expected: []config.Repository{
				{Owner: "default-owner", Name: "repo1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Defaults:     config.Defaults{Owner: tt.defaultOwner},
				Repositories: tt.repositories,
			}

			result := collectRepos(cfg.FilteredRepositories(tt.args))
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Helper to create a file for testing if it doesn't exist
	// os.MkdirAll("testdata", 0755)
	// if _, err := os.Stat("testdata/labels-config.yaml"); os.IsNotExist(err) {
	// 	content := []byte(`
	// source:
	//   type: gitea
	//   url: https://source-gitea.example.com
	// target:
	//   url: https://target-gitea.example.com
	// defaults:
	//   owner: default-org
	//   interval: 1h
	//   public-target: false
	//   labels: ["default-label", "team-infra"]
	// repositories:
	//   - name: repo-with-default-labels
	//   - name: repo-with-own-labels
	//     owner: custom-org
	//     labels: ["app-prod", "team-product"]
	//   - name: repo-with-empty-labels
	//     labels: []
	//   - name: repo-mixed-labels
	//     labels: ["default-label", "app-staging"]
	//   - name: repo-no-label-field
	// `)
	// 	_ = os.WriteFile("testdata/labels-config.yaml", content, 0644)
	// }

	tests := []struct {
		name          string
		configFiles   []string
		expectedError bool
		validate      func(*testing.T, *config.Config)
	}{
		{
			name:          "non-existent file",
			configFiles:   []string{"testdata/non-existent-file.yaml"},
			expectedError: true,
		},
		{
			name:        "valid config file",
			configFiles: []string{"testdata/valid-config.yaml"},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "gitea", cfg.Source.Type)
				assert.Equal(t, "https://source-gitea.example.com", cfg.Source.Url)
				assert.Equal(t, "https://remote-gitea.example.com", cfg.Source.RemoteUrl)
				assert.Equal(t, "https://target-gitea.example.com", cfg.Target.Url)
				assert.Equal(t, "test-org", cfg.Defaults.Owner)

				// Check interval is parsed correctly
				expectedInterval := 1 * time.Hour
				assert.Equal(t, expectedInterval, cfg.Defaults.Interval)

				// Check default settings
				assert.True(t, cfg.Defaults.PublicTarget)

				// Check repositories
				require.Equal(t, 2, len(cfg.Repositories))

				// First repo should inherit defaults
				assert.Equal(t, "test-org", cfg.Repositories[0].Owner)
				assert.Equal(t, "repo1", cfg.Repositories[0].Name)
				assert.NotNil(t, cfg.Repositories[0].Interval)
				assert.Equal(t, expectedInterval, *cfg.Repositories[0].Interval)
				assert.NotNil(t, cfg.Repositories[0].PublicTarget)
				assert.True(t, *cfg.Repositories[0].PublicTarget)

				// Second repo should have custom owner
				assert.Equal(t, "custom-org", cfg.Repositories[1].Owner)
				assert.Equal(t, "repo2", cfg.Repositories[1].Name)
				assert.NotNil(t, cfg.Repositories[1].Interval)
				assert.Equal(t, expectedInterval, *cfg.Repositories[1].Interval)
				assert.NotNil(t, cfg.Repositories[1].PublicTarget)
				assert.True(t, *cfg.Repositories[1].PublicTarget)
			},
		},
		{
			name:        "multiple config files with override",
			configFiles: []string{"testdata/base-config.yaml", "testdata/override-config.yaml"},
			validate: func(t *testing.T, cfg *config.Config) {
				// Base values that shouldn't change
				assert.Equal(t, "gitea", cfg.Source.Type)
				assert.Equal(t, "https://source-gitea.example.com", cfg.Source.Url)
				assert.Equal(t, "https://target-gitea.example.com", cfg.Target.Url)
				assert.Equal(t, "test-org", cfg.Defaults.Owner)

				// Check repositories were merged
				require.Equal(t, 3, len(cfg.Repositories))

				// First two repos from override-config.yaml
				assert.Equal(t, "test-org", cfg.Repositories[0].Owner) // default owner
				assert.Equal(t, "repo1", cfg.Repositories[0].Name)

				assert.Equal(t, "test-org", cfg.Repositories[1].Owner) // default owner
				assert.Equal(t, "repo2", cfg.Repositories[1].Name)

				// Third repo with custom owner
				assert.Equal(t, "different-org", cfg.Repositories[2].Owner)
				assert.Equal(t, "repo3", cfg.Repositories[2].Name)

				// All repos should have defaults applied
				expectedInterval := 1 * time.Hour
				for _, repo := range cfg.Repositories {
					assert.NotNil(t, repo.Interval)
					assert.Equal(t, expectedInterval, *repo.Interval)
				}
			},
		},
		{
			name:        "labels config file",
			configFiles: []string{"testdata/labels-config.yaml"},
			validate: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "default-org", cfg.Defaults.Owner)
				assert.Equal(t, []string{"default-label", "team-infra"}, cfg.Defaults.Labels)

				require.Len(t, cfg.Repositories, 5)

				// repo-with-default-labels
				assert.Equal(t, "repo-with-default-labels", cfg.Repositories[0].Name)
				assert.Equal(t, []string{"default-label", "team-infra"}, cfg.Repositories[0].Labels, "repo-with-default-labels should inherit default labels")

				// repo-with-own-labels
				assert.Equal(t, "repo-with-own-labels", cfg.Repositories[1].Name)
				assert.Equal(t, "custom-org", cfg.Repositories[1].Owner)
				assert.Equal(t, []string{"app-prod", "team-product"}, cfg.Repositories[1].Labels)

				// repo-with-empty-labels
				assert.Equal(t, "repo-with-empty-labels", cfg.Repositories[2].Name)
				assert.Empty(t, cfg.Repositories[2].Labels, "repo-with-empty-labels should have no labels")

				// repo-mixed-labels
				assert.Equal(t, "repo-mixed-labels", cfg.Repositories[3].Name)
				assert.Equal(t, []string{"default-label", "app-staging"}, cfg.Repositories[3].Labels)

				// repo-no-label-field
				assert.Equal(t, "repo-no-label-field", cfg.Repositories[4].Name)
				assert.Equal(t, []string{"default-label", "team-infra"}, cfg.Repositories[4].Labels, "repo-no-label-field should inherit default labels")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadConfig(tt.configFiles)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	// Load the valid config file which has defaults and repositories with some fields missing
	cfg, err := config.LoadConfig([]string{"testdata/valid-config.yaml"})
	require.NoError(t, err)

	// Verify defaults were applied correctly
	expectedInterval := 1 * time.Hour

	// Check first repository (should use all defaults)
	assert.Equal(t, "test-org", cfg.Repositories[0].Owner)
	assert.Equal(t, "repo1", cfg.Repositories[0].Name)
	assert.NotNil(t, cfg.Repositories[0].Interval)
	assert.Equal(t, expectedInterval, *cfg.Repositories[0].Interval)
	assert.NotNil(t, cfg.Repositories[0].PublicTarget)
	assert.True(t, *cfg.Repositories[0].PublicTarget)

	// Check second repository (has custom owner, should inherit interval and public)
	assert.Equal(t, "custom-org", cfg.Repositories[1].Owner)
	assert.Equal(t, "repo2", cfg.Repositories[1].Name)
	assert.NotNil(t, cfg.Repositories[1].Interval)
	assert.Equal(t, expectedInterval, *cfg.Repositories[1].Interval)
	assert.NotNil(t, cfg.Repositories[1].PublicTarget)
	assert.True(t, *cfg.Repositories[1].PublicTarget)
}

func TestForgeInterfaces(t *testing.T) {
	t.Run("Source implements Forge interface", func(t *testing.T) {
		source := &config.Source{
			Type:      "gitea",
			Url:       "https://source.example.com",
			RemoteUrl: "https://remote.example.com",
			Token:     "source-token",
		}

		assert.Equal(t, "gitea", source.GetType())
		assert.Equal(t, "https://source.example.com", source.GetUrl())
		assert.Equal(t, "https://remote.example.com", source.GetRemoteUrl())
		assert.Equal(t, "source-token", source.GetToken())
	})

	t.Run("Target implements Forge interface", func(t *testing.T) {
		target := &config.Target{
			Url:   "https://target.example.com",
			Token: "target-token",
		}

		assert.Equal(t, "gitea", target.GetType())
		assert.Equal(t, "https://target.example.com", target.GetUrl())
		assert.Equal(t, "", target.GetRemoteUrl())
		assert.Equal(t, "target-token", target.GetToken())
	})
}

func TestLabelledRepositories(t *testing.T) {
	defaultLabels := []string{"default", "all"}
	repo1 := config.Repository{Name: "repo1", Owner: "test", Labels: []string{"feat-a", "team-x"}}
	repo2 := config.Repository{Name: "repo2", Owner: "test", Labels: []string{"feat-b", "team-y"}}
	repo3 := config.Repository{Name: "repo3", Owner: "test", Labels: []string{"feat-a", "team-y"}}
	repo4 := config.Repository{Name: "repo4", Owner: "test", Labels: []string{}}
	repo5 := config.Repository{Name: "repo5", Owner: "test"} // No labels field, will take default
	repo6 := config.Repository{Name: "repo6", Owner: "test", Labels: defaultLabels}

	// The LabelledRepositories method operates on a Config instance, but primarily uses
	// its Repositories field. The Defaults.Labels are used by LoadConfig to populate
	// repository labels if they are missing. For testing LabelledRepositories in isolation,
	// we need to simulate the state of Repositories *after* LoadConfig would have run.

	// Apply defaults to repo5 for standalone testing of LabelledRepositories
	// In LoadConfig, this happens: if len(r.Labels) == 0 { config.Repositories[i].Labels = config.Defaults.Labels }
	// repo5 has len(.Labels) == 0 initially, so it would get defaultLabels.
	// repo4 has len(.Labels) == 0, but it's explicitly `[]`, so it should remain empty.
	// The logic in LoadConfig is: `if len(r.Labels) == 0 { r.Labels = defaults.Labels }`
	// This means repo4 (labels: []) would also get default labels if not careful. The actual code is:
	// for i, r := range config.Repositories {
	// 	 if len(r.Labels) == 0 { // This applies to both nil and empty slice if key was missing or explicitly empty in YAML
	// 		 config.Repositories[i].Labels = config.Defaults.Labels
	// 	 }
	// }
	// So, for repo4 (labels: []), it *will* inherit default labels during LoadConfig.
	// And repo5 (no labels field) will also inherit.
	// This means the test data for LabelledRepositories needs to reflect post-LoadConfig state.

	// Simulating post-LoadConfig state for labels:
	processedRepo4 := repo4
	processedRepo4.Labels = defaultLabels // repo4 with `labels: []` gets defaults
	processedRepo5 := repo5
	processedRepo5.Labels = defaultLabels // repo5 with no `labels` field gets defaults

	cfgForLabelTest := &config.Config{
		Repositories: []config.Repository{
			repo1,              // {feat-a, team-x}
			repo2,              // {feat-b, team-y}
			repo3,              // {feat-a, team-y}
			processedRepo4,     // {default, all} (originally empty, but inherits)
			processedRepo5,     // {default, all} (originally nil, inherits)
			repo6,              // {default, all}
		},
	}

	tests := []struct {
		name          string
		filterLabels  []string
		expectedRepos []config.Repository
	}{
		{
			name:         "filter by 'feat-a'",
			filterLabels: []string{"feat-a"},
			expectedRepos: []config.Repository{repo1, repo3},
		},
		{
			name:         "filter by 'team-y'",
			filterLabels: []string{"team-y"},
			expectedRepos: []config.Repository{repo2, repo3},
		},
		{
			name:         "filter by 'feat-a' OR 'feat-b'",
			filterLabels: []string{"feat-a", "feat-b"},
			expectedRepos: []config.Repository{repo1, repo2, repo3},
		},
		{
			name:         "filter by 'non-existent-label'",
			filterLabels: []string{"non-existent-label"},
			expectedRepos: []config.Repository{},
		},
		{
			name:         "filter by 'default' label",
			filterLabels: []string{"default"},
			expectedRepos: []config.Repository{processedRepo4, processedRepo5, repo6},
		},
		{
			name:         "filter by 'all' label",
			filterLabels: []string{"all"},
			expectedRepos: []config.Repository{processedRepo4, processedRepo5, repo6},
		},
		{
			name: "filter by specific and default label, e.g., 'feat-a' and 'all'",
			filterLabels: []string{"feat-a", "all"},
			expectedRepos: []config.Repository{repo1, repo3, processedRepo4, processedRepo5, repo6},
		},
		// As per root.go, LabelledRepositories is only called if filterLabels is not empty.
		// So, a test case for empty filterLabels might not be strictly necessary for CLI behavior
		// but for completeness of the function itself:
		// Based on current LabelledRepositories code: `if len(labels) == 0 || len(repo.Labels) > 0`
		// if filterLabels is empty, `len(labels)==0` is true. The inner loop `for _, label := range labels` does not run.
		// So, `LabelledRepositories([])` will return `[]`.
		{
			name:         "empty filter labels (direct call)",
			filterLabels: []string{},
			expectedRepos: []config.Repository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfgForLabelTest.LabelledRepositories(tt.filterLabels)
			assert.ElementsMatch(t, tt.expectedRepos, result)
		})
	}
}

func TestRepositoryStringMethods(t *testing.T) {
	repo := config.Repository{
		Owner: "test-owner",
		Name:  "test-repo",
	}

	t.Run("String()", func(t *testing.T) {
		assert.Equal(t, "test-owner/test-repo", repo.String())
	})

	t.Run("Success() without message", func(t *testing.T) {
		assert.Equal(t, "✅ test-owner/test-repo", repo.Success())
	})

	t.Run("Success() with message", func(t *testing.T) {
		assert.Equal(t, "✅ test-owner/test-repo: Operation complete", repo.Success("Operation complete"))
	})

	t.Run("Success() with multiple messages", func(t *testing.T) {
		assert.Equal(t, "✅ test-owner/test-repo: First message Second message",
			repo.Success("First message", "Second message"))
	})

	t.Run("Failure()", func(t *testing.T) {
		err := fmt.Errorf("something went wrong")
		assert.Equal(t, "❌ test-owner/test-repo: something went wrong", repo.Failure(err))
	})
}
