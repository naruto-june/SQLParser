package selecter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestOnOrderBy 测试 order by
func TestOnOrderBy(t *testing.T) {
	/*
		order by a
		order by a,b
		order by left(a,2),b desc
		order by a asc ,left(b,2)
		order by a,(left(b,2) desc,c),d
		order by (a,b),c,(d,e)
		order by ((a desc,b),c)
	*/
	ts := []Tokens{
		Tokens{
			Token{
				Raw:  "by",
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
				Raw:  "by",
				Type: "element",
			},
			Token{
				Raw:  "a",
				Type: "element",
			},
			Token{
				Raw:  "b",
				Type: "element",
			},
			//Token{
			//	Raw:  "limit",
			//	Type: "element",
			//},
		},
		Tokens{
			Token{
				Raw:  "by",
				Type: "element",
			},
			Token{
				Raw:  "left",
				Type: "element",
			},
			Token{
				Raw:  "a,2",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "a",
						Type: "element",
					},
					Token{
						Raw:  "2",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "b",
				Type: "element",
			},
			Token{
				Raw:  "desc",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "by",
				Type: "element",
			},
			Token{
				Raw:  "a",
				Type: "element",
			},
			Token{
				Raw:  "asc",
				Type: "element",
			},
			Token{
				Raw:  "left",
				Type: "element",
			},
			Token{
				Raw:  "b,2",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "b",
						Type: "element",
					},
					Token{
						Raw:  "2",
						Type: "element",
					},
				},
			},
			//Token{
			//	Raw:  "limit",
			//	Type: "element",
			//},
		},
		Tokens{
			Token{
				Raw:  "by",
				Type: "element",
			},
			Token{
				Raw:  "a",
				Type: "element",
			},
			Token{
				Raw:  "left(b,2) desc,c",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "left",
						Type: "element",
					},
					Token{
						Raw:  "b,2",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "b",
								Type: "element",
							},
							Token{
								Raw:  "2",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "desc",
						Type: "element",
					},
					Token{
						Raw:  "c",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "d",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "a,b",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "a",
						Type: "element",
					},
					Token{
						Raw:  "b",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "c",
				Type: "element",
			},
			Token{
				Raw:  "d,e",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "d",
						Type: "element",
					},
					Token{
						Raw:  "e",
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
				Raw:  "(a,b),c",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "a desc,b",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "a",
								Type: "element",
							},
							Token{
								Raw:  "desc",
								Type: "element",
							},
							Token{
								Raw:  "b",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "c",
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

	fmt.Println("----------orderby------------------")
	for k, v := range ts {
		fs, offset, err := parseOrderBy(v)
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
