package xtpl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func newVarName() (varName string) {
	return "$_XR_" + strconv.FormatInt(serialID.Next(), 10)
}

type xVarCollection struct {
	source     map[string]interface{}
	overSource map[string]interface{}
	variables  map[string]*xVar
	keyList    []string
}

func (xvc *xVarCollection) len() int {
	return len(xvc.keyList)
}

func (xvc *xVarCollection) keys() []string {
	return xvc.keyList
}

func (xvc *xVarCollection) getMultiVar(fields []string) *xVar {
	var res *xVar
	res = xvc.getVar(fields[0])
	for i := 1; i < len(fields); i++ {
		field := fields[i]
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
		if value, ok := xvc.variables[varName]; ok {
			return value
		}
		if value, ok := xvc.overSource[varName]; ok {
			return xvc.toVar(varName, value)
		}
		if value, ok := xvc.source[varName]; ok {
			return xvc.toVar(varName, value)
		}
	}
	return &xVar{}
}

func (xvc *xVarCollection) setVar(varName string, value interface{}) {
	varName = strings.TrimLeft(varName, "$")
	if _, ok := xvc.overSource[varName]; ok {
		delete(xvc.variables, varName)
	} else if _, ok := xvc.source[varName]; ok {
		delete(xvc.variables, varName)
	} else {
		xvc.keyList = append(xvc.keyList, varName)
	}
	xvc.overSource[varName] = value
}

func (xvc *xVarCollection) toVar(varName string, value interface{}) *xVar {
	var xv = xVarInit(varName, value)
	xvc.variables[varName] = xv
	return xv
}

type varType uint8

const (
	_ varType = iota
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
		source:     map[string]interface{}{},
		overSource: map[string]interface{}{},
		variables:  map[string]*xVar{},
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
		if strings.HasPrefix(s, "0") && s != "0" {
			xv.vType = varTypeString
			xv.valueString = s
		} else if n, err := strconv.ParseFloat(s, 64); err == nil {
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
		if valueOf.IsValid() {
			for i := 0; i < valueOf.NumField(); i++ {
				field := valueOf.Field(i)
				if field.IsValid() && field.CanInterface() {
					xv.Collection.setVar(valueOf.Type().Field(i).Name, field.Interface())
				}
			}
			for i := 0; i < valueOf.NumMethod(); i++ {
				method := valueOf.Method(i)
				if method.IsValid() && method.CanInterface() {
					xv.Collection.setVar(valueOf.Type().Method(i).Name, method.Interface())
				}
			}
		}
	case reflect.Ptr:
		if valueOf.Elem().IsValid() && valueOf.Elem().CanInterface() {
			xv = xVarInit(name, valueOf.Elem().Interface())
		}
	case reflect.Func:
		xv.vType = varTypeFunc
		xv.valueFunc = valueOf
	case reflect.Interface:
		xv.vType = varTypeInterface
	}

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
		return fmt.Sprintf("%s", xv.valueInterface)
	case varTypeMap:
		return fmt.Sprintf("%s", xv.valueInterface)
	case varTypeInterface:
		return fmt.Sprintf("%s", xv.valueInterface)
	case varTypeFunc:
		return fmt.Sprintf("%s", xv.valueInterface)
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

func (xv *xVar) toSliceWithInterface() []interface{} {
	switch xv.vType {
	case varTypeSlice,
		varTypeMap,
		varTypeInterface:
		var s = make([]interface{}, xv.Collection.len())
		for i, key := range xv.Collection.keys() {
			s[i] = xv.Collection.getVar(key).toInterface()
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

func (xv *xVar) toMapWithInterface() map[string]interface{} {
	switch xv.vType {
	case varTypeSlice,
		varTypeMap,
		varTypeInterface:
		var m = make(map[string]interface{})
		for _, key := range xv.Collection.keys() {
			m[key] = xv.Collection.getVar(key).toInterface()
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
