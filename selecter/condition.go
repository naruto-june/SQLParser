package selecter

import (
	"errors"
	"strings"
)

func parseCondition(tObj Tokens) ([]map[string]interface{}, int, error) {
	lastIdx := -1
	length := len(tObj)
	ret := make([]map[string]interface{}, 0)
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
		if "order" == token || "limit" == token || "group" == token {
			break
		}

		for {
			lastIdx = lastIdx + 1
			if lastIdx < length {
				nextToken := strings.ToLower(tObj[lastIdx].Raw)
				nextToken = strings.Trim(nextToken, " ")
				if "or" == nextToken || "and" == nextToken {
					pCond, err := parseOneCond(tObj[k:lastIdx], "$"+nextToken)
					if nil != err {
						return nil, 0, err
					}

					ret = append(ret, pCond)
					break
				}

				if "having" == nextToken {
					return nil, 0, errors.New("having not supported")
				}
				if "order" == nextToken || "limit" == nextToken || "group" == nextToken {
					pCond, err := parseOneCond(tObj[k:lastIdx], "")
					if nil != err {
						return nil, 0, err
					}

					ret = append(ret, pCond)
					if lastIdx+1 < length && "by" == strings.ToLower(strings.Trim(tObj[lastIdx+1].Raw, " ")) || "limit" == nextToken {
						return ret, lastIdx, nil
					}
					break
				}
			} else {
				pCond, err := parseOneCond(tObj[k:lastIdx], "")
				if nil != err {
					return nil, 0, err
				}

				ret = append(ret, pCond)
				break
			}
		}
	}
	return ret, lastIdx, nil
}

