package xtpl

import (
	"io"
	"sync"
)

// XtplCollection
type XtplCollection struct {
	m             sync.RWMutex
	collection    map[string]*xtpl
	functions     map[string]interface{}
	viewsPath     string
	viewExtension string
	cyclesLimit   uint
	debug         bool
}

// NewCollection Создание новой коллекции шаблонов
func NewCollection(viewsPath, viewExtension string) *XtplCollection {
	var xc = &XtplCollection{}
	xc.collection = make(map[string]*xtpl)
	xc.functions = xc.defaultFunctions()
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

// SetFunctions Загрузка пользовательских функций в шаблоны.
func (xc *XtplCollection) SetFunctions(functions map[string]interface{}) {
	xc.m.Lock()
	xc.functions = xc.defaultFunctions()
	for name, function := range functions {
		xc.functions[name] = function
	}
	for fileName := range xc.collection {
		xc.collection[fileName] = xtplInit(xc, fileName)
	}
	xc.m.Unlock()
}

// SetCycleLimit Установка ограничения, на максимальное количество итераций в циклах. По умолчанию 10000
func (xc *XtplCollection) SetCycleLimit(limit uint) {
	xc.m.Lock()
	xc.cyclesLimit = limit
	xc.m.Unlock()
}

// SetDebug Переключение в режим отладки.
// В этом режиме, все изменения в шаблонах подхватываются налету, однако обработка шаблона занимет больше времени
func (xc *XtplCollection) SetDebug(debug bool) {
	xc.debug = debug
}

// View Обработка шаблона
func (xc *XtplCollection) View(tplPath string, data map[string]interface{}, writer io.Writer) {
	if xc.debug {
		xtplInit(xc, tplPath).run(data, writer)
		return
	}

	xc.m.RLock()
	view, ok := xc.collection[tplPath]
	xc.m.RUnlock()

	if !ok {
		xc.m.Lock()
		view = xtplInit(xc, tplPath)
		xc.collection[tplPath] = view
		xc.m.Unlock()
	}
	view.run(data, writer)
}
