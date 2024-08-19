package handler

import (
	"github.com/nzin/lexneo4j/swagger_gen/restapi/operations"
	"github.com/nzin/lexneo4j/swagger_gen/restapi/operations/app"
	"github.com/nzin/lexneo4j/swagger_gen/restapi/operations/health"
)

// Setup initialize all the handler functions
func Setup(api *operations.Lexneo4jAPI) {
	c := NewCRUD()

	// healthcheck
	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(c.GetHealthcheck)

	// neo4j functions
	api.AppListMoviesHandler = app.ListMoviesHandlerFunc(c.ListMovies)
	api.AppDoCypherHandler = app.DoCypherHandlerFunc(c.DoCypher)
}
