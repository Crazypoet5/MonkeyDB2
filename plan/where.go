package plan

import (
	"../exe"
)

type WhereClause func(*exe.Relation) bool
