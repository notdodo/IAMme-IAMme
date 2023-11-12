package orm

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateNode creates nodes in Neo4j
func CreateNodes(session neo4j.SessionWithContext, labels []string, properties *[]map[string]interface{}) ([]map[string]interface{}, error) {
	ctx := context.TODO()
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		createNodeQuery := fmt.Sprintf("UNWIND $propsList AS props CREATE (n:%s) SET n += props RETURN id(n) as id", labelString(labels))
		parameters := map[string]interface{}{"propsList": filteredProperties(properties)}
		result, err := tx.Run(ctx, createNodeQuery, parameters)
		if err != nil {
			return nil, err
		}

		return collectResults(result, ctx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]map[string]interface{}), err
}

func filteredProperties(properties *[]map[string]interface{}) []map[string]interface{} {
	filteredProperties := make([]map[string]interface{}, len(*properties))
	for _, props := range *properties {
		filteredProp := make(map[string]interface{}, len(props))
		for key, value := range props {
			if isPrimitive(value) {
				filteredProp[key] = value
			}
		}
		filteredProperties = append(filteredProperties, filteredProp)
	}
	return filteredProperties
}

func collectResults(result neo4j.ResultWithContext, ctx context.Context) (results []map[string]interface{}, err error) {
	if records, err := result.Collect(ctx); err == nil {
		for _, k := range records {
			results = append(results, k.AsMap())
		}
	} else {
		return nil, err
	}
	return results, result.Err()
}

// isPrimitive checks if a value is a primitive type
func isPrimitive(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64, string, bool:
		return true
	default:
		return false
	}
}

// labelString generates a string representation of labels for use in Cypher queries.
func labelString(labels []string) string {
	return joinLabels(labels)
}

// joinLabels joins label strings with ":" separator.
func joinLabels(labels []string) string {
	return strings.Join(labels, ":")
}