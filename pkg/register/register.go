package register

import (
	"github.com/lvlcn-t/ChronoTemplify/pkg/handlers"
	"github.com/lvlcn-t/ChronoTemplify/pkg/pages"
	"github.com/lvlcn-t/ChronoTemplify/pkg/pages/projects"
)

var HandlerMap = map[string]handlers.HandlerConstructor{
	"/":         pages.NewHandler,
	"/projects": projects.NewHandler,
}
