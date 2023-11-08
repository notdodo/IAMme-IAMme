package neo4j

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jClient is an interface for interacting with the Neo4j database.
type Neo4jClient interface {
	Connect() neo4j.SessionWithContext
	Close() error
}

// Session is an interface for a Neo4j database session.
type Session interface {
	Run(cypher string, params map[string]interface{}) Result
	Close() error
}

// Result is an interface for a Neo4j query result.
type Result interface {
	Consume() (int, error)
}

type neo4jClient struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jClient(dbUri, username, password string) Neo4jClient {
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalln("Invalid Neo4j login", err.Error())
	}
	return &neo4jClient{
		driver: driver,
	}
}

func (c *neo4jClient) Connect() neo4j.SessionWithContext {
	return c.driver.NewSession(context.TODO(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
}

func (c *neo4jClient) Close() error {
	return c.driver.Close(context.TODO())
}
