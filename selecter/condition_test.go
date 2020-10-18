package selecter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestOnCondition 测试 where
func TestOnCondition(t *testing.T) {
	/*
		0 where field_a=1
		1 where field_a="1"
		2 where field_a='1'
		3 where field_a=1 and field_b="2"
		4 where field_a>=1 or field_b="2"
		5 where field_a is not null or field_b="2"
		6 where field_a is null or field_b="2"
		7 where field_a in (1,3,45) or field_b="2"
		8 where field_a in (select id from table_tmp where name="test") or field_b="2"
		9 where field_a<>3 or field_b not in (select name from table_tmp where id>10)
		10 where field_a != 3 or field_b not in (select name from table_tmp where id>10)
		11 where field_a in (1,3,45) or field_b not like "11%2%"
		12 where (eid="534472fd-7d53-4958-8132-d6a6242423d8" and ((date>="2017-06-03" and date<="2017-06-03") or (name='zhangsan' and name='lisi')) or id=10) or (id>=10 and id<=100)
		13 where (eid in (select eid from t_enterprise_0 where score>100) and ((date>="2017-06-03" and date<="2017-06-03") or (name='zhangsan' and name='lisi')) or id=10) or (id>=10 and id<=100)
	*/
	ts := []Tokens{
		Tokens{
			Token{
				Raw:  "field_a=1",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "=",
				Type: "element",
			},
			Token{
				Raw:  "1",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "=1",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a=",
				Type: "element",
			},
			Token{
				Raw:  "1",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "=",
				Type: "element",
			},
			Token{
				Raw:  "\"2\"",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "=",
				Type: "element",
			},
			Token{
				Raw:  "'1'",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a=1",
				Type: "element",
			},
			Token{
				Raw:  "and",
				Type: "element",
			},
			Token{
				Raw:  "field_b=\"2\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a>=1",
				Type: "element",
			},
			Token{
				Raw:  "and",
				Type: "element",
			},
			Token{
				Raw:  "field_b=\"2\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "is",
				Type: "element",
			},
			Token{
				Raw:  "not",
				Type: "element",
			},
			Token{
				Raw:  "null",
				Type: "element",
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b=\"2\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "is",
				Type: "element",
			},
			Token{
				Raw:  "null",
				Type: "element",
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b=\"2\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "in",
				Type: "element",
			},
			Token{
				Raw:  "1,3,45",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "1",
						Type: "element",
					},
					Token{
						Raw:  "3",
						Type: "element",
					},
					Token{
						Raw:  "45",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b=\"2\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "in",
				Type: "element",
			},
			Token{
				Raw:  "select id from table_tmp where name=\"test\"",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select",
						Type: "element",
					},
					Token{
						Raw:  "id",
						Type: "element",
					},
					Token{
						Raw:  "from",
						Type: "element",
					},
					Token{
						Raw:  "table_tmp",
						Type: "element",
					},
					Token{
						Raw:  "where",
						Type: "element",
					},
					Token{
						Raw:  "name=\"test\"",
						Type: "value",
					},
				},
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b=\"2\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "field_a<>3",
				Type: "element",
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b",
				Type: "element",
			},
			Token{
				Raw:  "in",
				Type: "element",
			},
			Token{
				Raw:  "select name from table_tmp where id > 10",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select",
						Type: "element",
					},
					Token{
						Raw:  "name",
						Type: "element",
					},
					Token{
						Raw:  "from",
						Type: "element",
					},
					Token{
						Raw:  "table_tmp",
						Type: "element",
					},
					Token{
						Raw:  "where",
						Type: "element",
					},
					Token{
						Raw:  "id",
						Type: "element",
					},
					Token{
						Raw:  ">",
						Type: "element",
					},
					Token{
						Raw:  "10",
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
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "!=",
				Type: "element",
			},
			Token{
				Raw:  "3",
				Type: "element",
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b",
				Type: "element",
			},
			Token{
				Raw:  "not",
				Type: "element",
			},
			Token{
				Raw:  "in",
				Type: "element",
			},
			Token{
				Raw:  "select name from table_tmp where id > 10",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "select",
						Type: "element",
					},
					Token{
						Raw:  "name",
						Type: "element",
					},
					Token{
						Raw:  "from",
						Type: "element",
					},
					Token{
						Raw:  "table_tmp",
						Type: "element",
					},
					Token{
						Raw:  "where",
						Type: "element",
					},
					Token{
						Raw:  "id",
						Type: "element",
					},
					Token{
						Raw:  ">",
						Type: "element",
					},
					Token{
						Raw:  "10",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},

		// 11 where field_a in (1,3,45) or field_b not like "11%2%"
		Tokens{
			Token{
				Raw:  "field_a",
				Type: "element",
			},
			Token{
				Raw:  "in",
				Type: "element",
			},
			Token{
				Raw:  "1,3,45",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "1",
						Type: "element",
					},
					Token{
						Raw:  "2",
						Type: "element",
					},
					Token{
						Raw:  "45",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "field_b",
				Type: "element",
			},
			Token{
				Raw:  "not",
				Type: "element",
			},
			Token{
				Raw:  "like",
				Type: "element",
			},
			Token{
				Raw:  "\"11%2%\"",
				Type: "value",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		// 12 where (eid="534472fd-7d53-4958-8132-d6a6242423d8" and ((date>="2017-06-03" and date<="2017-06-03") or (name='zhangsan' and name='lisi')) or id=10) or (id>=10 and id<=100)
		Tokens{
			Token{
				Raw:  "eid=\"534472fd-7d53-4958-8132-d6a6242423d8\" and ((date>=\"2017-06-03\" and date<=\"2017-06-03\") or (name='zhangsan' and name='lisi')) or id=10",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "eid=\"534472fd-7d53-4958-8132-d6a6242423d8\"",
						Type: "value",
					},
					Token{
						Raw:  "and",
						Type: "element",
					},
					Token{
						Raw:  "(date>=\"2017-06-03\" and date<=\"2017-06-03\") or (name='zhangsan' and name='lisi')",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "date>=\"2017-06-03\" and date<=\"2017-06-03\"",
								Type: "block",
								TS: Tokens{
									Token{
										Raw:  "date>=\"2017-06-03\"",
										Type: "value",
									},
									Token{
										Raw:  "and",
										Type: "element",
									},
									Token{
										Raw:  "date<=\"2017-06-03\"",
										Type: "value",
									},
								},
							},
							Token{
								Raw:  "or",
								Type: "element",
							},
							Token{
								Raw:  "name='zhangsan' and name='lisi'",
								Type: "block",
								TS: Tokens{
									Token{
										Raw:  "name='zhangsan'",
										Type: "value",
									},
									Token{
										Raw:  "and",
										Type: "element",
									},
									Token{
										Raw:  "name='lisi'",
										Type: "value",
									},
								},
							},
						},
					},
					Token{
						Raw:  "or",
						Type: "element",
					},
					Token{
						Raw:  "id=10",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "id>=10 and id<=100",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "id>=10",
						Type: "element",
					},
					Token{
						Raw:  "and",
						Type: "element",
					},
					Token{
						Raw:  "id<=100",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
		// 13 where (eid not in (select eid from t_enterprise_0 where score>100) and ((date>="2017-06-03" and date<="2017-06-03") or (name='zhangsan' and name='lisi')) or id=10) or (id>=10 and id<=100)
		Tokens{
			Token{
				Raw:  "eid not in (select eid from t_enterprise_0 where score>100) and ((date>=\"2017-06-03\" and date<=\"2017-06-03\") or (name='zhangsan' and name='lisi')) or id=10",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "eid",
						Type: "element",
					},
					Token{
						Raw:  "not",
						Type: "element",
					},
					Token{
						Raw:  "in",
						Type: "element",
					},
					Token{
						Raw:  "select eid from t_enterprise_0 where score>100",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "select",
								Type: "element",
							},
							Token{
								Raw:  "eid",
								Type: "element",
							},
							Token{
								Raw:  "from",
								Type: "element",
							},
							Token{
								Raw:  "t_enterprise_0",
								Type: "element",
							},
							Token{
								Raw:  "where",
								Type: "element",
							},
							Token{
								Raw:  "score>100",
								Type: "element",
							},
						},
					},
					Token{
						Raw:  "and",
						Type: "element",
					},
					Token{
						Raw:  "(date>=\"2017-06-03\" and date<=\"2017-06-03\") or (name='zhangsan' and name='lisi')",
						Type: "block",
						TS: Tokens{
							Token{
								Raw:  "date>=\"2017-06-03\" and date<=\"2017-06-03\"",
								Type: "block",
								TS: Tokens{
									Token{
										Raw:  "date>=\"2017-06-03\"",
										Type: "value",
									},
									Token{
										Raw:  "and",
										Type: "element",
									},
									Token{
										Raw:  "date<=\"2017-06-03\"",
										Type: "value",
									},
								},
							},
							Token{
								Raw:  "or",
								Type: "element",
							},
							Token{
								Raw:  "name='zhangsan' and name='lisi'",
								Type: "block",
								TS: Tokens{
									Token{
										Raw:  "name='zhangsan'",
										Type: "value",
									},
									Token{
										Raw:  "and",
										Type: "element",
									},
									Token{
										Raw:  "name='lisi'",
										Type: "value",
									},
								},
							},
						},
					},
					Token{
						Raw:  "or",
						Type: "element",
					},
					Token{
						Raw:  "id=10",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "or",
				Type: "element",
			},
			Token{
				Raw:  "id>=10 and id<=100",
				Type: "block",
				TS: Tokens{
					Token{
						Raw:  "id>=10",
						Type: "element",
					},
					Token{
						Raw:  "and",
						Type: "element",
					},
					Token{
						Raw:  "id<=100",
						Type: "element",
					},
				},
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},

		// where left(a,3) = 'san' and b=4
		Tokens{
			Token{
				Raw:  "left",
				Type: "element",
			},
			Token{
				Raw:  "a,3",
				Type: "block",
			},
			Token{
				Raw:  "=",
				Type: "element",
			},
			Token{
				Raw:  "'san'",
				Type: "value",
			},
			Token{
				Raw:  "and",
				Type: "element",
			},
			Token{
				Raw:  "b=4",
				Type: "element",
			},
			Token{
				Raw:  "limit",
				Type: "element",
			},
		},
	}

	fmt.Println("----------Condition------------------")
	for k, v := range ts {
		fs, offset, err := parseCondition(v)
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
