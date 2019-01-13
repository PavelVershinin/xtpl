package xtpl

import (
	"io"
	"sync"
)

type XtplCollection struct {
	m             sync.Mutex
	collection    map[string]*xtpl
	functions     map[string]interface{}
	viewsPath     string
	viewExtension string
	cyclesLimit   uint
	debug         bool
}

func NewCollection(viewsPath, viewExtension string) *XtplCollection {
	var xc = &XtplCollection{}
	xc.collection = make(map[string]*xtpl)
	xc.functions = make(map[string]interface{})
	xc.cyclesLimit = 10000
	if viewsPath != "" {
		xc.viewsPath = viewsPath
	} else {
		xc.viewsPath = "."
	}
	if viewExtension != "" {
		xc.viewExtension = viewExtension
	} else {
		xc.viewExtension = "tpl"
	}
	return xc
}

func (xc *XtplCollection) SetFunctions(functions map[string]interface{}) {
	xc.m.Lock()
	xc.functions = functions
	for fileName := range xc.collection {
		xc.collection[fileName] = xtplInit(xc, fileName)
	}
	xc.m.Unlock()
}

func (xc *XtplCollection) SetCycleLimit(limit uint) {
	xc.m.Lock()
	xc.cyclesLimit = limit
	xc.m.Unlock()
}

func (xc *XtplCollection) SetDebug(debug bool) {
	xc.debug = debug
}

func (xc *XtplCollection) View(tplPath string, data map[string]interface{}, writer io.Writer) {
	if xc.debug {
		xtplInit(xc, tplPath).run(data, writer)
		return
	}
start:
	if view, ok := xc.collection[tplPath]; ok {
		view.run(data, writer)
	} else {
		xc.m.Lock()
		xc.collection[tplPath] = xtplInit(xc, tplPath)
		xc.m.Unlock()
		goto start
	}
}
