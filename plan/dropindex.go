package plan

import (
	"errors"

	"../exe"
	"../index"
	"../sql/syntax"
	"../table"
)

func dropIndexPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	re := NewResult()
	if stn.Name != "dropindex" {
		return nil, nil, errors.New("Expected drop but get " + stn.Name)
	}
	indexName := string(stn.Value.([]byte))
	for _, v := range index.IndexList {
		if v.Name == indexName {
			tab := table.GetTableByName(v.Table)
			if tab != nil {
				for k, vv := range tab.Fields {
					if vv.Name == v.Key {
						tab.Fields[k].Index = nil
					}
				}
			}
			v.Delete()
			re.SetResult(0)
			return nil, re, nil
		}
	}
	return nil, re, errors.New("Cannot find index")
}
