package selecter

import (
	"errors"
	"strings"
)

type table struct {
	TName string `json:"tname,omitempty"`
	As    string `json:"as,omitempty"`
}

/*
	case1: zipkin_spans
	case2: zipkin_spans as a
	case3: zipkin_spans as a,zipkin_spans as b
	case4: zipkin_spans, (select xx from table where ...) as a

	case5: (zipkin_spans as a,zipkin_spans as b)
	case6: (select xx from table where ...) as a
	case7: (zipkin_spans, (select xx from table where ...) as a)
	case8: ((select xx from table where ...) as a, zipkin_spans)
	case9: ((select xx from table where ...) as a, (select xx from table where ...) as b)
*/
// 后接where/order/group/limit
func parseSubQuery(tObj Tokens) ([]SQLT, int, error) {
	lastIdx := -1
	elemTable := table{}
	length := len(tObj)
	retSQLT := []SQLT{}
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
		if "where" == token || "order" == token || "limit" == token || "group" == token {
			break
		}

		if "block" == v.Type {
			subToken := strings.ToLower(v.TS[0].Raw)
			subToken = strings.Trim(subToken, " ")
			if "select" == subToken {
				sql, err := ParseSQL2Obj(v.TS, v.Raw)
				if nil != err {
					return nil, 0, err
				}

				if k+1 < length {
					nextToken := strings.ToLower(tObj[k+1].Raw)
					nextToken = strings.Trim(nextToken, " ")
					if "as" == nextToken {
						lastIdx = lastIdx + 1

						if k+2 < length && "element" == tObj[k+2].Type {

							sql.As = tObj[k+2].Raw
							lastIdx = lastIdx + 1
						} else {
							return nil, 0, errors.New("subquery as name expected")
						}
					}
				}

				retSQLT = append(retSQLT, *sql)
			} else {
				subSQLT, _, err := parseSubQuery(v.TS)
				if nil != err {
					return nil, 0, err
				}

				if len(subSQLT) > 0 {
					retSQLT = append(retSQLT, subSQLT...)
				}
			}
		} else { // element
			if k+1 < length {
				if "element" == tObj[k+1].Type { // as or next table-name
					subToken := strings.ToLower(tObj[k+1].Raw)
					subToken = strings.Trim(subToken, " ")
					if "as" == subToken {
						lastIdx = lastIdx + 1
						if k+2 < length && "element" == tObj[k+2].Type {
							elemTable.As = tObj[k+2].Raw
							lastIdx = lastIdx + 1
						} else {
							return nil, 0, errors.New("as name expected")
						}
					}

					elemTable.TName = v.Raw
					retSQLT = append(retSQLT, SQLT{Table: &table{TName: elemTable.TName, As: elemTable.As}})
					elemTable = table{}
				} else { // next is block
					elemTable.TName = v.Raw
					retSQLT = append(retSQLT, SQLT{Table: &table{TName: elemTable.TName, As: elemTable.As}})
					elemTable = table{}
				}
			} else { //结束
				elemTable.TName = v.Raw
				retSQLT = append(retSQLT, SQLT{Table: &table{TName: elemTable.TName, As: elemTable.As}})
				elemTable = table{}
			}
		}
	}
	return retSQLT, lastIdx, nil
}
