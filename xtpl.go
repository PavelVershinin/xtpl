package xtpl

import (
	"io"
	"io/fs"
	"os"
)

type treeNode func(vars *xVarCollection) []byte
type xtpl struct {
	tree   []treeNode
	errors *errors
}

func xtplInit(tplPath string) *xtpl {
	var xTpl = &xtpl{}
	xTpl.errors = &errors{}
	src := xTpl.tplSource(tplPath)
	xTpl.parse(src)
	return xTpl
}

func xtplInitFromSource(src string) *xtpl {
	var xTpl = &xtpl{}
	xTpl.errors = &errors{}
	xTpl.parse(src)
	return xTpl
}

func (x *xtpl) tplSource(tplPath string) string {
	var file fs.File
	var err error

	path := viewsPath + string(os.PathSeparator) + tplPath + "." + viewExtension

	if embeddedFS != nil {
		file, err = embeddedFS.Open(path)
	} else {
		file, err = os.Open(path)
	}

	if err != nil {
		return x.errors.Add(err)
	}

	defer func() {
		_ = file.Close()
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		return x.errors.Add(err)
	}

	return string(b)
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
			x.errors.Add(err)
			return
		}
	}
}
