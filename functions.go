package xtpl

import "bytes"

var xtplFunctions =  map[string]interface{}{
	"len": func(slice []*xVar) int {
		return len(slice)
	},
	"join": func(a []*xVar, sep string) string {
		var buff bytes.Buffer
		for i, val := range a {
			buff.WriteString(val.toString())
			if i < len(a) - 1 {
				buff.WriteString(sep)
			}
		}
		return buff.String()
	},
	"in_array": func(arr []*xVar, needle string) bool {
		for _, v := range arr {
			if v.toString() == needle {
				return true
			}
		}
		return false
	},
}