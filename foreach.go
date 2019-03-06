package xtpl

import (
	"bytes"
	"strings"
)

func (x *xtpl) execForeach(src []rune) func(vars *xVarCollection) []byte {
	var sliceName, iteratorName, valueName string
	var hasIterator = false
	var bracketEnd = getOffset(src, ")", "", true, true)
	var expr = src[1:bracketEnd]
	var content = x.buildTree(src[bracketEnd+1:], true)
	var asPosition = getOffset(expr, " as ", "", true, true)

	if asPosition > 0 {
		sliceName = strings.TrimSpace(string(expr[:asPosition]))
	} else {
		return func(vars *xVarCollection) []byte {
			return []byte("Parse error")
		}
	}

	if toPosition := getOffset(expr[asPosition+4:], "=>", "", true, true); toPosition > 0 {
		iteratorName = strings.TrimSpace(string(expr[asPosition+4 : asPosition+4+toPosition]))
		valueName = strings.TrimSpace(string(expr[asPosition+4+toPosition+2 : bracketEnd-1]))
		hasIterator = true
	} else {
		valueName = strings.TrimSpace(string(expr[asPosition+4 : bracketEnd-1]))
	}

	var sliceFunc = x.exec([]rune(sliceName))

	return func(vars *xVarCollection) []byte {
		var buff = &bytes.Buffer{}
		var value = sliceFunc(vars)
		var pos = uint(0)
		switch value.vType {
		case varTypeMap:
			for key, val := range value.toMap() {
				if hasIterator {
					vars.setVar(iteratorName, key)
				}
				vars.setVar(valueName, val.toInterface())
				for _, f := range content {
					buff.Write(f(vars))
				}
				pos++
				if pos >= x.collection.cyclesLimit {
					break
				}
			}
		default:
			for i, val := range value.toSlice() {
				if hasIterator {
					vars.setVar(iteratorName, i)
				}
				vars.setVar(valueName, val.toInterface())
				for _, f := range content {
					buff.Write(f(vars))
				}
				pos++
				if pos >= x.collection.cyclesLimit {
					break
				}
			}
		}
		return buff.Bytes()
	}
}
