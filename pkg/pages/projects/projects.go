package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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
	repos, err := p.getRepositories(c.Request.Context())
	if err != nil {
		p.log.Error("Error while getting repositories", "error", err)
		c.HTML(http.StatusInternalServerError, "", "")
		return
	}
	page := layout.Layout(Projects(repos))
	c.HTML(http.StatusOK, "", page)
}

type Repository struct {
	Name        string `json:"name"`
	Url         string `json:"html_url"`
	Description string `json:"description"`
}

func (p *projectsHandler) getRepositories(ctx context.Context) ([]Repository, error) {
	var profile string
	for _, social := range p.config.Socials {
		if strings.HasPrefix(social.Url, githubURL) {
			profile = strings.TrimPrefix(social.Url, githubURL)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.github.com/users/%s/repos", profile), http.NoBody)
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

	var repos []Repository
	if res.StatusCode == http.StatusOK {
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
	}

	return repos, nil
}
