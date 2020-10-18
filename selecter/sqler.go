package selecter

import (
	"errors"
	"strings"
)

// SQLT 普通sql语句
type SQLT struct {
	Fields    []filedT                 `json:"fields,omitempty"`
	As        string                   `json:"as,omitempty"`
	SQLs      []SQLT                   `json:"sqls,omitempty"`
	Table     *table                   `json:"table,omitempty"`
	Condition []map[string]interface{} `json:"condition,omitempty"`
	OrderBy   []orderByT               `json:"orderby,omitempty"`
	GroupBy   []groupByT               `json:"groupby,omitempty"`
	LimitBy   *limitByT                `json:"limitby,omitempty"`
}

var supportSQLFunc = map[string]int{
	"ascii":             1,
	"char_length":       1,
	"character_length":  1,
	"concat":            1,
	"concat_ws":         1,
	"field":             1,
	"find_in_set":       1,
	"format":            1,
	"insert":            1,
	"locate":            1,
	"lcase":             1,
	"left":              1,
	"lower":             1,
	"lpad":              1,
	"ltrim":             1,
	"mid":               1,
	"position":          1,
	"repeat":            1,
	"replace":           1,
	"reverse":           1,
	"right":             1,
	"rpad":              1,
	"rtrim":             1,
	"space":             1,
	"strcmp":            1,
	"substr":            1,
	"substring":         1,
	"substring_index":   1,
	"trim":              1,
	"ucase":             1,
	"upper":             1,
	"abs":               1,
	"acos":              1,
	"asin":              1,
	"atan":              1,
	"atan2":             1,
	"avg":               1,
	"ceil":              1,
	"ceiling":           1,
	"cos":               1,
	"cot":               1,
	"count":             1,
	"degrees":           1,
	"exp":               1,
	"floor":             1,
	"greatest":          1,
	"least":             1,
	"ln":                1,
	"log":               1,
	"log10":             1,
	"log2":              1,
	"max":               1,
	"min":               1,
	"pow":               1,
	"power":             1,
	"radians":           1,
	"rand":              1,
	"round":             1,
	"sin":               1,
	"sign":              1,
	"sqrt":              1,
	"sum":               1,
	"tan":               1,
	"truncate":          1,
	"pi":                1,
	"adddate":           1,
	"addtime":           1,
	"curdate":           1,
	"current_date":      1,
	"current_time":      1,
	"current_timestamp": 1,
	"curtime":           1,
	"date":              1,
	"datediff":          1,
	"date_add":          1,
	"date_format":       1,
	"date_sub":          1,
	"day":               1,
	"dayname":           1,
	"dayofmonth":        1,
	"dayofweek":         1,
	"dayofyear":         1,
	"extract":           1,
	"from_days":         1,
	"hour":              1,
	"last_day":          1,
	"localtime":         1,
	"localtimestamp":    1,
	"makedate":          1,
	"maketime":          1,
	"microsecond":       1,
	"minute":            1,
	"monthname":         1,
	"month":             1,
	"now":               1,
	"period_add":        1,
	"period_diff":       1,
	"quarter":           1,
	"second":            1,
	"sec_to_time":       1,
	"str_to_date":       1,
	"subdate":           1,
	"subtime":           1,
	"sysdate":           1,
	"time":              1,
	"time_format":       1,
	"time_to_sec":       1,
	"timediff":          1,
	"timestamp":         1,
	"to_days":           1,
	"week":              1,
	"weekday":           1,
	"weekofyear":        1,
	"year":              1,
	"yearweek":          1,
	"bin":               1,
	"binary":            1,
	"cast":              1,
	"coalesce":          1,
	"connection_id":     1,
	"conv":              1,
	"charset":           1,
	"current_user":      1,
	"database":          1,
	"if":                1,
	"ifnull":            1,
	"isnull":            1,
	"last_insert_id":    1,
	"nullif":            1,
	"session_user":      1,
	"system_user":       1,
	"user":              1,
	"version":           1,
}

// ParseSQL2Obj 解析sql
func ParseSQL2Obj(ts Tokens, sql string) (*SQLT, error) {
	tObj := ts
	var err error
	if nil == tObj {
		tObj, err = Split2Tokens(sql)
		if nil != err {
			return nil, err
		}
	}

	ret := &SQLT{}
	lastIdx := -1
	for idx, v := range tObj {
		if idx <= lastIdx {
			continue
		}
		lastIdx = lastIdx + 1

		token := strings.ToLower(v.Raw)
		token = strings.Trim(token, " ")

		if "select" == token {
			fs, offsetIdx, err := parseFields(tObj[lastIdx+1:])
			if nil != err {
				return nil, err
			}
			lastIdx = lastIdx + offsetIdx
			ret.Fields = fs
			continue
		}

		if "from" == token {
			sqls, offsetIdx, err := parseSubQuery(tObj[lastIdx+1:])
			if nil != err {
				return nil, err
			}
			lastIdx = lastIdx + offsetIdx
			ret.SQLs = sqls

			continue
		}

		if "where" == token {
			conds, offsetIdx, err := parseCondition(tObj[lastIdx+1:])
			if nil != err {
				return nil, err
			}
			lastIdx = lastIdx + offsetIdx
			ret.Condition = conds

			continue
		}

		if "group" == token {
			nextToken := strings.ToLower(tObj[lastIdx+1].Raw)
			nextToken = strings.Trim(nextToken, " ")
			if "by" != nextToken {
				return nil, errors.New("by expected in groupby")
			}

			gby, offsetIdx, err := parseGroupBy(tObj[lastIdx+1:])
			if nil != err {
				return nil, err
			}
			lastIdx = lastIdx + offsetIdx
			ret.GroupBy = gby

			continue
		}

		if "order" == token {
			nextToken := strings.ToLower(tObj[lastIdx+1].Raw)
			nextToken = strings.Trim(nextToken, " ")
			if "by" != nextToken {
				return nil, errors.New("by expected in orderby")
			}

			oby, offsetIdx, err := parseOrderBy(tObj[lastIdx+1:])
			if nil != err {
				return nil, err
			}
			lastIdx = lastIdx + offsetIdx
			ret.OrderBy = oby

			continue
		}

		if "limit" == token {
			limitBy, offsetIdx, err := parseLimit(tObj[lastIdx+1:])
			if nil != err {
				return nil, err
			}
			lastIdx = lastIdx + offsetIdx
			ret.LimitBy = limitBy

			continue
		}
	}
	return ret, nil
}
