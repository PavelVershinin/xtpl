package xtpl

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"
)

func TestNewCollection(t *testing.T) {
	var testData = map[string]interface{}{
		"page_title": "Тестовая страница",
		"strings": []string{"Первый", "Второй", "Третий"},
		"numbers": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		"structs": []struct {
			ID int
			Name string
			Value string
		}{
			{
				1,
				"my name",
				"my value",
			},
			{
				2,
				"my name 2",
				"my value 2",
			},
		},
	}

	var collection = NewCollection("./templates_test", "tpl")
	collection.SetCycleLimit(100)
	collection.SetDebug(false)
	collection.SetFunctions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	var buff = &bytes.Buffer{}
	collection.View("index", testData, buff)

	var normal, err = ioutil.ReadFile("./templates_test/index.html")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !bytes.Equal(normal, buff.Bytes()) {
		t.Errorf("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
	}
}