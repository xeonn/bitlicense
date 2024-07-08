package fxserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/fx"
)

type CmdParam struct {
	Port int
}

var Content embed.FS

func StartFx(port int) {
	fx.New(
		fx.Provide(func() *CmdParam {
			return &CmdParam{
				Port: port,
			}
		}),
		// fx.Provide(config.NewAppConfigProducer),
		//
		// fx.Annotate is meant to be an alternative way of declaring values with
		// annotations without having to wrap them up in a separate structure.
		fx.Provide(NewHTTPServer, fx.Annotate(
			NewServerMux,
			fx.ParamTags(`group:"routes"`),
		),
			// Register REST Api here
			// AsRoute(handler.NewAppConfigHandler),
			// AsRoute(auth.NewAuthHandler),
			// AsRoute(transaction.NewFindAllTxnHandler),
			// AsRoute(transaction.NewFindByOrgTxnHandler),
		),
		fx.Invoke(func(*http.Server) {}),
		// fx.NopLogger,
	).Run()
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IRoute)),
		fx.ResultTags(`group:"routes"`),
	)
}

// handler declared by AsRoute above, must implement this interface
type IRoute interface {
	http.Handler

	Pattern() string
}

func NewHTTPServer(config *CmdParam, lc fx.Lifecycle, mux *mux.Router) *http.Server {
	// cors stuff
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodPut, http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodOptions},
		AllowCredentials: false,
	})

	svr := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: c.Handler(mux),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			conn, err := net.Listen("tcp", svr.Addr)
			if err != nil {
				return err
			}
			log.Println("Starting HTTP server at ", svr.Addr)
			go svr.Serve(conn)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return svr.Shutdown(ctx)
		},
	})
	return svr
}

func NewServerMux(routes []IRoute, config *CmdParam) *mux.Router {
	//mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true}))
	topmux := mux.NewRouter()

	// sort route by longest pattern first
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Pattern() < routes[j].Pattern()
	})

	for _, route := range routes {
		if strings.Contains(route.Pattern(), "/post") {
			topmux.Handle(route.Pattern(), route).Methods(http.MethodPost)
			continue
		}

		topmux.Handle(route.Pattern(), route)
	}

	// embedded file system
	fSys, err := fs.Sub(Content, "frontend")
	if err != nil {
		log.Fatalln(err)
	}

	fileserver := http.FileServer(http.FS(fSys))

	prefix := "/"
	topmux.PathPrefix(prefix).Handler(http.StripPrefix(prefix, fileserver))

	topmux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		paths, err := route.GetPathRegexp()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("PathTemplate ", paths)
		return nil
	})

	return topmux
}
