package selecter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestOnField 测试字段集合
func TestOnField(t *testing.T) {
	/*
		field1
		field1,field2
		field1 as a,field2
		field1,field2 as b
		field1,left(field2, 5) as b
		left(field1, 5) as a, field2
		FIND_IN_SET("c", "a,b,c,d,e") as a, field2
		SIN(RADIANS(field1)) as a, field2
		now() as a, field2

		field1,(field2 as b,field3),field4
		(field1 as a, field2), field3
	*/
	ts := []Tokens{
		Tokens{
			Token{
				Raw:  "field1",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field1",
				Type: "element",
			},
			Token{
				Raw:  "field2",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field1",
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
				Raw:  "field2",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field1",
				Type: "element",
			},
			Token{
				Raw:  "field2",
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
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field1",
				Type: "element",
			},
			Token{
				Raw:  "left",
				Type: "element",
			},
			Token{
				Raw:  "field2, 5",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "field2",
						Type: "element",
					},
					Token{
						Raw:  "5",
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
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "left",
				Type: "element",
			},
			Token{
				Raw:  "field1, 5",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "field1",
						Type: "element",
					},
					Token{
						Raw:  "5",
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
				Raw:  "field2",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "FIND_IN_SET",
				Type: "element",
			},
			Token{
				Raw:  "\"c\", \"a,b,c,d,e\"",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "\"c\"",
						Type: "element",
					},
					Token{
						Raw:  "\"a,b,c,d,e\"",
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
				Raw:  "field2",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "SIN",
				Type: "element",
			},
			Token{
				Raw:  "RADIANS(field1,3)",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "RADIANS",
						Type: "element",
					},
					Token{
						Raw:  "field1,3",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "field1",
								Type: "element",
							},
							Token{
								Raw:  "3",
								Type: "element",
							},
						},
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
				Raw:  "field2",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "now",
				Type: "element",
			},
			Token{
				Raw:  "",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "",
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
				Raw:  "field2",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field1",
				Type: "element",
			},
			Token{
				Raw:  "field2 as b,field3",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "field2",
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
						Raw:  "field3",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "field4",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field1 as a, field2",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "field1",
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
						Raw:  "field2",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "field3",
				Type: "element",
			},
			Token{
				Raw:  "from",
				Type: "element",
			},
		},
	}

	fmt.Println("----------fields------------------")
	for k, v := range ts {
		fs, offset, err := parseFields(v)
		if nil != err {
			fmt.Println(err)
			continue
		}
		//fmt.Printf("k:%d \t offset:%d \t fs:%#v \n", k, offset, fs)
		fmt.Printf("k:%d \t offset:%d \n", k, offset)

		format, err := json.Marshal(&fs)
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
