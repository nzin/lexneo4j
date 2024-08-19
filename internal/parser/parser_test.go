package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	t.Run("test parse properties", func(t *testing.T) {
		s := "{foo:'bar'}"
		parser := NewParser(s)
		props, err := parser.parseProperties()
		assert.Nil(t, err)
		assert.Equal(t, "bar", props["foo"])
	})

	t.Run("test parse node definition 1", func(t *testing.T) {
		s := "()"
		parser := NewParser(s)
		node, err := parser.parseNode()
		assert.Nil(t, err)
		assert.Nil(t, node.VariableName)
	})
	t.Run("test parse node definition 2", func(t *testing.T) {
		s := "(n)"
		parser := NewParser(s)
		node, err := parser.parseNode()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.VariableName)
	})
	t.Run("test parse node definition 3", func(t *testing.T) {
		s := "(n:Person)"
		parser := NewParser(s)
		node, err := parser.parseNode()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.VariableName)
		assert.Equal(t, "Person", *node.TypeName)
	})
	t.Run("test parse node definition 4", func(t *testing.T) {
		s := "(n:Person{foo:'bar'})"
		parser := NewParser(s)
		node, err := parser.parseNode()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.VariableName)
		assert.Equal(t, "Person", *node.TypeName)
		assert.Equal(t, "bar", node.Props["foo"])
	})
	t.Run("test parse node definition 5", func(t *testing.T) {
		s := "(:Person{foo:'bar'})"
		parser := NewParser(s)
		node, err := parser.parseNode()
		assert.Nil(t, err)
		assert.Nil(t, node.VariableName)
		assert.Equal(t, "Person", *node.TypeName)
		assert.Equal(t, "bar", node.Props["foo"])
	})

	t.Run("test parse relationship definition 1", func(t *testing.T) {
		s := "[]"
		parser := NewParser(s)
		r, err := parser.parseRelationshipProperties()
		assert.Nil(t, err)
		assert.Nil(t, r.VariableName)
	})
	t.Run("test parse relationship definition 2", func(t *testing.T) {
		s := "[n]"
		parser := NewParser(s)
		r, err := parser.parseRelationshipProperties()
		assert.Nil(t, err)
		assert.Equal(t, "n", *r.VariableName)
	})
	t.Run("test parse relationship definition 3", func(t *testing.T) {
		s := "[n:Person]"
		parser := NewParser(s)
		r, err := parser.parseRelationshipProperties()
		assert.Nil(t, err)
		assert.Equal(t, "n", *r.VariableName)
		assert.Equal(t, "Person", *r.TypeName)
	})
	t.Run("test parse relationship definition 4", func(t *testing.T) {
		s := "[n:Person{foo:'bar'}]"
		parser := NewParser(s)
		r, err := parser.parseRelationshipProperties()
		assert.Nil(t, err)
		assert.Equal(t, "n", *r.VariableName)
		assert.Equal(t, "Person", *r.TypeName)
		assert.Equal(t, "bar", r.Props["foo"])
	})
	t.Run("test parse relationship definition 5", func(t *testing.T) {
		s := "[:Person{foo:'bar'}]"
		parser := NewParser(s)
		r, err := parser.parseRelationshipProperties()
		assert.Nil(t, err)
		assert.Nil(t, r.VariableName)
		assert.Equal(t, "Person", *r.TypeName)
		assert.Equal(t, "bar", r.Props["foo"])
	})

	t.Run("test return 1", func(t *testing.T) {
		s := "a"
		parser := NewParser(s)
		ret, err := parser.parseReturn()
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ret))
		assert.Equal(t, "a", ret[0].VariableName)
	})

	t.Run("test return 2", func(t *testing.T) {
		s := "a.foo,b"
		parser := NewParser(s)
		ret, err := parser.parseReturn()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(ret))
		assert.Equal(t, "a", ret[0].VariableName)
		assert.Equal(t, "foo", *ret[0].Property)
		assert.Equal(t, "b", ret[1].VariableName)
	})

	t.Run("complete test 1", func(t *testing.T) {
		s := "MATCH (n) RETURN n"
		parser := NewParser(s)
		node, err := parser.parseQuery()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.MatchNode.VariableName)
		assert.Equal(t, 1, len(node.Return))
		assert.Equal(t, "n", node.Return[0].VariableName)
	})
	t.Run("complete test 2", func(t *testing.T) {
		s := "MATCH (n:Person{foo:'bar'}) RETURN n.foo"
		parser := NewParser(s)
		node, err := parser.parseQuery()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.MatchNode.VariableName)
		assert.Equal(t, "bar", node.MatchNode.Props["foo"])
		assert.Equal(t, 1, len(node.Return))
		assert.Equal(t, "n", node.Return[0].VariableName)
		assert.Equal(t, "foo", *node.Return[0].Property)
	})
	t.Run("complete test 3", func(t *testing.T) {
		s := "MATCH (n:Person{foo:'bar'})-[r]->(o:Person) RETURN n.foo"
		parser := NewParser(s)
		node, err := parser.parseQuery()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.MatchNode.VariableName)
		assert.Equal(t, "bar", node.MatchNode.Props["foo"])
		assert.Equal(t, 1, len(node.Return))
		assert.Equal(t, "n", node.Return[0].VariableName)
		assert.Equal(t, "foo", *node.Return[0].Property)
		assert.Equal(t, "o", *node.Relationship.Target.VariableName)
		assert.Equal(t, "r", *node.Relationship.Props.VariableName)
		assert.Equal(t, REL_TO, node.Relationship.Direction)
	})
	t.Run("complete test 4", func(t *testing.T) {
		s := "MATCH (n:Person{foo:'bar'})<-[r{foo:'bar2'}]-(o:Person) RETURN n.foo"
		parser := NewParser(s)
		node, err := parser.parseQuery()
		assert.Nil(t, err)
		assert.Equal(t, "n", *node.MatchNode.VariableName)
		assert.Equal(t, "bar", node.MatchNode.Props["foo"])
		assert.Equal(t, 1, len(node.Return))
		assert.Equal(t, "n", node.Return[0].VariableName)
		assert.Equal(t, "foo", *node.Return[0].Property)
		assert.Equal(t, "o", *node.Relationship.Target.VariableName)
		assert.Equal(t, "r", *node.Relationship.Props.VariableName)
		assert.Equal(t, "bar2", node.Relationship.Props.Props["foo"])
		assert.Equal(t, REL_FROM, node.Relationship.Direction)
	})

	t.Run("not happy complete test 1", func(t *testing.T) {
		s := "MATCH (n) RETURN n,"
		parser := NewParser(s)
		_, err := parser.parseQuery()
		assert.NotNil(t, err)
	})
	t.Run("not happy complete test 2", func(t *testing.T) {
		s := "MATCH (n RETURN n"
		parser := NewParser(s)
		_, err := parser.parseQuery()
		assert.NotNil(t, err)
	})
	t.Run("not happy complete test 3", func(t *testing.T) {
		s := "MATCH (n:) RETURN n"
		parser := NewParser(s)
		_, err := parser.parseQuery()
		assert.NotNil(t, err)
	})
	t.Run("not happy complete test 4", func(t *testing.T) {
		s := "MATCH (n:)-[]]->(o) RETURN n"
		parser := NewParser(s)
		_, err := parser.parseQuery()
		assert.NotNil(t, err)
	})
}

func TestCypherReturn(t *testing.T) {

	t.Run("complete test 1", func(t *testing.T) {
		s := "MATCH (n) RETURN n"
		parser := NewParser(s)
		query, err := parser.parseQuery()
		assert.Nil(t, err)
		str := query.ToStringWithTenant("TENANT")
		assert.Equal(t, "MATCH (n{tenant:'TENANT'}) RETURN n", str)
	})
	t.Run("complete test 2", func(t *testing.T) {
		s := "MATCH (n:Person{foo:'bar'})-[r]->(o:Person) RETURN n.foo"
		parser := NewParser(s)
		query, err := parser.parseQuery()
		assert.Nil(t, err)
		str := query.ToStringWithTenant("TENANT")
		assert.Equal(t, "MATCH (n:Person{foo:'bar',tenant:'TENANT'})-[r{tenant:'TENANT'}]->(o:Person{tenant:'TENANT'}) RETURN n.foo", str)
	})
}
