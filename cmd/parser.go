package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/naruto-june/SQLParser/selecter"
)

var (
	rawSQL string
)

func init() {
	flag.StringVar(&rawSQL, "s", "select a,b from table_name where a=1", "sql to be parsed")
}

func main() {
	flag.Parse()

	//rawSQL := "select a, left(b,3) as m from table_a,table_use,(select a,b from table_b where a> 10) as b where b.a>100 group by a,left(b,2) order by left(a,1),b desc limit 5,10"
	pSQL, err := selecter.ParseSQL2Obj(nil, rawSQL)
	if nil != err {
		fmt.Println(err)
	}

	formatSQLJson, err := json.Marshal(pSQL)
	if err != nil {
		fmt.Println(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, formatSQLJson, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	out.WriteTo(os.Stdout)
}
