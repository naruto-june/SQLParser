package selecter

import (
	"errors"
)

/*
SELECT
    [ DISTINCT ]
    col_name [, col_name] ...
    [FROM table_references
    [WHERE where_condition]
    [GROUP BY {col_name | expr | position}, ... [WITH ROLLUP]]
    [ORDER BY {col_name | expr | position}
      [ASC | DESC], ... [WITH ROLLUP]]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
*/

/*

[
    {SELECT element []}
    {* element []}
    {FROM element []}
    {t_lawsuits_relations_3 element []}
    {where element []}
    {(eid='00aa7cf1-af83-4534-aad4-c0598ce43338' AND u_tags IN ( '0','') ) AND ((doc_type='判决书') OR (doc_type='裁定书' AND trail_result IN ( '01','02','03','99','') )) block
        [
            {eid='00aa7cf1-af83-4534-aad4-c0598ce43338' AND u_tags IN ( '0','')  block
                [
                    {eid='00aa7cf1-af83-4534-aad4-c0598ce43338' value []}
                    {AND element []}
                    {u_tags element []}
                    {IN element []}
                    { '0','' block [{'0' value []} {'' value []}]}
                ]
            }
            {AND element []}
            {(doc_type='判决书') OR (doc_type='裁定书' AND trail_result IN ( '01','02','03','99','') ) block
                [
                    {doc_type='判决书' block
                        [
                            {doc_type='判决书' value []}
                        ]
                    }
                    {OR element []}
                    {doc_type='裁定书' AND trail_result IN ( '01','02','03','99','')  block
                        [
                            {doc_type='裁定书' value []}
                            {AND element []}
                            {trail_result element []}
                            {IN element []}
                            { '01','02','03','99','' block
                                [
                                    {'01' value []}
                                    {'02' value []}
                                    {'03' value []}
                                    {'99' value []}
                                    {'' value []}
                                ]
                            }
                        ]
                    }
                ]
            }
        ]
    }
    {ORDER element []}
    {BY element []}
    {pub_date element []}
    {DESC element []}
    {LIMIT element []}
    {5 element []}
]
*/

// Token 分词
type Token struct {
	Raw  string
	Type string // element/block/value
	TS   Tokens
}

// Tokens 分词集合
type Tokens []Token

// Split2Tokens 拆解sql
func Split2Tokens(sql string) (Tokens, error) {
	//fmt.Println(sql)

	byteIdx := -1
	runeIdx := -1
	word := make([]rune, 0)
	ret := make(Tokens, 0)
	for idx, char := range sql {
		if idx <= byteIdx {
			continue
		}
		byteIdx = idx
		if 0 > runeIdx {
			runeIdx = 0
		} else {
			runeIdx = runeIdx + 1
		}

		if '(' == char { //取出完整块
			if len(word) > 0 {
				t := Token{
					Raw:  string(word),
					Type: "element",
				}

				word = word[:0]
				ret = append(ret, t)
			}

			stack := make([]rune, 0) //长度参数务必传0
			stack = append(stack, char)
			tmpWord := make([]rune, 0)

			tmp := []rune(sql)
			length := len(tmp)
			for {
				runeIdx = runeIdx + 1
				if runeIdx >= length {
					return nil, errors.New("check to end but no matched bracket")
				}

				c := tmp[runeIdx]
				if ')' == c {
					stack = stack[:len(stack)-1]
					tmpWord = append(tmpWord, c)
					if len(stack) <= 0 {
						/*if len(word) <= 0 {
							return nil, errors.New("empty content in a bracket pair")
						}*/

						var t Token
						if len(word) <= 0 {
							t = Token{
								Raw:  "",
								Type: "block",
								TS: Tokens{
									Token{
										Raw:  "",
										Type: "element",
									},
								},
							}
						} else {
							ts, err := Split2Tokens(string(word))
							if nil != err {
								return nil, err
							}

							t = Token{
								Raw:  string(word),
								Type: "block",
								TS:   ts,
							}
						}

						byteIdx = idx + len(string(tmpWord))
						word = word[:0]
						ret = append(ret, t)
						break
					}

					word = append(word, c)
					continue
				}

				if '(' == c {
					stack = append(stack, c)
				}

				word = append(word, c)
				tmpWord = append(tmpWord, c)
			}
			continue
		}

		if '\'' == char { //取出完整字符串
			stack := make([]rune, 0)
			tmpWord := make([]rune, 0)
			stack = append(stack, char)
			word = append(word, char)

			tmp := []rune(sql)
			length := len(tmp)
			for {
				runeIdx = runeIdx + 1
				if runeIdx >= length {
					return nil, errors.New("check to end but no matched single quote:" + sql)
				}

				c := tmp[runeIdx]
				if '\'' == c && '\\' != tmp[runeIdx-1] { // 字符串结束
					stack = stack[:len(stack)-1]
					word = append(word, c)
					tmpWord = append(tmpWord, c)
					if len(stack) <= 0 {
						t := Token{
							Raw:  string(word),
							Type: "value",
						}

						byteIdx = idx + len(string(tmpWord))
						word = word[:0]
						ret = append(ret, t)
						break
					}

					return nil, errors.New("no matched single quote")
				}

				word = append(word, c)
				tmpWord = append(tmpWord, c)
			}
			continue
		}

		if '"' == char { //取出完整字符串
			stack := make([]rune, 0)
			stack = append(stack, char)
			tmpWord := make([]rune, 0)
			word = append(word, char)

			tmp := []rune(sql)
			length := len(tmp)
			for {
				runeIdx = runeIdx + 1
				if runeIdx >= length {
					return nil, errors.New("check to end but no matched quote")
				}

				c := tmp[runeIdx]
				if '"' == c && '\\' != tmp[runeIdx-1] { // 字符串结束
					stack = stack[:len(stack)-1]
					word = append(word, c)
					tmpWord = append(tmpWord, c)
					if len(stack) <= 0 {
						t := Token{
							Raw:  string(word),
							Type: "value",
						}

						byteIdx = idx + len(string(tmpWord))
						word = word[:0]
						ret = append(ret, t)
						break
					}

					return nil, errors.New("no matched quote")
				}

				word = append(word, c)
				tmpWord = append(tmpWord, c)
			}
			continue
		}

		if char == ' ' || char == ',' {
			if len(word) > 0 {
				t := Token{
					Raw:  string(word),
					Type: "element",
				}

				word = word[:0]
				ret = append(ret, t)
				continue
			} else {
				// fmt.Println("####")
			}

		} else {
			word = append(word, char)
		}
	}

	if len(word) > 0 {
		t := Token{
			Raw:  string(word),
			Type: "element",
		}

		word = word[:0]
		ret = append(ret, t)
	}
	//fmt.Println(ret)
	return ret, nil
}
