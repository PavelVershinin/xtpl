package xtpl

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"
	"time"
)

var testData1 = map[string]interface{}{
	"page_title": "Тестовая страница",
	"strings":    []string{"Первый", "Второй", "Третий"},
	"numbers":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	"structs": []struct {
		ID       int
		Name     string
		Value    string
		Function func(s string, i int) string
	}{
		{
			1,
			"my name",
			"my value",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " дней"
			},
		},
		{
			2,
			"my name 2",
			"my value 2",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " дней"
			},
		},
	},
}

var testData2 = map[string]interface{}{
	"page_title": "Тестовая страница 2",
	"strings":    []string{"Первый 2", "Второй 2", "Третий 2"},
	"numbers":    []int{0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	"structs": []struct {
		ID       int
		Name     string
		Value    string
		Function func(s string, i int) string
	}{
		{
			3,
			"my name 2",
			"my value 2",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " часов"
			},
		},
		{
			4,
			"my name 2 2",
			"my value 2 2",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " часов"
			},
		},
	},
}


func TestXtplCollection_View(t *testing.T) {
	var collection = NewCollection("./templates_test", "tpl")
	collection.SetCycleLimit(100)
	collection.SetDebug(false)
	collection.SetFunctions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	var wg = sync.WaitGroup{}
	wg.Add(200)

	for  i := 0; i < 100; i++ {
		go func() {
			var buff = &bytes.Buffer{}
			collection.View("index", testData1, buff)

			var normal, err = ioutil.ReadFile("./templates_test/index1.html")
			if err != nil {
				t.Errorf(err.Error())
			}
			if !bytes.Equal(normal, buff.Bytes()) {
				t.Errorf("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
			}
			wg.Done()
		}()
		go func() {
			var buff = &bytes.Buffer{}
			collection.View("index", testData2, buff)

			var normal, err = ioutil.ReadFile("./templates_test/index2.html")
			if err != nil {
				t.Errorf(err.Error())
			}
			if !bytes.Equal(normal, buff.Bytes()) {
				t.Errorf("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkView(b *testing.B) {
	var collection = NewCollection("./templates_test", "tpl")
	collection.SetCycleLimit(100)
	collection.SetDebug(false)
	collection.SetFunctions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	for i := 0; i < b.N; i++ {
		var buff = &bytes.Buffer{}
		collection.View("index", testData1, buff)
	}
}