func parseOneCond(tObj Tokens, logic string) (map[string]interface{}, error) {
	idx := -1
	lastIdx := -1
	length := len(tObj)
	ret := make(map[string]interface{}, 0)
	if "" != logic {
		ret["nextlinker"] = logic
	}
	ret["entity"] = make([]map[string]interface{}, 0)
	entElem := make(map[string]interface{}, 0)
	for k, v := range tObj {
		if k <= lastIdx {
			continue
		}
		lastIdx = k

		if "block" == v.Type { // block eid="534472fd-7d53-4958-8132-d6a6242423d8" and ((date>="2017-06-03" and date<="2017-06-03") or (name='zhangsan' and name='lisi')) or id=10
			conds, _, err := parseCondition(v.TS)
			if nil != err {
				return nil, err
			}

			ret["entity"] = append(ret["entity"].([]map[string]interface{}), conds...)
			idx = idx + len(conds)
		} else if "value" == v.Type { // value eid=\"534472fd-7d53-4958-8132-d6a6242423d8\"
			items := parseValue(v)
			if len(items) != 3 {
				return nil, errors.New("illegal value:" + v.Raw)
			}

			valueT := []byte(items[2])
			valueT = valueT[1 : len(valueT)-1]
			entElem[items[0]] = map[string]string{
				opMap[items[1]]: string(valueT),
			}

			idx = idx + 1
			ret["entity"] = append(ret["entity"].([]map[string]interface{}), entElem)
			entElem = make(map[string]interface{}, 0)
		} else { // element eid = \"534472fd-7d53-4958-8132-d6a6242423d8\" /  left(a,3) = 'san' / a =1 / a= 1
			token := strings.ToLower(v.Raw)
			token = strings.Trim(token, " ")
			if "or" == token || "and" == token {
				ret["entity"].([]map[string]interface{})[idx]["nextlinker"] = token
				continue
			}

			items := make([]string, 0)
			//如果是函数
			_, ok := supportSQLFunc[token]
			if ok {
				fs, _, err := parseFields(tObj[k : k+2])
				if nil != err {
					return nil, errors.New("function parameter expected")
				}
				if len(fs) <= 0 || len(fs) > 1 {
					return nil, errors.New("illegal condition field:" + v.Raw)
				}

				items = append(items, fs[0].Descr)
				lastIdx = lastIdx + 1
			} else {
				items = parseValue(v)
			}

			for {
				lastIdx = lastIdx + 1
				if lastIdx < length {
					nextToken := strings.ToLower(tObj[lastIdx].Raw)
					nextToken = strings.Trim(nextToken, " ")

					if "or" == nextToken || "and" == nextToken {
						goto over
					}

					if "in" == nextToken {
						if lastIdx+1 < length {
							if "block" == tObj[lastIdx+1].Type && "select" == strings.Trim(strings.ToLower(tObj[lastIdx+1].TS[0].Raw), " ") {
								pSQL, err := ParseSQL2Obj(tObj[lastIdx+1].TS, "")
								if nil != err {
									return nil, err
								}
								//format, _ := json.Marshal(pSQL)
								//fmt.Printf("field:%s\n sql:%s\n", items[0], format)

								idx = idx + 1
								entElem[items[0]] = map[string]SQLT{
									"$in": *pSQL,
								}
								ret["entity"] = append(ret["entity"].([]map[string]interface{}), entElem)
								entElem = make(map[string]interface{}, 0)
								lastIdx = lastIdx + 1

								if lastIdx+1 == length {
									return ret, nil
								}
								continue
							}
						} else {
							return nil, errors.New("in parameters expected")
						}
					} else if lastIdx+1 < length && "not" == nextToken && "in" == strings.Trim(strings.ToLower(tObj[lastIdx+1].Raw), " ") {
						if lastIdx+2 < length {
							if "block" == tObj[lastIdx+2].Type && "select" == strings.Trim(strings.ToLower(tObj[lastIdx+2].TS[0].Raw), " ") {
								pSQL, err := ParseSQL2Obj(tObj[lastIdx+2].TS, "")
								if nil != err {
									return nil, err
								}
								//format, _ := json.Marshal(pSQL)
								//fmt.Printf("field:%s\n sql:%s\n", items[0], format)

								idx = idx + 1
								entElem[items[0]] = map[string]SQLT{
									"$nin": *pSQL,
								}
								ret["entity"] = append(ret["entity"].([]map[string]interface{}), entElem)
								entElem = make(map[string]interface{}, 0)
								lastIdx = lastIdx + 2

								if lastIdx+1 == length {
									return ret, nil
								}
								continue
							}
						} else {
							return nil, errors.New("not in parameters expected")
						}
					}

					tmp := parseValue(tObj[lastIdx])
					items = append(items, tmp...)
				} else {
					if k+1 == lastIdx && 0 >= len(items) { //单元素
						items = parseValue(tObj[k])
					}

					goto over
				}
			}

		over:
			if len(items) != 3 && len(items) != 4 {
				return nil, errors.New("illegal value at end")
			}
			if len(items) == 3 { // like / is null
				op := strings.ToLower(items[1])
				if "in" == op {
					entElem[items[0]] = map[string]string{
						"$in": items[2],
					}
				} else if "like" == op {
					entElem[items[0]] = map[string]string{
						"$like": items[2],
					}
				} else if "is" == op {
					entElem[items[0]] = map[string]string{
						"$is": items[2],
					}
				} else {
					entElem[items[0]] = map[string]string{
						opMap[items[1]]: items[2],
					}
				}

				idx = idx + 1
				ret["entity"] = append(ret["entity"].([]map[string]interface{}), entElem)
				entElem = make(map[string]interface{}, 0)
			} else { // not like/ is not null
				op1 := strings.ToLower(items[1])
				op2 := strings.ToLower(items[2])
				if "not" == op1 {
					if "in" == op2 {
						entElem[items[0]] = map[string]string{
							"$nin": items[3],
						}
					} else if "like" == op2 {
						entElem[items[0]] = map[string]string{
							"$nlike": items[3],
						}
					} else {
						return nil, errors.New("like or in expected")
					}
				} else if "is" == op1 && "not" == op2 {
					entElem[items[0]] = map[string]string{
						"$nis": items[3],
					}
				} else {
					return nil, errors.New("not or is expected")
				}

				idx = idx + 1
				ret["entity"] = append(ret["entity"].([]map[string]interface{}), entElem)
				entElem = make(map[string]interface{}, 0)
			}
		}

	}

	return ret, nil
}

/*
	$in                 IN
	$nin                NOT IN
	$gt                 >
	$lt                 <
	$gte                >=
	$lte                <=
	$ne                 != / <>
	$like               like
	$nlike              not like
*/
var opArr = []string{
	">=",
	"<=",
	"!=",
	"<>",
	">",
	"<",
	"=",
}

var opMap = map[string]string{
	">=": "$gte",
	"<=": "$lte",
	"!=": "$ne",
	"<>": "$ne",
	">":  "$gt",
	"<":  "$lt",
	"=":  "$eq",
}

func idxOp(t string) (int, string) {
	for _, v := range opArr {
		idx := strings.Index(t, v)
		if idx < 0 {
			continue
		} else {
			return idx, v
		}
	}

	return -1, ""
}

func parseValue(item Token) []string {
	idx, op := idxOp(item.Raw)
	if idx < 0 {
		return []string{item.Raw}
	}

	ret := make([]string, 0)
	if idx == 0 {
		if len(item.Raw) == len(op) {
			ret = append(ret, op)
		} else {
			ret = append(ret, op)
			ret = append(ret, string([]byte(item.Raw)[len(op):]))
		}

	} else {
		if len(item.Raw)-1 == idx {
			ret = append(ret, string([]byte(item.Raw)[0:idx]))
			ret = append(ret, op)
		} else {
			ret = append(ret, string([]byte(item.Raw)[0:idx]))
			ret = append(ret, op)
			ret = append(ret, string([]byte(item.Raw)[idx+len(op):]))
		}
	}

	return ret
}
