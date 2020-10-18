package selecter

import (
	"errors"
	"strconv"
	"strings"
)

type groupBy struct {
	Descr string `json:"descr,omitempty"`
}

type groupByT map[string]groupBy

/*
	group by a,b,c
	从by关键字开始，直到遇到order/limit或结束
*/
func parseGroupBy(tObj Tokens) ([]groupByT, int, error) {
	lastIdx := -1
	length := len(tObj)
	elem := make(groupByT, 0)
	ret := make([]groupByT, 0)
	for k, v := range tObj {
		if k <= lastIdx {
			continue
		}
		lastIdx = k

		token := strings.ToLower(v.Raw)
		token = strings.Trim(token, " ")
		if "having" == token {
			return nil, 0, errors.New("having not supported")
		}
		if "order" == token || "limit" == token {
			break
		}
		if "by" == token {
			continue
		}

		if "block" == v.Type { // group by (a,b),c   group by (left(a,3),c)
			subGB, _, err := parseGroupBy(v.TS)
			if nil != err {
				return nil, 0, err
			}
			if len(subGB) <= 0 {
				return nil, 0, errors.New("illegal group")
			}
			ret = append(ret, subGB...)

			/*subLastIdx := 0
			subLength := len(v.TS)
			for ik, iv := range v.TS {
				if ik <= subLastIdx {
					continue
				}
				subLastIdx = ik

				if "block" == iv.Type {
					return nil, 0, errors.New("not support more than 1 embed")
				}

				if ik+1 < subLength && "block" == v.TS[ik+1].Type {
					subToken := strings.ToLower(v.TS[ik+1].TS[0].Raw)
					subToken = strings.Trim(subToken, " ")
					if _, err := strconv.ParseFloat(subToken, 64); nil == err {
						return nil, 0, errors.New("not support such groupby model in block")
					}

					elem[subToken] = groupBy{Descr: iv.Raw + "(" + v.TS[ik+1].Raw + ")"}
					subLastIdx = subLastIdx + 1
				} else {
					elem[iv.Raw] = groupBy{Descr: iv.Raw}
				}

				ret = append(ret, elem)
				elem = make(groupByT, 0)
			}*/
		} else {
			if k+1 < length && "block" == tObj[k+1].Type {
				_, ok := supportSQLFunc[token]
				if ok {
					subToken := strings.ToLower(tObj[k+1].TS[0].Raw)
					subToken = strings.Trim(subToken, " ")
					if _, err := strconv.ParseFloat(subToken, 64); nil == err {
						return nil, 0, errors.New("not support such groupby model")
					}

					elem[subToken] = groupBy{Descr: v.Raw + "(" + tObj[k+1].Raw + ")"}
					lastIdx = lastIdx + 1
				} else {
					elem[v.Raw] = groupBy{Descr: v.Raw}
					//return nil, 0, errors.New("funcion expect or function not support")
				}

			} else {
				elem[v.Raw] = groupBy{Descr: v.Raw}
			}

			ret = append(ret, elem)
			elem = make(groupByT, 0)
		}
	}

	return ret, lastIdx, nil
}
