package xtpl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var varNameIter = 0

func newVarName() (varName string) {
	varNameIter++
	return "$_XR_" + strconv.Itoa(varNameIter)
}

type xVarCollection struct {
	source    map[string]interface{}
	variables []*xVar
	keyList   []string
}

func (xvc *xVarCollection) len() int {
	return len(xvc.source)
}

func (xvc *xVarCollection) keys() []string {
	if len(xvc.keyList) > 0 {
		return xvc.keyList
	}
	xvc.keyList = make([]string, len(xvc.source))
	var keys = make([]string, xvc.len())
	var i = 0
	for key := range xvc.source {
		xvc.keyList[i] = key
		keys[i] = key
		i++
	}
	return keys
}

func (xvc *xVarCollection) getMultiVar(fields []string) *xVar {
	var res *xVar
	res = xvc.getVar(fields[0])
	for i := 1; i < len(fields); i++ {
		var field = fields[i]
		if strings.HasPrefix(field, "$") {
			field = xvc.getVar(field).toString()
		}
		res = res.Collection.getVar(field)
	}
	return res
}

func (xvc *xVarCollection) getVar(varName string) *xVar {
	if xvc != nil {
		varName = strings.TrimLeft(varName, "$")
		for i := 0; i < len(xvc.variables); i++ {
			if xvc.variables[i].name == varName {
				return xvc.variables[i]
			}
		}
		if value, ok := xvc.source[varName]; ok {
			return xvc.toVar(varName, value)
		}
	}
	return &xVar{}
}

func (xvc *xVarCollection) setVar(varName string, value interface{}) {
	varName = strings.TrimLeft(varName, "$")
	xvc.source[varName] = value
	var tmp []*xVar
	for i := 0; i < len(xvc.variables); i++ {
		if xvc.variables[i].name != varName {
			tmp = append(tmp, xvc.variables[i])
		}
	}
	xvc.variables = tmp

	for i := 0; i < len(xvc.keyList); i++ {
		if xvc.keyList[i] == varName {
			return
		}
	}
	xvc.keyList = append(xvc.keyList, varName)
}

func (xvc *xVarCollection) toVar(varName string, value interface{}) *xVar {
	var xv = xVarInit(varName, value)
	xvc.variables = append(xvc.variables, xv)
	return xv
}

type varType uint8

const (
	varTypeInvalid varType = iota
	varTypeBool
	varTypeInt
	varTypeFloat
	varTypeString
	varTypeSlice
	varTypeMap
	varTypeInterface
	varTypeFunc
)

type xVar struct {
	name           string
	vType          varType
	valueBool      bool
	valueInt       int64
	valueFloat     float64
	valueString    string
	valueInterface interface{}
	valueFunc      reflect.Value
	Collection     *xVarCollection
}

func xVarInit(name string, value interface{}) *xVar {
	var xv = &xVar{}
	xv.Collection = &xVarCollection{
		source: map[string]interface{}{},
	}
	var valueOf = reflect.ValueOf(value)
	var valueKind = valueOf.Kind()

	switch valueKind {
	case reflect.Bool:
		xv.vType = varTypeBool
		xv.valueBool = value.(bool)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		xv.vType = varTypeInt
		switch valueKind {
		case reflect.Int:
			xv.valueInt = int64(value.(int))
		case reflect.Int8:
			xv.valueInt = int64(value.(int8))
		case reflect.Int16:
			xv.valueInt = int64(value.(int16))
		case reflect.Int32:
			xv.valueInt = int64(value.(int32))
		case reflect.Int64:
			xv.valueInt = value.(int64)
		case reflect.Uint:
			xv.valueInt = int64(value.(uint))
		case reflect.Uint8:
			xv.valueInt = int64(value.(uint8))
		case reflect.Uint16:
			xv.valueInt = int64(value.(uint16))
		case reflect.Uint32:
			xv.valueInt = int64(value.(uint32))
		case reflect.Uint64:
			xv.valueInt = int64(value.(uint64))
		}
	case reflect.Float32,
		reflect.Float64:
		xv.vType = varTypeFloat
		switch valueKind {
		case reflect.Float32:
			xv.valueFloat = float64(value.(float32))
		case reflect.Float64:
			xv.valueFloat = value.(float64)
		}
	case reflect.String:
		s := fmt.Sprintf("%s", value)
		if n, err := strconv.ParseFloat(s, 64); err == nil {
			xv.vType = varTypeFloat
			xv.valueFloat = n
		} else {
			xv.vType = varTypeString
			xv.valueString = s
		}
	case reflect.Array,
		reflect.Slice:
		xv.vType = varTypeSlice
		for i := 0; i < valueOf.Len(); i++ {
			xv.Collection.setVar(strconv.Itoa(i), valueOf.Index(i).Interface())
		}
	case reflect.Map:
		xv.vType = varTypeMap
		for _, key := range valueOf.MapKeys() {
			xv.Collection.setVar(key.String(), valueOf.MapIndex(key).Interface())
		}
	case reflect.Struct:
		xv.vType = varTypeMap
		for i := 0; i < valueOf.NumField(); i++ {
			field := valueOf.Field(i)
			if field.IsValid() && field.CanInterface() {
				xv.Collection.setVar(valueOf.Type().Field(i).Name, field.Interface())
			}
		}
	case reflect.Ptr:
		xv.vType = varTypeMap
		if valueOf.Elem().IsValid() {
			for i := 0; i < valueOf.Elem().NumField(); i++ {
				field := valueOf.Elem().Field(i)
				if field.IsValid() && field.CanInterface() {
					xv.Collection.setVar(valueOf.Elem().Type().Field(i).Name, field.Interface())
				}
			}
		}
	case reflect.Func:
		xv.vType = varTypeFunc
		xv.valueFunc = valueOf
	case reflect.Interface:
		xv.vType = varTypeInterface
	}

	xv.name = name
	xv.valueInterface = value
	return xv
}

