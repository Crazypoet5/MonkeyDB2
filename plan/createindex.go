package plan

import (
	"errors"

	"../exe"
	"../index"
	"../table"

	"../sql/syntax"
)

func CreateIndexPlan(stn *syntax.SyntaxTreeNode) (*exe.Relation, *Result, error) {
	res := NewResult()
	if stn.Name != "createindex" {
		return nil, nil, errors.New("Expected createtable but get " + stn.Name)
	}
	values := stn.Value.([][]byte)
	tab := table.GetTableByName(string(values[1]))
	if tab == nil {
		return nil, nil, errors.New("Table not exists")
	}
	var f *table.Field
	var kField int
	for k, v := range tab.Fields {
		if v.Name == string(values[2]) {
			f = &tab.Fields[k]
			kField = k
		}
	}
	if f == nil {
		return nil, nil, errors.New("field not exists")
	}
	indexName := string(values[0])
	f.Index = index.CreateIndex(index.UNIQUE, "db", string(values[1]), string(values[2]))
	reader := tab.FirstPage.NewReader()
	if reader.LoadIndex(f.Index, kField) != nil {
		f.Index.Delete()
		f.Index = nil
		return nil, nil, errors.New("Cannot create index on this column")
	}
	f.Index.Name = indexName
	//TODO: Load Index
	res.SetResult(0)
	return nil, res, nil
}
