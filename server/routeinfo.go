package server

import (
	"fmt"
	"net/http"
)

type RouteInfo struct {
	Summary     string
	Description string
	Tags        []string
	Middlewares []func(http.Handler) http.Handler
	Params      []OpenAPIParam
	Responses   map[string]OpenAPIResponse
}

type RouteOption func(*RouteInfo)

func WithResponse(responseCode int, response OpenAPIResponse) RouteOption {
	return func(ri *RouteInfo) {
		ri.Responses[fmt.Sprintf("%v", responseCode)] = response
	}
}

func WithSummary(summary string) RouteOption {
	return func(ri *RouteInfo) {
		ri.Summary = summary
	}
}

func WithDescription(description string) RouteOption {
	return func(ri *RouteInfo) {
		ri.Description = description
	}
}

func WithTags(tags ...string) RouteOption {
	return func(ri *RouteInfo) {
		ri.Tags = tags
	}
}

func WithParams(params ...OpenAPIParam) RouteOption {
	return func(ri *RouteInfo) {
		ri.Params = append(ri.Params, params...)
	}
}

func WithMiddlewares(mws ...func(http.Handler) http.Handler) RouteOption {
	return func(ri *RouteInfo) {
		ri.Middlewares = append(ri.Middlewares, mws...)
	}
}

func getRouteInfos(opts ...RouteOption) RouteInfo {
	routeInfo := initRouteInfo()
	for _, opt := range opts {
		opt(&routeInfo)
	}
	return routeInfo
}
