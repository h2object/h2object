package page

import (
	"reflect"
	"github.com/h2object/cast"
)

type PageSort struct{
	Field string
	Ascend bool
	Pages []*Page
}

func (ps PageSort) Len() int {
	return len(ps.Pages)
}
// Less reports whether the element with
// index i should sort before the element with index j.
func (ps PageSort) Less(i, j int) bool {
	pi := ps.Pages[i]
	pj := ps.Pages[j]

	fi, ok := pi.meta[ps.Field]
	if !ok {
		return false
	}
	fj, ok := pj.meta[ps.Field]
	if !ok {
		return true
	}

	tpi := reflect.TypeOf(fi)
	tpj := reflect.TypeOf(fj)

	if tpi != tpj {
		if tpi.Kind() < tpj.Kind() {
			if ps.Ascend {
				return true
			} else {
				return false
			}
		} else {
			if ps.Ascend {
				return false
			} else {
				return true
			}
		}
	} else {
		switch tpi.Kind() {
		case reflect.Bool:
			vi := cast.ToBool(fi)
			vj := cast.ToBool(fj)
			if vi == false {
				if ps.Ascend {
					return true
				} else {
					return false
				}
			}
			if vj == false {
				if ps.Ascend {
					return false
				} else {
					return true
				}
			}
		case reflect.Int:
			fallthrough
        case reflect.Int8:
			fallthrough
        case reflect.Int16:
			fallthrough
        case reflect.Int32:
			fallthrough
        case reflect.Int64:
			fallthrough
        case reflect.Uint:
			fallthrough
        case reflect.Uint8:
			fallthrough
        case reflect.Uint16:
			fallthrough
        case reflect.Uint32:
			fallthrough
        case reflect.Uint64:
			vi := cast.ToInt(fi)
			vj := cast.ToInt(fj)
			if vi <= vj {
				if ps.Ascend {
					return true
				} else {
					return false
				}
			} else {
				if ps.Ascend {
					return false
				} else {
					return true
				}
			}
        case reflect.Float32:
			fallthrough
        case reflect.Float64:
			vi := cast.ToFloat64(fi)
			vj := cast.ToFloat64(fj)
			if vi <= vj {
				if ps.Ascend {
					return true
				} else {
					return false
				}
			} else {
				if ps.Ascend {
					return false
				} else {
					return true
				}
			}
		case reflect.String:
			vi := cast.ToString(fi)
			vj := cast.ToString(fj)
			if vi <= vj {
				if ps.Ascend {
					return true
				} else {
					return false
				}
			} else {
				if ps.Ascend {
					return false
				} else {
					return true
				}
			}
		default:
			if vi, err := cast.ToDurationE(fi); err == nil {
				vj := cast.ToDuration(fj)
				if vi <= vj {
					if ps.Ascend {
						return true
					} else {
						return false
					}
				} else {
					if ps.Ascend {
						return false
					} else {
						return true
					}
				}
			}

			if vi, err := cast.ToTimeE(fi); err == nil {
				vj := cast.ToTime(fj)
				if vi.UnixNano() <= vj.UnixNano() {
					if ps.Ascend {
						return true
					} else {
						return false
					}
				} else {
					if ps.Ascend {
						return false
					} else {
						return true
					}
				}
			}
		}
	}

	return false
}
// Swap swaps the elements with indexes i and j.
func (ps PageSort) Swap(i, j int) {
	ps.Pages[i], ps.Pages[j] = ps.Pages[j], ps.Pages[i]
}
