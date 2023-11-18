package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j/orm"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"
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
	log    logging.LogManager
}

/* #nosec */
//nolint:all
func (c *neo4jClient) setUpDb(session neo4j.SessionWithContext) {
	c.log.Info("Flushing the database")
	session.Run(context.TODO(), "MATCH (n) DETACH DELETE n;", nil)
	c.log.Info("Creating indexes")
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS ON (u:User) ASSERT u.Id IS UNIQUE;", nil)
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS ON (g:Group) ASSERT g.Id IS UNIQUE;", nil)
}

func NewNeo4jClient(dbUri, username, password string) Neo4jClient {
	logger := logging.GetLogManager()
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		logger.Error("Invalid Neo4j login", "err", err)
	}
	logger.Info("Valid Neo4j Connection")
	client := &neo4jClient{
		driver: driver,
		log:    logger,
	}
	client.setUpDb(client.Connect())
	return client
}

func (c *neo4jClient) Connect() neo4j.SessionWithContext {
	c.log.Info("Creating new Neo4j session")
	return c.driver.NewSession(context.TODO(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
}

func (c *neo4jClient) Close() error {
	return c.driver.Close(context.TODO())
}

func (c *neo4jClient) CreateNodes(labels []string, properties *[]map[string]interface{}) ([]map[string]interface{}, error) {
	c.log.Debug("Creating new nodes", "count", len(labels), "params", *properties)
	c.log.Info("Creating new nodes", "count", len(*properties))
	nodeIDs, err := orm.CreateNodes(c.Connect(), labels, properties)
	if err != nil {
		c.log.Error("Failed creating nodes on Neo4J", "err", err)
	}
	return nodeIDs, err
}
