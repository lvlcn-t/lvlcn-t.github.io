package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lvlcn-t/ChronoTemplify/pkg/components/layout"
	"github.com/lvlcn-t/ChronoTemplify/pkg/config"
	"github.com/lvlcn-t/ChronoTemplify/pkg/handlers"
)

const githubURL = "https://github.com/"

var _ handlers.Handler = (*projectsHandler)(nil)

type projectsHandler struct {
	log    *slog.Logger
	config *config.Data
}

func NewHandler(cfg *config.Data) handlers.Handler {
	return &projectsHandler{
		log:    slog.Default(),
		config: cfg,
	}
}

func (p *projectsHandler) View(c *gin.Context) {
	repos, err := p.listRepos(c.Request.Context())
	if err != nil {
		p.log.Error("Error while getting repositories", "error", err)
		c.HTML(http.StatusInternalServerError, "", "")
		return
	}

	repos = filterForks(repos)
	sortRepos(repos)

	page := layout.Layout(Projects(repos))
	c.HTML(http.StatusOK, "", page)
}

// Organization represents a GitHub organization
type Organization struct {
	// Name is the name of the organization
	Name string `json:"login"`
	// Description is the description of the organization
	Description string `json:"description"`
	// Url is the html url of the organization
	Url string `json:"html_url"`
	// ReposUrl is the url to get the repositories of the organization
	ReposUrl string `json:"repos_url"`
}

// Repository represents a GitHub repository
type Repository struct {
	// Name is the name of the repository
	Name string `json:"name"`
	// Organization is the name of the organization the repository belongs to (if any)
	Organization string
	// Url is the html url of the repository
	Url string `json:"html_url"`
	// Description is the description of the repository
	Description string `json:"description"`
	// Fork is a boolean indicating if the repository is a fork
	Fork bool `json:"fork"`
	// Stargazers is the number of stars the repository has
	Stargazers int `json:"stargazers_count"`
	// Language is the primary language of the repository
	Language string `json:"language"`
	// ContributersUrl is the url to get the contributors of the repository
	ContributersUrl string `json:"contributors_url"`
}

// listRepos gets the repositories of the configured GitHub profile
func (p *projectsHandler) listRepos(ctx context.Context) ([]Repository, error) {
	var profile string
	for _, social := range p.config.Socials {
		if strings.HasPrefix(social.Url, githubURL) {
			profile = strings.TrimPrefix(social.Url, githubURL)
		}
	}

	if profile == "" {
		return nil, fmt.Errorf("github profile not found")
	}

	repos, err := p.requestRepositories(ctx, fmt.Sprintf("https://api.github.com/users/%s/repos", profile))
	if err != nil {
		p.log.Error("Error while requesting repositories", "error", err)
		return nil, err
	}

	orgs, err := p.requestOrganizations(ctx, fmt.Sprintf("https://api.github.com/users/%s/orgs", profile))
	if err != nil {
		p.log.Error("Error while requesting organizations", "error", err)
		return nil, err
	}

	for _, org := range orgs {
		orgRepos, err := p.requestRepositories(ctx, org.ReposUrl)
		if err != nil {
			p.log.Error("Error while requesting repositories", "error", err)
			return nil, err
		}

		for _, repo := range orgRepos {
			if repo.IsUserContributor(ctx, profile) {
				repo.Organization = org.Name
				repos = append(repos, repo)
			}
		}
	}

	return repos, nil
}

// requestRepositories requests the repositories of the provided url
func (p *projectsHandler) requestRepositories(ctx context.Context, url string) ([]Repository, error) { //nolint:dupl // no need to refactor for now
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		p.log.Error("Error while creating request", "error", err)
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	res, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		p.log.Error("Error while requesting repositories", "error", err)
		return nil, err
	}
	defer func(b io.ReadCloser) {
		if err = b.Close(); err != nil {
			p.log.Error("Error while closing response body", "error", err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	var repos []Repository
	b, err := io.ReadAll(res.Body)
	if err != nil {
		p.log.Error("Error while reading response", "error", err)
		return nil, err
	}
	err = json.Unmarshal(b, &repos)
	if err != nil {
		p.log.Error("Error while unmarshaling response", "error", err)
		return nil, err
	}

	return repos, nil
}

// requestOrganizations requests the organizations the user is a member of
func (p *projectsHandler) requestOrganizations(ctx context.Context, url string) ([]Organization, error) { //nolint:dupl // no need to refactor for now
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		p.log.Error("Error while creating request", "error", err)
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	res, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		p.log.Error("Error while requesting organizations", "error", err)
		return nil, err
	}
	defer func(b io.ReadCloser) {
		if err = b.Close(); err != nil {
			p.log.Error("Error while closing response body", "error", err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	var orgs []Organization
	b, err := io.ReadAll(res.Body)
	if err != nil {
		p.log.Error("Error while reading response", "error", err)
		return nil, err
	}
	err = json.Unmarshal(b, &orgs)
	if err != nil {
		p.log.Error("Error while unmarshaling response", "error", err)
		return nil, err
	}

	return orgs, nil
}

// IsUserContributor checks if the user is a contributor to the repository
func (r *Repository) IsUserContributor(ctx context.Context, user string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.ContributersUrl, http.NoBody)
	if err != nil {
		return false
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	res, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		return false
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return false
	}

	var contributors []struct {
		Login string `json:"login"`
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return false
	}
	err = json.Unmarshal(b, &contributors)
	if err != nil {
		return false
	}

	for _, contributor := range contributors {
		if contributor.Login == user {
			return true
		}
	}

	return false
}

// filterForks filters out forked repositories
func filterForks(repos []Repository) []Repository {
	var filtered []Repository
	for _, repo := range repos {
		if !repo.Fork {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

// sortRepos sorts repositories by stars and then by name
func sortRepos(repos []Repository) {
	sort.Slice(repos, func(i, j int) bool {
		if repos[i].Stargazers == repos[j].Stargazers {
			return repos[i].Name < repos[j].Name
		}
		return repos[i].Stargazers > repos[j].Stargazers
	})
}
