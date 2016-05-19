package plan

import (
	"errors"

	"../recovery"

	"../exe"
	"../sql/syntax"
)

func DirectPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	defer recovery.RestoreFrame()
	switch stn.Name {
	case "createtable":
		return CreatePlan(stn)
	case "createkv":
		return createKVPlan(stn)
	case "createindex":
		return CreateIndexPlan(stn)
	case "insert":
		return insertPlan(stn)
	case "dump":
		return dumpPlan(stn)
	case "select":
		return selectPlan(stn)
	case "delete":
		return deletePlan(stn)
	case "drop":
		return dropPlan(stn)
	case "dropindex":
		return dropIndexPlan(stn)
	case "update":
		return updatePlan(stn)
	case "setkv":
		return setKVPlan(stn)
	case "getkv":
		return getKVPlan(stn)
	case "removekv":
		return removeKVPlan(stn)
	case "showtable":
		return showTablePlan(stn)
	case "showtables":
		return showTablesPlan(stn)
	}
	return nil, nil, errors.New("Unsopprted plan.")
}
