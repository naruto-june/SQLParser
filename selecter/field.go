package selecter

import (
	"errors"
	"strings"
)

type filedT struct {
	Distinct int    `json:"distinct,omitempty"`
	Descr    string `json:"descr,omitempty"`
	As       string `json:"as,omitempty"`
}

func parseFields(tObj Tokens) ([]filedT, int, error) {
	idx := -1
	lastIdx := -1
	length := len(tObj)
	elem := filedT{}
	ret := make([]filedT, 0)
	for k, v := range tObj {
		if k <= lastIdx {
			continue
		}
		lastIdx = k

		token := strings.ToLower(v.Raw)
		token = strings.Trim(token, " ")
		if "from" == token {
			break
		}
		if "select" == token {
			return nil, 0, errors.New("subquery not support in fields")
		}
		/*if "" == token {
			fmt.Println("find a space in field")
			continue
		}*/

		if "block" == v.Type {
			subFT, _, err := parseFields(v.TS)
			if nil != err {
				return nil, 0, err
			}
			if len(subFT) <= 0 {
				return nil, 0, errors.New("illegal field")
			}

			ret = append(ret, subFT...)
			continue
		}

		if "distinct" == token {
			elem.Distinct = 1
		} else if "as" == token {
			if k+1 >= length {
				return nil, 0, errors.New("idx beyond tObj")
			}
			ret[idx].As = tObj[k+1].Raw
			lastIdx = lastIdx + 1
		} else {
			if k+1 < length && "block" == tObj[k+1].Type {
				//TODO 判断下当前raw是不是mysql函数名称
				_, ok := supportSQLFunc[token]
				if ok {
					elem.Descr = token + "(" + tObj[k+1].Raw + ")"
					lastIdx = lastIdx + 1
				} else {
					elem.Descr = token
					//return nil, 0, errors.New("funcion expect or function not support")
				}

			} else {
				elem.Descr = token
			}

			idx = idx + 1
			ret = append(ret, elem)
			elem = filedT{}
		}
	}
	return ret, lastIdx, nil
}
