package main

import (
	"SQLParser/selecter"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/sql/util/parser", parseHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

// curl -i "localhost:8000/sql/util/parser" -d'select a,b from table_name where a>100'
func parseHandler(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "read body failed")
		return
	}

	//rawSQL := "select a, left(b,3) as m from table_a,table_use,(select a,b from table_b where a> 10) as b where b.a>100 group by a,left(b,2) order by left(a,1),b desc limit 5,10"
	rawSQL := string(bs)
	pSQL, err := selecter.ParseSQL2Obj(nil, rawSQL)
	if nil != err {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, err.Error())
		return
	}

	formatSQLJson, err := json.Marshal(pSQL)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, err.Error())
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, formatSQLJson, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	out.WriteTo(w)
}
