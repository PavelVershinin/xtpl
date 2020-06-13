package xtpl

import (
	"fmt"
	"reflect"
)

func (x *xtpl) execFunction(src []rune, function func(vars *xVarCollection) *xVar) func(vars *xVarCollection) *xVar {
	var closeBracketPosition = getOffset(src, ")", "", true, true)
	var lastCommaPosition = 1
	var arguments []func(vars *xVarCollection) *xVar
	for i := 1; i < closeBracketPosition; {
		if commaPosition := getOffset(src[i:], ",", "", true, true); commaPosition > 0 {
			arguments = append(arguments, x.exec(src[i:i+commaPosition]))
			i += commaPosition + 1
			lastCommaPosition = i
		} else {
			i++
		}
	}

	if lastCommaPosition < closeBracketPosition {
		arguments = append(arguments, x.exec(src[lastCommaPosition:closeBracketPosition]))
	}

	return func(vars *xVarCollection) (result *xVar) {
		defer func() {
			if r := recover(); r != nil {
				result = xVarInit("", fmt.Sprint(r))
			}
		}()
		call, argumentTypes := function(vars).toFunc()

		var args = prepareArguments(arguments, argumentTypes, vars)
		var ret = call(args)
		if len(ret) > 0 {
			result = xVarInit("", ret[0].Interface())
		}

		return
	}
}

func (x *xtpl) execUserFunction(src []rune, function interface{}) func(vars *xVarCollection) *xVar {
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

		var args = prepareArguments(arguments, argumentTypes, vars)
		var ret = functionOf.Call(args)
		if len(ret) > 0 {
			result = xVarInit("", ret[0].Interface())
		}

		return
	}
}

func prepareArguments(arguments []func(vars *xVarCollection) *xVar, argumentTypes []reflect.Kind, vars *xVarCollection) []reflect.Value {
	var interfaceOf = func(t reflect.Kind, v *xVar) interface{} {
		switch t {
		case reflect.Bool:
			return v.toBool()
		case reflect.Int:
			return int(v.toInt())
		case reflect.Int8:
			return int8(v.toInt())
		case reflect.Int16:
			return int16(v.toInt())
		case reflect.Int32:
			return int32(v.toInt())
		case reflect.Int64:
			return v.toInt()
		case reflect.Uint:
			return uint(v.toInt())
		case reflect.Uint8:
			return uint8(v.toInt())
		case reflect.Uint16:
			return uint16(v.toInt())
		case reflect.Uint32:
			return uint32(v.toInt())
		case reflect.Uint64:
			return uint64(v.toInt())
		case reflect.Float32:
			return float32(v.toFloat())
		case reflect.Float64:
			return v.toFloat()
		case reflect.String:
			return v.toString()
		case reflect.Slice:
			return v.toSlice()
		case reflect.Map:
			return v.toMap()
		default:
			return v.toInterface()
		}
	}
	var argLen = len(argumentTypes)
	var args = make([]reflect.Value, argLen)
	for i := 0; i < argLen; i++ {
		v := arguments[i](vars)
		args[i] = reflect.ValueOf(interfaceOf(argumentTypes[i], v))
	}
	return args
}
