package selecter

import (
	"errors"
	"strconv"
	"strings"
)

type limitByT struct {
	Offset float64 `json:"offset,omitempty"`
	Number float64 `json:"number,omitempty"`
}

// 只支持 limit 10 或者limit 0,10两种格式
func parseLimit(tObj Tokens) (*limitByT, int, error) {
	length := len(tObj)
	ret := &limitByT{}
	if 1 == length {
		token := strings.ToLower(tObj[0].Raw)
		token = strings.Trim(token, " ")
		tmp, err := strconv.ParseFloat(token, 64)
		if nil != err {
			return nil, 0, errors.New("limit.number not accept non-number")
		}
		ret.Number = tmp
	} else if 2 == length {
		token := strings.ToLower(tObj[0].Raw)
		token = strings.Trim(token, " ")
		tmp, err := strconv.ParseFloat(token, 64)
		if nil != err {
			return nil, 0, errors.New("limit.offset not accept non-number")
		}
		ret.Offset = tmp

		token = strings.ToLower(tObj[1].Raw)
		token = strings.Trim(token, " ")
		tmp, err = strconv.ParseFloat(token, 64)
		if nil != err {
			return nil, 0, errors.New("limit.number not accept non-number")
		}
		ret.Number = tmp
	} else {
		return nil, 0, errors.New("not support such format limit")
	}

	return ret, length, nil
}
