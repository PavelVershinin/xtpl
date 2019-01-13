package xtpl

import (
	"bytes"
)

type ifStruct struct {
	equal func(vars *xVarCollection) bool
	tree  []treeNode
}

func (x *xtpl) execIf(src []rune) func(vars *xVarCollection) []byte {
	var variants []ifStruct

	for len(src) > 0 {
		var variant ifStruct
		bracketBegin := getOffset(src, "(", "", true, true)
		bracketEnd := getOffset(src, ")", "", true, true)
		if bracketBegin > -1 && bracketEnd > bracketBegin {
			var equal = x.exec(src[bracketBegin : bracketEnd+1])
			variant.equal = func(vars *xVarCollection) bool {
				return equal(vars).toBool()
			}
		} else {
			variant.equal = func(vars *xVarCollection) bool {
				return true
			}
		}
		contentBegin := bracketEnd + 1
		if contentEnd := getOffset(src, "@else", "@if", false, false); contentEnd < 0 {
			variant.tree = x.buildTree(src[contentBegin:], true)
			src = nil
		} else {
			variant.tree = x.buildTree(src[contentBegin:contentEnd], true)
			src = src[contentEnd+5:]
		}
		variants = append(variants, variant)
	}
	return func(vars *xVarCollection) []byte {
		for _, variant := range variants {
			if variant.equal(vars) {
				var buff = &bytes.Buffer{}
				for _, f := range variant.tree {
					buff.Write(f(vars))
				}
				return buff.Bytes()
			}
		}
		return nil
	}
}
