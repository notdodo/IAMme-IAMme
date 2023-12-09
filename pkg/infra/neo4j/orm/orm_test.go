package orm

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/notdodo/goflat/v2"
)

type User struct {
	Username string
	Email    string
}

type Member struct {
	User   *User
	Role   string
	Active bool
}

type Group struct {
	Name    string
	Members []*Member
}

func TestCreateNodesQuery(t *testing.T) {
	labels := []string{"A", "B"}
	members := []*Member{
		{User: &User{Username: "john_doe", Email: "john@example.com"}, Role: "Admin", Active: true},
		{User: &User{Username: "jane_doe", Email: "jane@example.com"}, Role: "User", Active: false},
	}
	group := Group{Name: "Admins", Members: members}
	flattenedMap := []map[string]interface{}{goflat.FlatStruct(group, goflat.FlattenerConfig{
		Prefix:    "",
		Separator: ".",
		SortKeys:  true,
		OmitEmpty: true,
	})}
	query, parameters := createNodesQuery(labels, flattenedMap)

	expectedQuery := fmt.Sprintf("UNWIND $propsList AS props CREATE (n:%s) SET n += props RETURN id(n) as id", strings.Join(labels, ":"))
	if expectedQuery != query {
		t.Errorf("expected: %s\ngot: %s", expectedQuery, query)
	}
	expectedParameters := map[string]interface{}{
		"propsList": []map[string]interface{}{
			{
				"Name":                    "Admins",
				"Members.0.User.Username": "john_doe",
				"Members.0.User.Email":    "john@example.com",
				"Members.0.Role":          "Admin",
				"Members.0.Active":        true,
				"Members.1.User.Username": "jane_doe",
				"Members.1.User.Email":    "jane@example.com",
				"Members.1.Role":          "User",
				"Members.1.Active":        false,
			},
		},
	}

	if !reflect.DeepEqual(expectedParameters, parameters) {
		t.Errorf("expected: %v\ngot: %v", expectedParameters, parameters)
	}
}

func TestCreateRelationsQuery(t *testing.T) {
	label := "RelLabel"
	aLabels := []string{"A"}
	bLabels := []string{"B"}
	properties := []map[string]interface{}{
		{
			"property": "1",
		},
		{
			"property": "2",
		},
	}
	query, parameters := createRelationsQuery(label, aLabels, bLabels, properties)

	expectedQuery := fmt.Sprintf(`UNWIND $propsList AS props MATCH (a:%s), (b:%s) WHERE a[props.left_key] = props.left_value AND b[props.right_key] = props.right_value CREATE (a)-[r:%s]->(b) SET r += apoc.map.fromPairs([[props.left_key, props.left_value], [props.right_key, props.right_value]]) RETURN id(r) as id`, strings.Join(aLabels, ":"), strings.Join(bLabels, ":"), label)
	if expectedQuery != query {
		t.Errorf("expected: %s\ngot: %s", expectedQuery, query)
	}
	expectedParameters := map[string]interface{}{
		"propsList": []map[string]interface{}{
			{
				"property": "1",
			},
			{
				"property": "2",
			},
		},
	}

	if !reflect.DeepEqual(expectedParameters, parameters) {
		t.Errorf("expected: %v\ngot: %v", expectedParameters, parameters)
	}
}
