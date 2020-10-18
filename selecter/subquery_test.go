package selecter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestOnSubQuery 测试 subquery or table
func TestOnSubQuery(t *testing.T) {
	/*
		0 table_a
		1 table_a as a
		2 table_a as a,table_b as b
		3 table_a, (select xx from table_b) as b

		4 (table_a as a,table_b as b)
		5 (select xx from table_a) as a
		6 (table_a, (select xx from table_b) as b)
		7 ((select xx from table_a) as a, table_b)
		8 ((select xx from table_a) as a, (select xx from table_b) as b)
	*/
	ts := []Tokens{
		Tokens{
			Token{
				Raw:  "table_a",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "table_a",
				Type: "element",
			},
			Token{
				Raw:  "as",
				Type: "element",
			},
			Token{
				Raw:  "a",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "table_a",
				Type: "element",
			},
			Token{
				Raw:  "as",
				Type: "element",
			},
			Token{
				Raw:  "a",
				Type: "element",
			},
			Token{
				Raw:  "table_b",
				Type: "element",
			},
			Token{
				Raw:  "as",
				Type: "element",
			},
			Token{
				Raw:  "b",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "table_a",
				Type: "element",
			},
			Token{
				Raw:  "select xx from table_b",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select",
						Type: "element",
					},
					Token{
						Raw:  "xx",
						Type: "element",
					},
					Token{
						Raw:  "from",
						Type: "element",
					},
					Token{
						Raw:  "table_b",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "as",
				Type: "element",
			},
			Token{
				Raw:  "b",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "table_a as a,table_b as b",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "table_a",
						Type: "element",
					},
					Token{
						Raw:  "as",
						Type: "element",
					},
					Token{
						Raw:  "a",
						Type: "element",
					},
					Token{
						Raw:  "table_b",
						Type: "element",
					},
					Token{
						Raw:  "as",
						Type: "element",
					},
					Token{
						Raw:  "b",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "select xx from table_a",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select",
						Type: "element",
					},
					Token{
						Raw:  "xx",
						Type: "element",
					},
					Token{
						Raw:  "from",
						Type: "element",
					},
					Token{
						Raw:  "table_a",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "as",
				Type: "element",
			},
			Token{
				Raw:  "a",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "table_a, (select xx from table_b) as b",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "table_a",
						Type: "element",
					},
					Token{
						Raw:  "select xx from table_b",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "select",
								Type: "element",
							},
							Token{
								Raw:  "xx",
								Type: "element",
							},
							Token{
								Raw:  "from",
								Type: "element",
							},
							Token{
								Raw:  "table_b",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "as",
						Type: "element",
					},
					Token{
						Raw:  "b",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "(select xx from table_a) as a, table_b",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select xx from table_a",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "select",
								Type: "element",
							},
							Token{
								Raw:  "xx",
								Type: "element",
							},
							Token{
								Raw:  "from",
								Type: "element",
							},
							Token{
								Raw:  "table_a",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "as",
						Type: "element",
					},
					Token{
						Raw:  "a",
						Type: "element",
					},
					Token{
						Raw:  "table_b",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "(select xx from table_a) as a, (select xx from table_b) as b",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select xx from table_a",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "select",
								Type: "element",
							},
							Token{
								Raw:  "xx",
								Type: "element",
							},
							Token{
								Raw:  "from",
								Type: "element",
							},
							Token{
								Raw:  "table_a",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "as",
						Type: "element",
					},
					Token{
						Raw:  "a",
						Type: "element",
					},
					Token{
						Raw:  "select xx from table_b",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "select",
								Type: "element",
							},
							Token{
								Raw:  "xx",
								Type: "element",
							},
							Token{
								Raw:  "from",
								Type: "element",
							},
							Token{
								Raw:  "table_b",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "as",
						Type: "element",
					},
					Token{
						Raw:  "b",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
	}

	fmt.Println("----------subquery------------------")
	for k, v := range ts {
		sqls, offset, err := parseSubQuery(v)
		if nil != err {
			fmt.Println(err)
			continue
		}
		//fmt.Printf("k:%d \t offset:%d \t sqls:%#v \n", k, offset, sqls)
		fmt.Printf("k:%d \t offset:%d \n", k, offset)

		format, err := json.Marshal(&sqls)
		if err != nil {
			fmt.Println(err)
		}
		var out bytes.Buffer
		err = json.Indent(&out, format, "", "\t")
		if err != nil {
			fmt.Println(err)
		}
		out.WriteTo(os.Stdout)
	}
}
