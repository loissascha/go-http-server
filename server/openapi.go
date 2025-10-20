package server

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type OpenAPIDescription struct {
	Version    string                    `json:"openapi"`
	Info       OpenAPIDescriptionInfo    `json:"info"`
	Servers    []OpenAPIServer           `json:"servers"`
	Paths      map[string]map[string]any `json:"paths"`
	Components map[string]any            `json:"components"`
}

type OpenAPIServer struct {
	Url string `json:"url"`
}

type OpenAPIDescriptionInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type OpenAPIPath struct {
	Summary     string                     `json:"summary"`
	Description string                     `json:"description"`
	Tags        []string                   `json:"tags"`
	OperationId string                     `json:"operationId"`
	Responses   map[string]OpenAPIResponse `json:"responses"`
	Security    []map[string]any           `json:"security"`
}

type OpenAPIResponse struct {
	Description string `json:"description"`
}

type OpenAPIParam struct {
	Description string         `json:"description"`
	Name        string         `json:"name"`
	In          string         `json:"in"`
	Required    bool           `json:"required"`
	Schema      map[string]any `json:"schema"`
}

func getRouteParams(route string) []string {
	result := []string{}
	routeSplit := strings.SplitSeq(route, "")
	status := 0
	readingStatusName := ""
	for char := range routeSplit {
		if status == 0 {
			if char == "{" {
				status = 1
			}
		} else if status == 1 {
			if char == "}" {
				status = 0
				// fmt.Println("status name: ", readingStatusName)
				result = append(result, readingStatusName)
				readingStatusName = ""
			} else {
				readingStatusName += char
			}
		}
	}
	return result
}

func createRoutePaths(path *ServerPath, operationId string) map[string]any {
	routePaths := make(map[string]any)

	// read parameters out of route and if necessary create OpenAPIParam config
	routeParams := getRouteParams(path.Route)
	if len(routeParams) > 0 || len(path.Info.Params) > 0 {
		parameters := []OpenAPIParam{}
		for _, param := range routeParams {
			if slices.ContainsFunc(path.Info.Params, func(e OpenAPIParam) bool {
				return e.Name == param
			}) {
				continue
			}
			paramSchema := map[string]any{}
			paramSchema["type"] = "string"
			thisParam := OpenAPIParam{
				Description: "",
				Name:        param,
				In:          "path",
				Required:    true,
				Schema:      paramSchema,
			}
			parameters = append(parameters, thisParam)
		}
		for _, param := range path.Info.Params {
			thisParam := OpenAPIParam{
				Description: param.Description,
				Name:        param.Name,
				In:          param.In,
				Required:    param.Required,
				Schema:      param.Schema,
			}
			parameters = append(parameters, thisParam)
		}
		routePaths["parameters"] = parameters
	}

	switch path.Method {
	case METHOD_GET:
		routePaths["get"] = OpenAPIPath{
			Summary:     path.Info.Summary,
			Description: path.Info.Description,
			Tags:        path.Info.Tags,
			OperationId: operationId,
			Responses:   path.Info.Responses,
			Security:    []map[string]any{},
		}
	case METHOD_POST:
		routePaths["post"] = OpenAPIPath{
			Summary:     path.Info.Summary,
			Description: path.Info.Description,
			Tags:        path.Info.Tags,
			OperationId: operationId,
			Responses:   path.Info.Responses,
			Security:    []map[string]any{},
		}
	case METHOD_PUT:
		routePaths["put"] = OpenAPIPath{
			Summary:     path.Info.Summary,
			Description: path.Info.Description,
			Tags:        path.Info.Tags,
			OperationId: operationId,
			Responses:   path.Info.Responses,
			Security:    []map[string]any{},
		}
	case METHOD_DELETE:
		routePaths["delete"] = OpenAPIPath{
			Summary:     path.Info.Summary,
			Description: path.Info.Description,
			Tags:        path.Info.Tags,
			OperationId: operationId,
			Responses:   path.Info.Responses,
			Security:    []map[string]any{},
		}
	}
	return routePaths
}

func (s *Server) CreateOpenAPIJson(port string) {
	openApiObj := OpenAPIDescription{
		Version: "3.1.0",
		Info: OpenAPIDescriptionInfo{
			Title:   "Title Desc",
			Version: "1.0",
		},
		Servers:    []OpenAPIServer{},
		Paths:      make(map[string]map[string]any),
		Components: make(map[string]any),
	}

	openApiObj.Servers = append(openApiObj.Servers, OpenAPIServer{
		Url: "http://localhost:" + port,
	})

	operationId := 1
	parameters := make(map[string]OpenAPIParam)
	for _, paths := range s.Paths {
		for _, path := range paths {
			routePaths := createRoutePaths(&path, fmt.Sprintf("%v", operationId))
			openApiObj.Paths[path.Route] = routePaths
			operationId++
		}
	}
	openApiObj.Components["parameters"] = parameters

	m, err := json.MarshalIndent(openApiObj, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("openapi.json", m, 0777)
	if err != nil {
		panic(err)
	}
}
