package xtpl

import (
	"math"
	"regexp"
	"strings"
)

var (
	regVarFuncs = regexp.MustCompile(`(?is)^\$([a-z0-9_\[\]."']+)$`)

	regAssign = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?=[\s]?(.*)`)
	regPlusPlus = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\+\+`)
	regMinusMinus = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?--`)
	regShortStyle = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?[\+\-\*\/\\\%\^]{1}=[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regMultiple = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\*[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regExponentiation = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\^[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regDivision = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?/[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regDivisionWithoutRemainder = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\\[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regDivisionRemainder = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?%[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regAddition = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\+[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regSubtraction = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?-[\s]?(\$[a-z0-9_]+|[0-9.]+)`)

	regEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?==[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regNotEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?(!=|<>)[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regMore = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?>[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regLess = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?<[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regMoreOrEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?>=[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regLessOrEqual = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?<=[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regOr = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\|\|[\s]?(\$[a-z0-9_]+|[0-9.]+)`)
	regAnd = regexp.MustCompile(`(?is)(\$[a-z0-9_]+|[0-9.]+)[\s]?\&\&[\s]?(\$[a-z0-9_]+|[0-9.]+)`)

	reMultiVars = regexp.MustCompile(`(?is)(\$[a-z0-9_]+\.[\$a-z0-9_\.]+)`)
)

func (x *xtpl) exec(src []rune) func(vars *xVarCollection) *xVar {
	var functions []treeNode

	src = trim(src)

	// Все строки в кавычках переводим в переменные
	for i := 0; i < len(src); i++ {
		openQuote := -1
		closeQuote := -1
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

	// Разбираем структуры в мапы
	for i := 0; i < len(src); i++ {
		openBracket := -1
		closeBracket := -1
		if openBracket = getOffset(src[i:], "[", "", true, false); openBracket > -1 {
			closeBracket = getOffset(src[i+openBracket+1:], "]", "", true, true)
		}
		if openBracket > -1 && closeBracket > -1 {
			value := strings.TrimSpace(string(src[i+openBracket+1 : i+openBracket+1+closeBracket]))
			ms := make(map[string]string)
			mi := make(map[string]interface{})
			for _, line := range strings.Split(value, ",") {
				if a := strings.Split(line, "=>"); len(a) == 2 {
					ms[strings.TrimSpace(a[0])] = strings.TrimSpace(a[1])
				}
			}
			varName := newVarName()
			functions = append(functions, func(vars *xVarCollection) []byte {
				for k, v := range ms {
					mi[vars.getVar(k).toString()] = vars.getVar(v).toInterface()
				}
				vars.setVar(varName, mi)
				return nil
			})
			src = append(src[:i+openBracket], append([]rune(varName), src[i+openBracket+1+closeBracket+1:]...)...)
		}
	}

	// Выгребаем пользовательские функции
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

	// Выгребаем функции переданные в переменных
	for i := 0; i < len(src); i++ {
		if src[i] == '$' {
			openBracketPosition := getOffset(src[i:], "(", "", true, false)
			if openBracketPosition == -1 {
				break
			}

			if regVarFuncs.MatchString(string(src[i : i+openBracketPosition])) {
				closeBracketPosition := getOffset(src[i+openBracketPosition:], ")", "", true, true)
				if closeBracketPosition == -1 {
					break
				}

				var varName = newVarName()
				var function = x.execFunction(src[i+openBracketPosition:i+openBracketPosition+closeBracketPosition+1], x.exec(src[i:i+openBracketPosition]))
				functions = append(functions, func(vars *xVarCollection) []byte {
					vars.setVar(varName, function(vars).toInterface())
					return nil
				})

				src = append(src[:i], append([]rune(varName), src[i+openBracketPosition+closeBracketPosition+1:]...)...)
			}
		}
	}

	// Выгребаем все вычисления в скобках
	for i := 0; i < len(src); i++ {
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

	// Инкремент
	for regPlusPlus.MatchString(expr) {
		expr = regPlusPlus.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = strings.TrimSpace(strings.Split(s, "++")[0])
			return varName + " = " + varName + " + 1"
		})
	}

	// Декремент
	for regMinusMinus.MatchString(expr) {
		expr = regMinusMinus.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = strings.TrimSpace(strings.Split(s, "--")[0])
			return varName + " = " + varName + " - 1"
		})
	}

	for regShortStyle.MatchString(expr) {
		expr = regShortStyle.ReplaceAllStringFunc(expr, func(s string) string {
			var firstVar, secondVar, operator string
			var arr []string
			switch {
			case strings.Contains(s, "+="):
				arr = strings.Split(s, "+=")
				operator = "+"
			case strings.Contains(s, "-="):
				arr = strings.Split(s, "-=")
				operator = "-"
			case strings.Contains(s, "*="):
				arr = strings.Split(s, "*=")
				operator = "*"
			case strings.Contains(s, "/="):
				arr = strings.Split(s, "/=")
				operator = "/"
			case strings.Contains(s, "%="):
				arr = strings.Split(s, "%=")
				operator = "%"
			case strings.Contains(s, "\\="):
				arr = strings.Split(s, "\\=")
				operator = "\\"
			case strings.Contains(s, "^="):
				arr = strings.Split(s, "^=")
				operator = "^"
			default:
				return s
			}
			firstVar = strings.TrimSpace(arr[0])
			secondVar = strings.TrimSpace(arr[1])
			return firstVar + " = " + firstVar + operator + secondVar
		})
	}

	// Возведение в степень
	for regExponentiation.MatchString(expr) {
		expr = regExponentiation.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "^")
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
				vars.setVar(varName, math.Pow(firstValue.toFloat(), secondValue.toFloat()))
				return nil
			})
			return varName
		})
	}

	// Умножение
	for regMultiple.MatchString(expr) {
		expr = regMultiple.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regDivision.MatchString(expr) {
		expr = regDivision.ReplaceAllStringFunc(expr, func(s string) string {
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

	//Деление без остатка
	for regDivisionWithoutRemainder.MatchString(expr) {
		expr = regDivisionWithoutRemainder.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "\\")
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
					vars.setVar(varName, int64(firstValue.toInt()/sVal))
				}
				return nil
			})
			return varName
		})
	}

	//Остаток от деления
	for regDivisionRemainder.MatchString(expr) {
		expr = regDivisionRemainder.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regAddition.MatchString(expr) {
		expr = regAddition.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regSubtraction.MatchString(expr) {
		expr = regSubtraction.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regEqual.MatchString(expr) {
		expr = regEqual.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regNotEqual.MatchString(expr) {
		expr = regNotEqual.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regMore.MatchString(expr) {
		expr = regMore.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regLess.MatchString(expr) {
		expr = regLess.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regLessOrEqual.MatchString(expr) {
		expr = regLessOrEqual.ReplaceAllStringFunc(expr, func(s string) string {
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
	for regMoreOrEqual.MatchString(expr) {
		expr = regMoreOrEqual.ReplaceAllStringFunc(expr, func(s string) string {
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

	// И
	for regAnd.MatchString(expr) {
		expr = regAnd.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "&&")
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
				vars.setVar(varName, firstValue.toBool() && secondValue.toBool())
				return nil
			})
			return varName
		})
	}

	// Или
	for regOr.MatchString(expr) {
		expr = regOr.ReplaceAllStringFunc(expr, func(s string) string {
			var varName = newVarName()
			var arr = strings.Split(s, "||")
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
				vars.setVar(varName, firstValue.toBool() || secondValue.toBool())
				return nil
			})
			return varName
		})
	}

	// Присвоение
	for regAssign.MatchString(expr) {
		expr = regAssign.ReplaceAllStringFunc(expr, func(s string) string {
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
