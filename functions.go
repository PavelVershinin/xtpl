package xtpl

import (
	"bytes"
)

func (xc *XtplCollection) defaultFunctions() map[string]interface{} {
	return map[string]interface{}{
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
	}
}
