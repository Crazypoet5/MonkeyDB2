package plan

import (
	"errors"

	"../recovery"

	"../exe"
	"../sql/syntax"
)

func DirectPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	defer recovery.RestoreFrame()
	//	if stn.Name == "createtable" {
	//		return CreatePlan(stn)
	//	}
	//	if stn.Name == "createkv" {
	//		return CreateKVPlan(stn)
	//	}
	//	if stn.Name == "createindex" {
	//		return CreateIndexPlan(stn)
	//	}
	//	if stn.Name == "insert" {
	//		return insertPlan(stn)
	//	}
	//	if stn.Name == "dump" {
	//		return dumpPlan(stn)
	//	}
	//	if stn.Name == "select" {
	//		return selectPlan(stn)
	//	}
	//	if stn.Name == "delete" {
	//		return deletePlan(stn)
	//	}
	//	if stn.Name == "drop" {
	//		return dropPlan(stn)
	//	}
	//	if stn.Name == "dropindex" {
	//		return dropIndexPlan(stn)
	//	}
	//	if stn.Name == "update" {
	//		return updatePlan(stn)
	//	}
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
	}
	return nil, nil, errors.New("Unsopprted plan.")
}
