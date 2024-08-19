package handler

import (
	"fmt"
	"time"

	"github.com/nzin/lexneo4j/internal/config"
	"github.com/nzin/lexneo4j/internal/parser"
	"github.com/nzin/lexneo4j/swagger_gen/models"
	"github.com/nzin/lexneo4j/swagger_gen/restapi/operations/app"
	"github.com/nzin/lexneo4j/swagger_gen/restapi/operations/health"
	"github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// CRUD is the CRUD interface
type CRUD interface {
	// healthcheck
	GetHealthcheck(health.GetHealthParams) middleware.Responder
	ListMovies(app.ListMoviesParams) middleware.Responder
	DoCypher(app.DoCypherParams) middleware.Responder
}

// NewCRUD creates a new CRUD instance
func NewCRUD() CRUD {
	neo4jdriver, err := neo4j.NewDriver(config.Config.Neo4jURL, neo4j.BasicAuth(config.Config.Neo4jUsername, config.Config.Neo4jPassword, ""), func(config *neo4j.Config) {
		config.MaxConnectionLifetime = 1 * time.Minute
		config.MaxConnectionPoolSize = 10
		config.ConnectionAcquisitionTimeout = time.Minute
		config.SocketKeepalive = true
	})
	if err != nil {
		panic(err)
	}

	return &crud{
		neo4jdriver: neo4jdriver,
	}
}

type crud struct {
	neo4jdriver neo4j.Driver
}

func (c *crud) GetHealthcheck(params health.GetHealthParams) middleware.Responder {
	return health.NewGetHealthOK().WithPayload(&models.Health{Status: "OK"})
}

func (c *crud) ListMovies(app.ListMoviesParams) middleware.Responder {
	session := c.neo4jdriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	movies, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		moviesList := make(map[string]int64)

		result, err := tx.Run("MATCH (m:Movie) RETURN m.title,m.released", nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			values := result.Record().Values
			if len(values) > 1 {
				title := values[0].(string)
				released := values[1].(int64)
				moviesList[title] = released
			}
		}

		if err = result.Err(); err != nil {
			logrus.Errorf("error reading result: %v", err)
			return nil, err
		}

		return moviesList, nil
	})
	if err != nil {
		return app.NewListMoviesDefault(500).WithPayload(
			ErrorMessage("cannot list movies: %v", err))
	}
	listMovies := []*models.Movie{}
	for title, released := range movies.(map[string]int64) {
		listMovies = append(listMovies, &models.Movie{
			Title:    &title,
			Released: released,
		})
	}

	return app.NewListMoviesOK().WithPayload(
		&app.ListMoviesOKBody{
			Movies: listMovies,
		},
	)
}

func (c *crud) DoCypher(params app.DoCypherParams) middleware.Responder {
	parser := parser.NewParser(params.Body.Cmd)
	query, err := parser.Parse()
	if err != nil {
		return app.NewDoCypherDefault(500).WithPayload(
			ErrorMessage("cannot parse query: %v", err))
	}

	session := c.neo4jdriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		resList := make([]map[string]interface{}, 0)

		result, err := tx.Run(query.ToString(), nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			record := result.Record()
			res := make(map[string]interface{})
			for _, k := range record.Keys {
				res[k], _ = record.Get(k)
			}
			resList = append(resList, res)
		}

		if err = result.Err(); err != nil {
			logrus.Errorf("error reading result: %v", err)
			return nil, err
		}

		return resList, nil
	})
	if err != nil {
		return app.NewListMoviesDefault(500).WithPayload(
			ErrorMessage("cannot list movies: %v", err))
	}

	results := make([]*app.DoCypherOKBodyResultItems0, 0)
	for _, r := range res.([]map[string]interface{}) {
		s := ""
		for k, v := range r {
			if s != "" {
				s += ","
			}
			s += fmt.Sprintf("%s:%v", k, v)
		}
	}

	return app.NewDoCypherOK().WithPayload(
		&app.DoCypherOKBody{
			Result: results,
		},
	)
}
