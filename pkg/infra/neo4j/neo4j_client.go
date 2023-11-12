package neo4j

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j/orm"
)

// Neo4jClient is an interface for interacting with the Neo4j database.
type Neo4jClient interface {
	Connect() neo4j.SessionWithContext
	Close() error
	CreateNodes([]string, *[]map[string]interface{}) ([]map[string]interface{}, error)
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

/* #nosec */
//nolint:all
func (c *neo4jClient) setUpDb(session neo4j.SessionWithContext) {
	session.Run(context.TODO(), `MATCH (n) DETACH DELETE n;`, nil)
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS ON (u:User) ASSERT u.Id IS UNIQUE", nil)
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS ON (g:Group) ASSERT g.Id IS UNIQUE", nil)
}

func NewNeo4jClient(dbUri, username, password string) Neo4jClient {
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalln("Invalid Neo4j login", err.Error())
	}
	client := &neo4jClient{
		driver: driver,
	}
	client.setUpDb(client.Connect())
	return client
}

func (c *neo4jClient) Connect() neo4j.SessionWithContext {
	return c.driver.NewSession(context.TODO(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
}

func (c *neo4jClient) Close() error {
	return c.driver.Close(context.TODO())
}

func (c *neo4jClient) CreateNodes(labels []string, properties *[]map[string]interface{}) ([]map[string]interface{}, error) {
	nodeIDs, err := orm.CreateNodes(c.Connect(), []string{"User"}, properties)
	if err != nil {
		log.Fatalln(err)
	}
	return nodeIDs, err
}
