package selecter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestOnLimitBy 测试 limit
func TestOnLimitBy(t *testing.T) {
	/*
		limit 10
		limit 0,10
	*/
	ts := []Tokens{
		Tokens{
			Token{
				Raw:  "100",
				Type: "element",
			},
		},
		Tokens{
			Token{
				Raw:  "5",
				Type: "element",
			},
			Token{
				Raw:  "10",
				Type: "element",
			},
		},
	}

	fmt.Println("----------limitby------------------")
	for k, v := range ts {
		fs, offset, err := parseLimit(v)
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
