package plan

import (
	"errors"
	"unsafe"

	"../exe"
	"../sql/syntax"
)

type WhereClause func(*exe.Relation) *exe.BitSet //Return sorted

func wherePlan(stn *syntax.SyntaxTreeNode) (WhereClause, error) {
	if stn.Name != "where" {
		return nil, errors.New("Expected where but get:" + stn.Name)
	}
	logical, err := logicalPlan(stn.Child[0])
	if err != nil {
		return nil, err
	}
	return logical, nil
}

func relationPlan(stn *syntax.SyntaxTreeNode) (WhereClause, error) {
	if stn.Name != "relations" {
		return nil, errors.New("Expected relations but get:" + stn.Name)
	}
	l := stn.Child[0]
	r := stn.Child[1]

	switch string(stn.Value.([]byte)) {
	case "<":
		switch l.Name {
		case "value":
			if l.ValueType == syntax.INT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(int)
						rn := r.Value.(int)
						if ln < rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(int)
						rn := r.Value.(float64)
						if float64(ln) < rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'<' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if float64(l.Value.(int)) < f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
			if l.ValueType == syntax.FLOAT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(float64)
						rn := float64(r.Value.(int))
						if ln < rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(float64)
						rn := r.Value.(float64)
						if float64(ln) < rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'<' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if l.Value.(float64) < f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
		case "string":
			switch r.Name {
			case "value":
				return nil, errors.New("'<' cannot be used between value and string")
			case "identical":
				return nil, errors.New("'<' cannot be used between string and identical")
			case "string":
				return nil, errors.New("'<' cannot be used betwrrn string and string")
			}
		case "identical":
			id := string(l.Value.([]byte))
			switch r.Name {
			case "value":
				if r.ValueType == syntax.INT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f < float64(r.Value.(int)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
				if r.ValueType == syntax.FLOAT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f < float64(r.Value.(float64)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			case "string":
				return nil, errors.New("'<' cannot be used between identical and string")
			case "identical":
				idr := string(r.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v, v2 *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}
						if v2 = rel.GetFieldByName(i, idr); v2 == nil {
							return ret
						}
						var f, f2 float64
						switch v.Kind {
						case exe.INT:
							f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
						case exe.FLOAT:
							f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
						case exe.STRING:
							return ret
						}
						switch v2.Kind {
						case exe.INT:
							f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
						case exe.FLOAT:
							f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
						case exe.STRING:
							return ret
						}
						if f < f2 {
							ret.Set(i)
						}
					}
					return ret
				}, nil
			}

		}
	case "<=":
		switch l.Name {
		case "value":
			if l.ValueType == syntax.INT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(int)
						rn := r.Value.(int)
						if ln <= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(int)
						rn := r.Value.(float64)
						if float64(ln) <= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'<=' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if float64(l.Value.(int)) <= f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
			if l.ValueType == syntax.FLOAT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(float64)
						rn := float64(r.Value.(int))
						if ln <= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(float64)
						rn := r.Value.(float64)
						if float64(ln) <= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'<=' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if l.Value.(float64) <= f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
		case "string":
			switch r.Name {
			case "value":
				return nil, errors.New("'<=' cannot be used between value and string")
			case "identical":
				return nil, errors.New("'<=' cannot be used between string and identical")
			case "string":
				return nil, errors.New("'<=' cannot be used betwrrn string and string")
			}
		case "identical":
			id := string(l.Value.([]byte))
			switch r.Name {
			case "value":
				if r.ValueType == syntax.INT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f <= float64(r.Value.(int)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
				if r.ValueType == syntax.FLOAT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f <= float64(r.Value.(float64)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			case "string":
				return nil, errors.New("'<=' cannot be used between identical and string")
			case "identical":
				idr := string(r.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v, v2 *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}
						if v2 = rel.GetFieldByName(i, idr); v2 == nil {
							return ret
						}
						var f, f2 float64
						switch v.Kind {
						case exe.INT:
							f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
						case exe.FLOAT:
							f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
						case exe.STRING:
							return ret
						}
						switch v2.Kind {
						case exe.INT:
							f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
						case exe.FLOAT:
							f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
						case exe.STRING:
							return ret
						}
						if f <= f2 {
							ret.Set(i)
						}
					}
					return ret
				}, nil
			}

		}
	case "=":
		switch l.Name {
		case "value":
			if l.ValueType == syntax.INT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(int)
						rn := r.Value.(int)
						if ln == rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(int)
						rn := r.Value.(float64)
						if float64(ln) == rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'==' cannot be used between value and string")
				case "identical":
					id := string(l.Value.([]byte))
					switch r.Name {
					case "value":
						if r.ValueType == syntax.INT {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									var v *exe.Value
									if v = rel.GetFieldByName(i, id); v == nil {
										return ret
									}
									var f float64
									switch v.Kind {
									case exe.INT:
										f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
									case exe.FLOAT:
										f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
									}
									if f == float64(r.Value.(int)) {
										ret.Set(i)
									}
								}
								return ret
							}, nil
						}
						if r.ValueType == syntax.FLOAT {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									var v *exe.Value
									if v = rel.GetFieldByName(i, id); v == nil {
										return ret
									}
									var f float64
									switch v.Kind {
									case exe.INT:
										f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
									case exe.FLOAT:
										f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
									}
									if f == float64(r.Value.(float64)) {
										ret.Set(i)
									}
								}
								return ret
							}, nil
						}
					case "string":
						return nil, errors.New("'==' cannot be used between value and string")
					case "identical":
						idr := string(r.Value.([]byte))
						return func(rel *exe.Relation) *exe.BitSet {
							ret := &exe.BitSet{}
							for i := 0; i < len(rel.Rows); i++ {
								var v, v2 *exe.Value
								if v = rel.GetFieldByName(i, id); v == nil {
									return ret
								}
								if v2 = rel.GetFieldByName(i, idr); v2 == nil {
									return ret
								}
								var f, f2 float64
								switch v.Kind {
								case exe.INT:
									f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
								case exe.FLOAT:
									f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
								case exe.STRING:
									return ret
								}
								switch v2.Kind {
								case exe.INT:
									f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
								case exe.FLOAT:
									f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
								case exe.STRING:
									return ret
								}
								if f == f2 {
									ret.Set(i)
								}
							}
							return ret
						}, nil
					}

				}
			}
			if l.ValueType == syntax.FLOAT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(float64)
						rn := float64(r.Value.(int))
						if ln == rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(float64)
						rn := r.Value.(float64)
						if float64(ln) == rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'<' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if l.Value.(float64) == f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
		case "string":
			switch r.Name {
			case "value":
				return nil, errors.New("'==' cannot be used between value and string")
			case "identical":
				id := string(r.Value.([]byte))
				str := string(l.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}

						if v.Kind == exe.STRING && string(v.Raw) == str {
							ret.Set(i)
						}
					}
					return ret
				}, nil
				//return nil, errors.New("'==' cannot be used between string and identical")
			case "string":
				return nil, errors.New("'==' cannot be used betwrrn string and string")
			}
		case "identical":
			id := string(l.Value.([]byte))
			switch r.Name {
			case "value":
				if r.ValueType == syntax.INT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f == float64(r.Value.(int)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
				if r.ValueType == syntax.FLOAT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f == float64(r.Value.(float64)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			case "string":

				str := string(r.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}

						if v.Kind == exe.STRING && string(v.Raw) == str {
							ret.Set(i)
						}
					}
					return ret
				}, nil
				//return nil, errors.New("'==' cannot be used between identical and string")
			case "identical":
				idr := string(r.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v, v2 *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}
						if v2 = rel.GetFieldByName(i, idr); v2 == nil {
							return ret
						}
						var f, f2 float64
						switch v.Kind {
						case exe.INT:
							f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
						case exe.FLOAT:
							f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
						case exe.STRING:
							return ret
						}
						switch v2.Kind {
						case exe.INT:
							f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
						case exe.FLOAT:
							f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
						case exe.STRING:
							return ret
						}
						if f == f2 {
							ret.Set(i)
						}
					}
					return ret
				}, nil
			}
		}
	case ">":
		switch l.Name {
		case "value":
			if l.ValueType == syntax.INT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(int)
						rn := r.Value.(int)
						if ln > rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(int)
						rn := r.Value.(float64)
						if float64(ln) > rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'>' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					switch r.Name {
					case "value":
						if r.ValueType == syntax.INT {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									var v *exe.Value
									if v = rel.GetFieldByName(i, id); v == nil {
										return ret
									}
									var f float64
									switch v.Kind {
									case exe.INT:
										f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
									case exe.FLOAT:
										f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
									}
									if f > float64(r.Value.(int)) {
										ret.Set(i)
									}
								}
								return ret
							}, nil
						}
						if r.ValueType == syntax.FLOAT {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									var v *exe.Value
									if v = rel.GetFieldByName(i, id); v == nil {
										return ret
									}
									var f float64
									switch v.Kind {
									case exe.INT:
										f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
									case exe.FLOAT:
										f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
									}
									if f > float64(r.Value.(float64)) {
										ret.Set(i)
									}
								}
								return ret
							}, nil
						}
					case "string":
						return nil, errors.New("'>' cannot be used between identical and string")
					case "identical":
						idr := string(r.Value.([]byte))
						return func(rel *exe.Relation) *exe.BitSet {
							ret := &exe.BitSet{}
							for i := 0; i < len(rel.Rows); i++ {
								var v, v2 *exe.Value
								if v = rel.GetFieldByName(i, id); v == nil {
									return ret
								}
								if v2 = rel.GetFieldByName(i, idr); v2 == nil {
									return ret
								}
								var f, f2 float64
								switch v.Kind {
								case exe.INT:
									f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
								case exe.FLOAT:
									f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
								case exe.STRING:
									return ret
								}
								switch v2.Kind {
								case exe.INT:
									f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
								case exe.FLOAT:
									f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
								case exe.STRING:
									return ret
								}
								if f > f2 {
									ret.Set(i)
								}
							}
							return ret
						}, nil
					}

				}
			}
			if l.ValueType == syntax.FLOAT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(float64)
						rn := float64(r.Value.(int))
						if ln > rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(float64)
						rn := r.Value.(float64)
						if float64(ln) > rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'>' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if l.Value.(float64) > f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
		case "string":
			switch r.Name {
			case "value":
				return nil, errors.New("'>' cannot be used between value and string")
			case "identical":
				return nil, errors.New("'>' cannot be used between string and identical")
			case "string":
				return nil, errors.New("'>' cannot be used betwrrn string and string")
			}
		case "identical":
			id := string(l.Value.([]byte))
			switch r.Name {
			case "value":
				if r.ValueType == syntax.INT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f > float64(r.Value.(int)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
				if r.ValueType == syntax.FLOAT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f > float64(r.Value.(float64)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			case "string":
				return nil, errors.New("'>' cannot be used between identical and string")
			case "identical":
				idr := string(r.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v, v2 *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}
						if v2 = rel.GetFieldByName(i, idr); v2 == nil {
							return ret
						}
						var f, f2 float64
						switch v.Kind {
						case exe.INT:
							f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
						case exe.FLOAT:
							f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
						case exe.STRING:
							return ret
						}
						switch v2.Kind {
						case exe.INT:
							f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
						case exe.FLOAT:
							f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
						case exe.STRING:
							return ret
						}
						if f > f2 {
							ret.Set(i)
						}
					}
					return ret
				}, nil
			}
		}
	case ">=":
		switch l.Name {
		case "value":
			if l.ValueType == syntax.INT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(int)
						rn := r.Value.(int)
						if ln >= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(int)
						rn := r.Value.(float64)
						if float64(ln) >= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'>=' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					switch r.Name {
					case "value":
						if r.ValueType == syntax.INT {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									var v *exe.Value
									if v = rel.GetFieldByName(i, id); v == nil {
										return ret
									}
									var f float64
									switch v.Kind {
									case exe.INT:
										f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
									case exe.FLOAT:
										f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
									}
									if f >= float64(r.Value.(int)) {
										ret.Set(i)
									}
								}
								return ret
							}, nil
						}
						if r.ValueType == syntax.FLOAT {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									var v *exe.Value
									if v = rel.GetFieldByName(i, id); v == nil {
										return ret
									}
									var f float64
									switch v.Kind {
									case exe.INT:
										f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
									case exe.FLOAT:
										f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
									}
									if f >= float64(r.Value.(float64)) {
										ret.Set(i)
									}
								}
								return ret
							}, nil
						}
					case "string":
						return nil, errors.New("'>=' cannot be used between identical and string")
					case "identical":
						idr := string(r.Value.([]byte))
						return func(rel *exe.Relation) *exe.BitSet {
							ret := &exe.BitSet{}
							for i := 0; i < len(rel.Rows); i++ {
								var v, v2 *exe.Value
								if v = rel.GetFieldByName(i, id); v == nil {
									return ret
								}
								if v2 = rel.GetFieldByName(i, idr); v2 == nil {
									return ret
								}
								var f, f2 float64
								switch v.Kind {
								case exe.INT:
									f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
								case exe.FLOAT:
									f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
								case exe.STRING:
									return ret
								}
								switch v2.Kind {
								case exe.INT:
									f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
								case exe.FLOAT:
									f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
								case exe.STRING:
									return ret
								}
								if f >= f2 {
									ret.Set(i)
								}
							}
							return ret
						}, nil
					}

				}
			}
			if l.ValueType == syntax.FLOAT {
				switch r.Name {
				case "value":
					if r.ValueType == syntax.INT {
						ln := l.Value.(float64)
						rn := float64(r.Value.(int))
						if ln >= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else if r.ValueType == syntax.FLOAT {
						ln := l.Value.(float64)
						rn := r.Value.(float64)
						if float64(ln) >= rn {
							return func(rel *exe.Relation) *exe.BitSet {
								ret := &exe.BitSet{}
								for i := 0; i < len(rel.Rows); i++ {
									ret.Set(i)
								}
								return ret
							}, nil
						} else {
							return func(rel *exe.Relation) *exe.BitSet {
								return &exe.BitSet{}
							}, nil
						}
					} else {
						return nil, errors.New("Unexpedted value type.")
					}
				case "string":
					return nil, errors.New("'>=' cannot be used between value and string")
				case "identical":
					id := string(r.Value.([]byte))
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if l.Value.(float64) >= f {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			}
		case "string":
			switch r.Name {
			case "value":
				return nil, errors.New("'>=' cannot be used between value and string")
			case "identical":
				return nil, errors.New("'>=' cannot be used between string and identical")
			case "string":
				return nil, errors.New("'>=' cannot be used betwrrn string and string")
			}
		case "identical":
			id := string(l.Value.([]byte))
			switch r.Name {
			case "value":
				if r.ValueType == syntax.INT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f >= float64(r.Value.(int)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
				if r.ValueType == syntax.FLOAT {
					return func(rel *exe.Relation) *exe.BitSet {
						ret := &exe.BitSet{}
						for i := 0; i < len(rel.Rows); i++ {
							var v *exe.Value
							if v = rel.GetFieldByName(i, id); v == nil {
								return ret
							}
							var f float64
							switch v.Kind {
							case exe.INT:
								f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
							case exe.FLOAT:
								f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
							}
							if f >= float64(r.Value.(float64)) {
								ret.Set(i)
							}
						}
						return ret
					}, nil
				}
			case "string":
				return nil, errors.New("'>=' cannot be used between identical and string")
			case "identical":
				idr := string(r.Value.([]byte))
				return func(rel *exe.Relation) *exe.BitSet {
					ret := &exe.BitSet{}
					for i := 0; i < len(rel.Rows); i++ {
						var v, v2 *exe.Value
						if v = rel.GetFieldByName(i, id); v == nil {
							return ret
						}
						if v2 = rel.GetFieldByName(i, idr); v2 == nil {
							return ret
						}
						var f, f2 float64
						switch v.Kind {
						case exe.INT:
							f = float64(*((*int)(unsafe.Pointer(&v.Raw[0]))))
						case exe.FLOAT:
							f = *((*float64)(unsafe.Pointer(&v.Raw[0])))
						case exe.STRING:
							return ret
						}
						switch v2.Kind {
						case exe.INT:
							f2 = float64(*((*int)(unsafe.Pointer(&v2.Raw[0]))))
						case exe.FLOAT:
							f2 = *((*float64)(unsafe.Pointer(&v2.Raw[0])))
						case exe.STRING:
							return ret
						}
						if f >= f2 {
							ret.Set(i)
						}
					}
					return ret
				}, nil
			}
		}

	}
	return nil, errors.New("Unexpected relations.")
}

func logicalPlan(stn *syntax.SyntaxTreeNode) (WhereClause, error) {
	switch stn.Name {
	case "relations":
		relation, err := relationPlan(stn)
		if err != nil {
			return nil, err
		}
		return relation, err
	case "logical":
		switch string(stn.Value.([]byte)) {
		case "not":
			logical, err := logicalPlan(stn.Child[0])
			if err != nil {
				return nil, err
			}
			return func(r *exe.Relation) *exe.BitSet {
				ac := logical(r)
				all := &exe.BitSet{}
				all.SetRange(0, len(r.Rows))
				all.Difference(ac)
				return all
			}, nil
		case "or":
			logical1, err := logicalPlan(stn.Child[0])
			if err != nil {
				return nil, err
			}
			logical2, err := logicalPlan(stn.Child[1])
			if err != nil {
				return nil, err
			}
			return func(r *exe.Relation) *exe.BitSet {
				l, right := logical1(r), logical2(r)
				l.Union(right)
				return l
			}, nil
		case "and":
			logical1, err := logicalPlan(stn.Child[0])
			if err != nil {
				return nil, err
			}
			logical2, err := logicalPlan(stn.Child[1])
			if err != nil {
				return nil, err
			}
			return func(r *exe.Relation) *exe.BitSet {
				l, right := logical1(r), logical2(r)
				l.Intersect(right)
				return l
			}, nil
		}
	}
	return nil, errors.New("Expect logical but get:" + stn.Name)
}
