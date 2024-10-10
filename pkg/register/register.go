package register

import (
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/handlers"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/pages"
	"github.com/lvlcn-t/lvlcn-t.github.io/pkg/pages/projects"
)

var HandlerMap = map[string]handlers.HandlerConstructor{
	"/":         pages.NewHandler,
	"/projects": projects.NewHandler,
}
