package xtpl

import (
	"io"
	"io/ioutil"
	"log"
)

type treeNode func(vars *xVarCollection) []byte
type xtpl struct {
	tree []treeNode
}

func xtplInit(tplPath string) *xtpl {
	var xTpl = &xtpl{}
	src := xTpl.tplSource(tplPath)
	xTpl.parse(src)
	return xTpl
}

func xtplInitFromSource(src string) *xtpl {
	var xTpl = &xtpl{}
	xTpl.parse(src)
	return xTpl
}

func (x *xtpl) tplSource(tplPath string) string {
	b, err := ioutil.ReadFile(viewsPath + "/" + tplPath + "." + viewExtension)
	if err == nil {
		return string(b)
	}
	return err.Error()
}

func (x *xtpl) parse(src string) {
	x.tree = x.buildTree(x.preBuild(src), true)
}

func (x *xtpl) run(data map[string]interface{}, writer io.Writer) {
	if data == nil {
		data = make(map[string]interface{})
	}
	var vars = &xVarCollection{
		source:     data,
		overSource: map[string]interface{}{},
		variables:  map[string]*xVar{},
	}
	for _, f := range x.tree {
		if _, err := writer.Write(f(vars)); err != nil {
			log.Println(err.Error())
			return
		}
	}
}
