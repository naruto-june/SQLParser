package selecter

import (
	"errors"
	"strconv"
	"strings"
)

type orderBy struct {
	Descr  string `json:"descr,omitempty"`
	Direct string `json:"direct,omitempty"`
}

type orderByT map[string]orderBy

func parseOrderBy(tObj Tokens) ([]orderByT, int, error) {
	lastIdx := -1
	length := len(tObj)
	elem := make(orderByT, 0)
	ret := make([]orderByT, 0)
	for k, v := range tObj {
		if k <= lastIdx {
			continue
		}
		lastIdx = k

		token := strings.ToLower(v.Raw)
		token = strings.Trim(token, " ")
		if "limit" == token {
			break
		}
		if "by" == token {
			continue
		}

		if "block" == v.Type { // order by (a,b),c desc  order by (left(a,3) desc,c)
			subOB, _, err := parseOrderBy(v.TS)
			if nil != err {
				return nil, 0, err
			}
			if len(subOB) <= 0 {
				return nil, 0, errors.New("illegal order")
			}
			ret = append(ret, subOB...)

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
						return nil, 0, errors.New("not support such orderby model in block")
					}

					tmp := orderBy{Descr: iv.Raw + "(" + v.TS[ik+1].Raw + ")"}
					if ik+2 < subLength {
						nextToken := strings.ToLower(v.TS[k+2].Raw)
						nextToken = strings.Trim(nextToken, " ")

						if "desc" == nextToken || "asc" == nextToken {
							tmp.Direct = nextToken
							subLastIdx = subLastIdx + 1
						}
					}

					elem[subToken] = tmp
					subLastIdx = subLastIdx + 1
				} else {
					tmp := orderBy{Descr: iv.Raw}
					if k+1 < subLength {
						nextToken := strings.ToLower(v.TS[ik+1].Raw)
						nextToken = strings.Trim(nextToken, " ")

						if "desc" == nextToken || "asc" == nextToken {
							tmp.Direct = nextToken
							subLastIdx = subLastIdx + 1
						}
					}

					elem[iv.Raw] = tmp
				}

				ret = append(ret, elem)
				elem = make(orderByT, 0)
			}*/
		} else {
			if k+1 < length && "block" == tObj[k+1].Type {
				_, ok := supportSQLFunc[token]
				if ok {
					subToken := strings.ToLower(tObj[k+1].TS[0].Raw)
					subToken = strings.Trim(subToken, " ")
					if _, err := strconv.ParseFloat(subToken, 64); nil == err {
						return nil, 0, errors.New("not support such orderby model")
					}

					tmp := orderBy{Descr: v.Raw + "(" + tObj[k+1].Raw + ")"}
					if k+2 < length {
						nextToken := strings.ToLower(tObj[k+2].Raw)
						nextToken = strings.Trim(nextToken, " ")

						if "desc" == nextToken || "asc" == nextToken {
							tmp.Direct = nextToken
							lastIdx = lastIdx + 1
						}
					}

					elem[subToken] = tmp
					lastIdx = lastIdx + 1
				} else {
					tmp := orderBy{Descr: v.Raw}
					if k+1 < length {
						nextToken := strings.ToLower(tObj[k+1].Raw)
						nextToken = strings.Trim(nextToken, " ")

						if "desc" == nextToken || "asc" == nextToken {
							tmp.Direct = nextToken
							lastIdx = lastIdx + 1
						}
					}

					elem[v.Raw] = tmp
					//return nil, 0, errors.New("funcion expect or function not support")
				}
			} else {
				tmp := orderBy{Descr: v.Raw}
				if k+1 < length {
					nextToken := strings.ToLower(tObj[k+1].Raw)
					nextToken = strings.Trim(nextToken, " ")

					if "desc" == nextToken || "asc" == nextToken {
						tmp.Direct = nextToken
						lastIdx = lastIdx + 1
					}
				}

				elem[v.Raw] = tmp
			}

			ret = append(ret, elem)
			elem = make(orderByT, 0)
		}
	}

	return ret, lastIdx, nil
}
