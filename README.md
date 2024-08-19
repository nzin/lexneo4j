# lexneo4j

This is project is an example of how to create a REST API, that can talk to a Neo4j application.
In particular, how to expose a REST API to safely talk to Neo4j via sanitized CYPHER command.

To do so, the REST API server can receive a CYPHER command (on the /api/v1/cypher), will parse it to see if it is only a MATCH command, and execute it.

## How to run it
```
docker-compose up -d
make
curl http://localhost:18000/api/v1/cypher -H 'Content-type: application/json' -d '{"cmd":"MATCH (m:Movie) RETURN m.title,m.released"}' | jq .
```

## Parsing Cypher commands

To safely be able to execute CYPHER (readonly) commands, we parse the command via a lexer/parser. The code is in internal/parser directory

Usually you use it like

```
parser := NewParser(s)
query, err := parser.parseQuery()
assert.Nil(t, err)
str := query.ToString("TENANT")

```

Note: there is a tenant variant. This tenant variant can be used in a multi-tenant environment where all Neo4j have been tagged with a "TENANT" property. In that case you could use it like:
```
parser := NewParser(s)
query, err := parser.parseQuery()
assert.Nil(t, err)
str := query.ToStringWithTenant("TENANT")
```



