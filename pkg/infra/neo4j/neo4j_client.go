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
	CreateNodes([]string, []map[string]interface{}) ([]map[string]interface{}, error)
	CreateRelationsAtoB([]string, []string, []string, []map[string]interface{}) ([]map[string]interface{}, error)
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
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.User_Id IS UNIQUE;", nil)
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS FOR (g:Group) REQUIRE g.Group_Id IS UNIQUE;", nil)
	session.Run(context.TODO(), "CREATE CONSTRAINT IF NOT EXISTS FOR (r:Rule) REQUIRE r.GroupRule_Id IS UNIQUE;", nil)
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

func (c *neo4jClient) CreateNodes(labels []string, properties []map[string]interface{}) ([]map[string]interface{}, error) {
	c.log.Debug("Creating new nodes", "count", len(properties), "params", properties)
	c.log.Info("Creating new nodes", "count", len(properties))
	nodeIDs, err := orm.CreateNodes(c.Connect(), labels, properties)
	if err != nil {
		c.log.Error("Failed creating nodes on Neo4J", "err", err)
	}
	c.log.Debug("Created nodes", "count", len(nodeIDs), "ids", nodeIDs)
	return nodeIDs, err
}

func (c *neo4jClient) CreateRelationsAtoB(labels []string, aLabels []string, bLabels []string, properties []map[string]interface{}) ([]map[string]interface{}, error) {
	c.log.Debug("Creating new relationships", "count", len(properties), "params", properties)
	c.log.Info("Creating new relationships", "count", len(properties))
	relIDs, err := orm.CreateRelationsAtoB(c.Connect(), labels, aLabels, bLabels, properties)
	if err != nil {
		c.log.Error("Failed creating relationships on Neo4J", "err", err)
	}
	c.log.Debug("Created relationships", "count", len(relIDs), "ids", relIDs)
	return relIDs, err
}
