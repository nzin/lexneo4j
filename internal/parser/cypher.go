package parser

import "fmt"

const (
	REL_BOTH int = iota
	REL_TO
	REL_FROM
)

type CypherQuery struct {
	MatchNode    CypherNode
	Relationship *CypherRelationShip
	Return       CypherReturn
}

type CypherRelationShip struct {
	Direction int
	Props     *CypherNode
	Target    CypherNode
}

type CypherNode struct {
	VariableName *string
	TypeName     *string
	Props        map[string]string
}

type CypherReturn []CypherVariableReturn

type CypherVariableReturn struct {
	VariableName string
	Property     *string
}

func (r *CypherVariableReturn) ToString() string {
	if r.Property == nil {
		return r.VariableName
	} else {
		return fmt.Sprintf("%s.%s", r.VariableName, *r.Property)
	}
}

// func (n *CypherNode) ToString(tenant string) string {
// 	str := ""
// 	if n.VariableName != nil {
// 		str = *n.VariableName
// 	}
// 	if n.TypeName != nil {
// 		str += ":" + *n.TypeName
// 	}
// 	props := make(map[string]string)
// 	for k, v := range n.Props {
// 		props[k] = v
// 	}
// 	props["tenant"] = tenant

// 	str += "{"
// 	firstProp := true
// 	for k, v := range props {
// 		if !firstProp {
// 			str += ","
// 		}
// 		str += fmt.Sprintf("%s:'%s'", k, v)
// 		firstProp = false
// 	}
// 	str += "}"

// 	return str
// }

// func (q *CypherQuery) ToString(tenant string) string {
// 	str := "MATCH "
// 	str += fmt.Sprintf("(%s)", q.MatchNode.ToString(tenant))

// 	if q.Relationship != nil {
// 		if q.Relationship.Direction == REL_FROM {
// 			str += "<-"
// 		} else {
// 			str += "-"
// 		}
// 		if q.Relationship.Props != nil {
// 			str += fmt.Sprintf("[%s]", q.Relationship.Props.ToString(tenant))
// 		}
// 		if q.Relationship.Direction == REL_TO {
// 			str += "->"
// 		} else {
// 			str += "-"
// 		}
// 		str += fmt.Sprintf("(%s)", q.Relationship.Target.ToString(tenant))
// 	}

// 	if q.Return != nil {
// 		str += " RETURN "
// 		firstRet := true
// 		for _, ret := range q.Return {
// 			if !firstRet {
// 				str += ","
// 			}
// 			str += ret.ToString()
// 			firstRet = true
// 		}
// 	}
// 	return str
// }

func (n *CypherNode) ToString() string {
	str := ""
	if n.VariableName != nil {
		str = *n.VariableName
	}
	if n.TypeName != nil {
		str += ":" + *n.TypeName
	}
	props := make(map[string]string)
	for k, v := range n.Props {
		props[k] = v
	}

	str += "{"
	firstProp := true
	for k, v := range props {
		if !firstProp {
			str += ","
		}
		str += fmt.Sprintf("%s:'%s'", k, v)
		firstProp = false
	}
	str += "}"

	return str
}

func (q *CypherQuery) ToString() string {
	str := "MATCH "
	str += fmt.Sprintf("(%s)", q.MatchNode.ToString())

	if q.Relationship != nil {
		if q.Relationship.Direction == REL_FROM {
			str += "<-"
		} else {
			str += "-"
		}
		if q.Relationship.Props != nil {
			str += fmt.Sprintf("[%s]", q.Relationship.Props.ToString())
		}
		if q.Relationship.Direction == REL_TO {
			str += "->"
		} else {
			str += "-"
		}
		str += fmt.Sprintf("(%s)", q.Relationship.Target.ToString())
	}

	if q.Return != nil {
		str += " RETURN "
		firstRet := true
		for _, ret := range q.Return {
			if !firstRet {
				str += ","
			}
			str += ret.ToString()
			firstRet = true
		}
	}
	return str
}
