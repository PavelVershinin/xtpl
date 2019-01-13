package xtpl

import (
	"fmt"
	"reflect"
)

func (x *xtpl) execUserFunction(src []rune, function interface{}, functionName string) func(vars *xVarCollection) *xVar {
	var functionOf = reflect.ValueOf(function)

	if functionOf.Kind() != reflect.Func {
		return func(vars *xVarCollection) *xVar {
			return xVarInit("", function)
		}
	}

	var closeBracketPosition = getOffset(src, ")", "", true, true)
	var lastCommaPosition = 1
	var stringArguments [][]rune
	for i := 1; i < closeBracketPosition; {
		if commaPosition := getOffset(src[i:], ",", "", true, true); commaPosition > 0 {
			stringArguments = append(stringArguments, src[i:i+commaPosition])
			i += commaPosition + 1
			lastCommaPosition = i
		} else {
			i++
		}
	}

	if lastCommaPosition < closeBracketPosition {
		stringArguments = append(stringArguments, src[lastCommaPosition:closeBracketPosition])
	}

	var argumentTypes []reflect.Kind
	var arguments []func(vars *xVarCollection) *xVar
	for i := 0; i < functionOf.Type().NumIn(); i++ {
		var argumentType = functionOf.Type().In(i).Kind()
		var argument func(vars *xVarCollection) *xVar
		if len(stringArguments) > i {
			i := i
			argument = x.exec(stringArguments[i])
		} else {
			argument = func(vars *xVarCollection) *xVar {
				return nil
			}
		}
		arguments = append(arguments, argument)
		argumentTypes = append(argumentTypes, argumentType)
	}

	return func(vars *xVarCollection) (result *xVar) {
		defer func() {
			if r := recover(); r != nil {
				result = xVarInit("", fmt.Sprint(r))
			}
		}()

		var argLen = len(arguments)
		var args = make([]reflect.Value, argLen, argLen)
		for i := 0; i < argLen; i++ {
			v := arguments[i](vars)
			switch argumentTypes[i] {
			case reflect.Bool:
				args[i] = reflect.ValueOf(v.toBool())
			case reflect.Int:
				args[i] = reflect.ValueOf(int(v.toInt()))
			case reflect.Int8:
				args[i] = reflect.ValueOf(int8(v.toInt()))
			case reflect.Int16:
				args[i] = reflect.ValueOf(int16(v.toInt()))
			case reflect.Int32:
				args[i] = reflect.ValueOf(int32(v.toInt()))
			case reflect.Int64:
				args[i] = reflect.ValueOf(v.toInt())
			case reflect.Uint:
				args[i] = reflect.ValueOf(uint(v.toInt()))
			case reflect.Uint8:
				args[i] = reflect.ValueOf(uint8(v.toInt()))
			case reflect.Uint16:
				args[i] = reflect.ValueOf(uint16(v.toInt()))
			case reflect.Uint32:
				args[i] = reflect.ValueOf(uint32(v.toInt()))
			case reflect.Uint64:
				args[i] = reflect.ValueOf(uint64(v.toInt()))
			case reflect.Float32:
				args[i] = reflect.ValueOf(float32(v.toFloat()))
			case reflect.Float64:
				args[i] = reflect.ValueOf(v.toFloat())
			case reflect.String:
				args[i] = reflect.ValueOf(v.toString())
			case reflect.Slice:
				args[i] = reflect.ValueOf(v.toSlice())
			case reflect.Map:
				args[i] = reflect.ValueOf(v.toMap())
			default:
				args[i] = reflect.ValueOf(v.toInterface())
			}
		}
		var ret = functionOf.Call(args)
		if len(ret) > 0 {
			result = xVarInit("", ret[0].Interface())
		}

		return
	}
}
