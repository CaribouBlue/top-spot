package mux

import (
	"net/http"

	"github.com/CaribouBlue/top-spot/internal/core"
	"github.com/CaribouBlue/top-spot/internal/server/middleware"
)

type RootMux struct {
	*http.ServeMux
	Services   RootMuxServices
	Middleware []middleware.Middleware
	Children   RootMuxChildren
}

type RootMuxServices struct {
	UserService *core.UserService
}

type RootMuxChildren struct {
	AuthMux   *AuthMux
	AppMux    *AppMux
	StaticMux *StaticMux
}

func NewRootMux(services RootMuxServices, middleware []middleware.Middleware, children RootMuxChildren) *RootMux {
	mux := &RootMux{
		http.NewServeMux(),
		services,
		middleware,
		children,
	}

	authPathPrefix := mux.Children.AuthMux.Opts.PathPrefix
	mux.Handle(authPathPrefix+"/", http.StripPrefix(authPathPrefix, mux.Children.AuthMux))

	appPathPrefix := mux.Children.AppMux.Opts.PathPrefix
	mux.Handle(appPathPrefix+"/", http.StripPrefix(appPathPrefix, mux.Children.AppMux))

	staticPathPrefix := mux.Children.StaticMux.Opts.PathPrefix
	mux.Handle(staticPathPrefix+"/", http.StripPrefix(staticPathPrefix, mux.Children.StaticMux))

	return mux
}

func (mux *RootMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.Apply(mux.ServeMux, mux.Middleware...).ServeHTTP(w, r)
}
