package xtpl

import (
	"regexp"
	"strings"
)

func (x *xtpl) exec(src []rune) func(vars *xVarCollection) *xVar {
	var functions []treeNode

	src = trim(src)

	// Все строки в кавычках переводим в переменные
	for i := 0; i < len(src); i++ {
		var openQuote = -1
		var closeQuote = -1
		if openQuote = getOffset(src[i:], `"`, "", true, true); openQuote > -1 {
			closeQuote = getOffset(src[i+openQuote+1:], `"`, "", true, true)
		} else if openQuote = getOffset(src[i:], `'`, "", true, true); openQuote > -1 {
			closeQuote = getOffset(src[i+openQuote+1:], `'`, "", true, true)
		}
		if openQuote > -1 && closeQuote > -1 {
			var value = string(src[i+openQuote+1 : i+openQuote+1+closeQuote])
			var varName = newVarName()
			functions = append(functions, func(vars *xVarCollection) []byte {
				vars.setVar(varName, value)
				return nil
			})
			src = append(src[:i+openQuote], append([]rune(varName), src[i+openQuote+1+closeQuote+1:]...)...)
		}
	}

	// Выгребаем функции
	for i := 0; i < len(src); i++ {
		if f, offset := x.userFunction(src[i:], false); f != nil {
			var varName = newVarName()
			functions = append(functions, func(vars *xVarCollection) []byte {
				vars.setVar(varName, f(vars).toInterface())
				return nil
			})
			src = append(src[:i], append([]rune(varName), src[i+offset:]...)...)
		}
	}

	// Выгребаем все вычисления в скобках
	for {
		closeBracketPosition := getOffset(src, ")", "", true, false)
		if closeBracketPosition == -1 {
			break
		}
		openBracketPosition := closeBracketPosition
		for ; openBracketPosition >= 0; openBracketPosition-- {
			if src[openBracketPosition] == '(' {
				break
			}
		}

		if openBracketPosition == closeBracketPosition {
			return func(vars *xVarCollection) *xVar {
				return xVarInit("", "parse error")
			}
		}
		var varName = newVarName()
		var f = x.exec(src[openBracketPosition+1 : closeBracketPosition])
		functions = append(functions, func(vars *xVarCollection) []byte {
			vars.setVar(varName, f(vars).toInterface())
			return nil
		})
		src = append(src[:openBracketPosition], append([]rune(varName), src[closeBracketPosition+1:]...)...)
	}

	// К этому моменту уже должно оставаться чистое выражение, без вызовов функций и скобок
	var expr = string(src)

	var reAssign = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?=[\s]?(.*)`)
	var rePlusPlus = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\+\+`)
	var reMinusMinus = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?--`)
	var reMultiple = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\*[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reDivision = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?/[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reDivisionRemainder = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?%[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reAddition = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\+[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reSubtraction = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?-[\s]?(\$[a-z0-9_]+|[0-9.]+)`)

	var reEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?==[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reNotEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?(!=|<>)[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reMore = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?>[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reLess = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?<[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reMoreOrEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?>=[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	var reLessOrEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?<=[\s]?(\$[a-z0-9_]+|[0-9.]+)`)

	var reMultiVars = regexp.MustCompile(`(?is)(\$[a-z0-9_]+\.[a-z0-9_\.]+)`)
	var reArrayVars = regexp.MustCompile(`(?is)(\$[a-z0-9_]+\[[^\]]+\][a-z0-9_\]\[\$\.]*)`)

	// Массивы, доступ через квадратные скобки
	for reArrayVars.MatchString(expr) {
		expr = reArrayVars.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var fields = strings.Split(s, "[")
			for i := 0; i < len(fields); i++ {
				fields[i] = strings.TrimRight(fields[i], "]")
			}
			functions = append(functions, func(vars *xVarCollection) []byte {
				vars.setVar(varName, vars.getMultiVar(fields).toInterface())
				return nil
			})
			return varName
		})
	}

	// Структуры, доступ через точку
	for reMultiVars.MatchString(expr) {
		expr = reMultiVars.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var fields = strings.Split(s, ".")
			functions = append(functions, func(vars *xVarCollection) []byte {
				vars.setVar(varName, vars.getMultiVar(fields).toInterface())
				return nil
			})
			return varName
		})
	}

	// $var++
	for rePlusPlus.MatchString(expr) {
		expr = rePlusPlus.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = strings.TrimSpace(strings.Split(s, "++")[0])
			return varName + " = " + varName + " + 1"
		})
	}

	// $var--
	for reMinusMinus.MatchString(expr) {
		expr = reMinusMinus.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = strings.TrimSpace(strings.Split(s, "--")[0])
			return varName + " = " + varName + " - 1"
		})
	}

	// Умножение
	for reMultiple.MatchString(expr) {
		expr = reMultiple.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "*")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toFloat()*secondValue.toFloat())
				return nil
			})
			return varName
		})
	}

	//Деление
	for reDivision.MatchString(expr) {
		expr = reDivision.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "/")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				sVal := secondValue.toFloat()
				if sVal == 0 {
					vars.setVar(varName, "Error: divide by zero")
				} else {
					vars.setVar(varName, firstValue.toFloat()/sVal)
				}
				return nil
			})
			return varName
		})
	}

	//Остаток от деления
	for reDivisionRemainder.MatchString(expr) {
		expr = reDivisionRemainder.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "%")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				sVal := secondValue.toInt()
				if sVal == 0 {
					vars.setVar(varName, "Error: divide by zero")
				} else {
					vars.setVar(varName, firstValue.toInt()%sVal)
				}
				return nil
			})
			return varName
		})
	}

	// Сложение
	for reAddition.MatchString(expr) {
		expr = reAddition.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "+")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				if firstValue.vType == varTypeString || secondValue.vType == varTypeString {
					vars.setVar(varName, firstValue.toString()+secondValue.toString())
				} else {
					vars.setVar(varName, firstValue.toFloat()+secondValue.toFloat())
				}
				return nil
			})
			return varName
		})
	}

	// Вычитание
	for reSubtraction.MatchString(expr) {
		expr = reSubtraction.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "-")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toFloat()-secondValue.toFloat())
				return nil
			})
			return varName
		})
	}

	// Равно
	for reEqual.MatchString(expr) {
		expr = reEqual.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "==")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toString() == secondValue.toString())
				return nil
			})
			return varName
		})
	}

	// НЕ Равно
	for reNotEqual.MatchString(expr) {
		expr = reNotEqual.ReplaceAllStringFunc(expr, func(s string) string {
			var sep = "!="
			if !strings.Contains(s, sep) {
				sep = "<>"
			}
			var varName = newVarName()
			var arr = strings.Split(s, sep)
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toString() != secondValue.toString())
				return nil
			})
			return varName
		})
	}

	// Больше
	for reMore.MatchString(expr) {
		expr = reMore.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, ">")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toFloat() > secondValue.toFloat())
				return nil
			})
			return varName
		})
	}

	// Меньше
	for reLess.MatchString(expr) {
		expr = reLess.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "<")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toFloat() < secondValue.toFloat())
				return nil
			})
			return varName
		})
	}

	// Меньше или равно
	for reLessOrEqual.MatchString(expr) {
		expr = reLessOrEqual.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "<=")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toFloat() <= secondValue.toFloat())
				return nil
			})
			return varName
		})
	}

	// Больше или равно
	for reMoreOrEqual.MatchString(expr) {
		expr = reMoreOrEqual.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, ">=")
			var firstVar = strings.TrimSpace(arr[0])
			var secondVar = strings.TrimSpace(arr[1])

			functions = append(functions, func(vars *xVarCollection) []byte {
				var firstValue, secondValue *xVar
				if strings.HasPrefix(firstVar, "$") {
					firstValue = vars.getVar(firstVar)
				} else {
					firstValue = xVarInit("", firstVar)
				}
				if strings.HasPrefix(secondVar, "$") {
					secondValue = vars.getVar(secondVar)
				} else {
					secondValue = xVarInit("", secondVar)
				}
				vars.setVar(varName, firstValue.toFloat() >= secondValue.toFloat())
				return nil
			})
			return varName
		})
	}

	// Присвоение
	for reAssign.MatchString(expr) {
		expr = reAssign.ReplaceAllStringFunc(expr, func(s string) string {
			var arr = strings.SplitN(s, "=", 2)
			var varName = strings.TrimSpace(arr[0])
			var value = x.exec([]rune(strings.TrimSpace(arr[1])))
			functions = append(functions, func(vars *xVarCollection) []byte {
				vars.setVar(varName, value(vars).toInterface())
				return nil
			})
			return varName
		})
	}

	expr = strings.TrimSpace(expr)

	// На этом этапе должен остаться лишь один фрагмент выражения имя переменной
	var result func(vars *xVarCollection) *xVar
	if strings.HasPrefix(expr, "$") {
		result = func(vars *xVarCollection) *xVar {
			return vars.getVar(expr)
		}
	} else {
		result = func(vars *xVarCollection) *xVar {
			return xVarInit("", expr)
		}
	}

	return func(vars *xVarCollection) *xVar {
		for _, fn := range functions {
			fn(vars)
		}
		return result(vars)
	}
}
