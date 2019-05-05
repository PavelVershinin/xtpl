package xtpl

import (
	"bytes"
)

func (x *xtpl) execFor(src []rune) func(vars *xVarCollection) []byte {
	var bracketEnd = getOffset(src, ")", "", true, true)
	var expr = src[1:bracketEnd]
	var content = x.buildTree(src[bracketEnd+1:], true)

	var initial = func(vars *xVarCollection) *xVar {
		return nil
	}
	var equal = func(vars *xVarCollection) *xVar {
		return xVarInit("", true)
	}
	var iteration = func(vars *xVarCollection) *xVar {
		return nil
	}
	var step = 0
	for i := 0; i < len(expr) && step <= 2; {
		commaPoint := getOffset(expr[i:], ";", "", true, true)
		if commaPoint < 0 {
			commaPoint = len(expr) - i
		}
		switch step {
		case 0:
			initial = x.exec(trim(expr[i : i+commaPoint]))
		case 1:
			equal = x.exec(trim(expr[i : i+commaPoint]))
		case 2:
			iteration = x.exec(trim(expr[i : i+commaPoint]))
		}
		i = i + commaPoint + 1
		step++
	}

	return func(vars *xVarCollection) []byte {
		var buff = &bytes.Buffer{}
		var pos = uint(0)
		for initial(vars); equal(vars).toBool(); iteration(vars) {
			for _, f := range content {
				buff.Write(f(vars))
			}
			pos++
			if pos >= cyclesLimit {
				break
			}
		}
		return buff.Bytes()
	}
}