func (xv *xVar) toBool() bool {
	switch xv.vType {
	case varTypeBool:
		return xv.valueBool
	case varTypeInt:
		return xv.valueInt != 0
	case varTypeFloat:
		return xv.valueFloat != 0
	case varTypeString:
		b, _ := strconv.ParseBool(xv.valueString)
		return b
	case varTypeSlice:
		return xv.Collection.len() > 0
	case varTypeMap:
		return xv.Collection.len() > 0
	case varTypeInterface:
		return xv.valueInterface != nil
	case varTypeFunc:
		return !xv.valueFunc.IsNil()
	}
	return false
}

func (xv *xVar) toInt() int64 {
	switch xv.vType {
	case varTypeBool:
		if xv.valueBool {
			return 1
		}
		return 0
	case varTypeInt:
		return xv.valueInt
	case varTypeFloat:
		return int64(xv.valueFloat)
	case varTypeString:
		i, _ := strconv.ParseInt(xv.valueString, 10, 64)
		return i
	case varTypeSlice:
		return int64(xv.Collection.len())
	case varTypeMap:
		return int64(xv.Collection.len())
	case varTypeInterface:
		if xv.valueInterface != nil {
			return 1
		}
		return 0
	case varTypeFunc:
		if !xv.valueFunc.IsNil() {
			return 1
		}
		return 0
	}
	return 0
}

func (xv *xVar) toFloat() float64 {
	switch xv.vType {
	case varTypeBool:
		if xv.valueBool {
			return 1
		}
		return 0
	case varTypeInt:
		return float64(xv.valueInt)
	case varTypeFloat:
		return xv.valueFloat
	case varTypeString:
		f, _ := strconv.ParseFloat(xv.valueString, 64)
		return f
	case varTypeSlice:
		return float64(xv.Collection.len())
	case varTypeMap:
		return float64(xv.Collection.len())
	case varTypeInterface:
		if xv.valueInterface != nil {
			return 1
		}
		return 0
	case varTypeFunc:
		if !xv.valueFunc.IsNil() {
			return 1
		}
		return 0
	}
	return 0
}

func (xv *xVar) toString() string {
	switch xv.vType {
	case varTypeBool:
		return strconv.FormatBool(xv.valueBool)
	case varTypeInt:
		return strconv.FormatInt(xv.valueInt, 10)
	case varTypeFloat:
		return strconv.FormatFloat(xv.valueFloat, 'f', -1, 64)
	case varTypeString:
		return xv.valueString
	case varTypeSlice:
		return fmt.Sprintf("%#v", xv.valueInterface)
	case varTypeMap:
		return fmt.Sprintf("%#v", xv.valueInterface)
	case varTypeInterface:
		return fmt.Sprintf("%#v", xv.valueInterface)
	case varTypeFunc:
		return fmt.Sprintf("%#v", xv.valueInterface)
	}
	return ""
}

func (xv *xVar) toBytes() []byte {
	return []byte(xv.toString())
}

func (xv *xVar) toSlice() []*xVar {
	switch xv.vType {
	case varTypeSlice,
		varTypeMap,
		varTypeInterface:
		var s = make([]*xVar, xv.Collection.len())
		for i, key := range xv.Collection.keys() {
			s[i] = xv.Collection.getVar(key)
		}
		return s
	}
	return nil
}

func (xv *xVar) toMap() map[string]*xVar {
	switch xv.vType {
	case varTypeSlice,
		varTypeMap,
		varTypeInterface:
		var m = make(map[string]*xVar)
		for _, key := range xv.Collection.keys() {
			m[key] = xv.Collection.getVar(key)
		}
		return m
	}
	return nil
}

func (xv *xVar) toInterface() interface{} {
	return xv.valueInterface
}

func (xv *xVar) toFunc() (call func(in []reflect.Value) []reflect.Value, arguments []reflect.Kind) {
	switch xv.vType {
	case varTypeFunc:
		call = xv.valueFunc.Call
		for i := 0; i < xv.valueFunc.Type().NumIn(); i++ {
			arguments = append(arguments, xv.valueFunc.Type().In(i).Kind())
		}
	}
	return
}
