package xtpl_test

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/PavelVershinin/xtpl"
)

var testData1 = map[string]interface{}{
	"page_title": "Тестовая страница",
	"strings":    []string{"Первый", "Второй", "Третий", "048465450"},
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
	"strings":    []string{"Первый 2", "Второй 2", "Третий 2", "07869786"},
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

var result1 = `<!DOCTYPE html>
<html lang="ru">
<head>
    <title>Тестовая страница</title>
</head>
<body>
    <h1>Тестовая страница<h1/>

    <h2>Вывод переменных в шаблон</h2>
    <ul>
        <li>&lt;b&gt;Экранированный вывод&lt;/b&gt;</li>
        <li><b>НЕ экранированный вывод</b></li>
    </ul>

    <h2>Объявление переменных и и прочие операции, без вывода результата в шаблон</h2>
    <ul>
        <li>Тут объявляю переменную и присваиваю ей результат математической операции, в шаблоне не вывожу </li>
        <li>Тут выведу значение сохранённое в переменной 50</li>
    </ul>

    <h2>Цикл for</h2>
    <ul>
    
        <li>
            
                В переменной $strings, под индексом 0 содержится значение Первый
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 1 содержится значение Второй
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 2 содержится значение Третий
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 3 содержится значение 048465450
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 4 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 5 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 6 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 7 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 8 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 9 нет значения
            
        </li>
    
    </ul>

    <h2>Цикл foreach</h2>
    <ul>
    
        <li>$strings[0] == Первый</li>
    
        <li>$strings[1] == Второй</li>
    
        <li>$strings[2] == Третий</li>
    
        <li>$strings[3] == 048465450</li>
    
    </ul>

    <h2>Условия</h2>
    <ul>
    
        
            
            
                <li>1 <= 1</li>
            
        
            
            
                <li>1 < 2</li>
            
        
            
            
                <li>1 < 3</li>
            
        
            
            
                <li>1 < 4</li>
            
        
            
            
                <li>1 < 5</li>
            
        
            
            
                <li>1 < 6</li>
            
        
            
            
                <li>1 < 7</li>
            
        
            
            
                <li>1 < 8</li>
            
        
            
                <li>9 > 8 && 1 < 8 || 9 < 1 && 1 > 1</li>
            
            
                <li>1 < 9</li>
            
        
            
            
                <li>1 >= 0</li>
            
        
    
        
            
            
                <li>2 >= 1</li>
            
        
            
            
                <li>2 <= 2</li>
            
        
            
            
                <li>2 < 3</li>
            
        
            
            
                <li>2 < 4</li>
            
        
            
            
                <li>2 < 5</li>
            
        
            
            
                <li>2 < 6</li>
            
        
            
            
                <li>2 < 7</li>
            
        
            
            
                <li>2 < 8</li>
            
        
            
                <li>9 > 8 && 2 < 8 || 9 < 1 && 2 > 1</li>
            
            
                <li>2 < 9</li>
            
        
            
                <li>0 > 8 && 2 < 8 || 0 < 1 && 2 > 1</li>
            
            
                <li>2 >= 0</li>
            
        
    
        
            
            
                <li>3 >= 1</li>
            
        
            
            
                <li>3 >= 2</li>
            
        
            
            
                <li>3 <= 3</li>
            
        
            
            
                <li>3 < 4</li>
            
        
            
            
                <li>3 < 5</li>
            
        
            
            
                <li>3 < 6</li>
            
        
            
            
                <li>3 < 7</li>
            
        
            
            
                <li>3 < 8</li>
            
        
            
                <li>9 > 8 && 3 < 8 || 9 < 1 && 3 > 1</li>
            
            
                <li>3 < 9</li>
            
        
            
                <li>0 > 8 && 3 < 8 || 0 < 1 && 3 > 1</li>
            
            
                <li>3 >= 0</li>
            
        
    
        
            
            
                <li>4 >= 1</li>
            
        
            
            
                <li>4 >= 2</li>
            
        
            
            
                <li>4 >= 3</li>
            
        
            
            
                <li>4 <= 4</li>
            
        
            
            
                <li>4 < 5</li>
            
        
            
            
                <li>4 < 6</li>
            
        
            
            
                <li>4 < 7</li>
            
        
            
            
                <li>4 < 8</li>
            
        
            
                <li>9 > 8 && 4 < 8 || 9 < 1 && 4 > 1</li>
            
            
                <li>4 < 9</li>
            
        
            
                <li>0 > 8 && 4 < 8 || 0 < 1 && 4 > 1</li>
            
            
                <li>4 >= 0</li>
            
        
    
        
            
            
                <li>5 >= 1</li>
            
        
            
            
                <li>5 >= 2</li>
            
        
            
            
                <li>5 >= 3</li>
            
        
            
            
                <li>5 >= 4</li>
            
        
            
            
                <li>5 <= 5</li>
            
        
            
            
                <li>5 < 6</li>
            
        
            
            
                <li>5 < 7</li>
            
        
            
            
                <li>5 < 8</li>
            
        
            
                <li>9 > 8 && 5 < 8 || 9 < 1 && 5 > 1</li>
            
            
                <li>5 < 9</li>
            
        
            
                <li>0 > 8 && 5 < 8 || 0 < 1 && 5 > 1</li>
            
            
                <li>5 >= 0</li>
            
        
    
        
            
            
                <li>6 >= 1</li>
            
        
            
            
                <li>6 >= 2</li>
            
        
            
            
                <li>6 >= 3</li>
            
        
            
            
                <li>6 >= 4</li>
            
        
            
            
                <li>6 >= 5</li>
            
        
            
            
                <li>6 <= 6</li>
            
        
            
            
                <li>6 < 7</li>
            
        
            
            
                <li>6 < 8</li>
            
        
            
                <li>9 > 8 && 6 < 8 || 9 < 1 && 6 > 1</li>
            
            
                <li>6 < 9</li>
            
        
            
                <li>0 > 8 && 6 < 8 || 0 < 1 && 6 > 1</li>
            
            
                <li>6 >= 0</li>
            
        
    
        
            
            
                <li>7 >= 1</li>
            
        
            
            
                <li>7 >= 2</li>
            
        
            
            
                <li>7 >= 3</li>
            
        
            
            
                <li>7 >= 4</li>
            
        
            
            
                <li>7 >= 5</li>
            
        
            
            
                <li>7 >= 6</li>
            
        
            
            
                <li>7 <= 7</li>
            
        
            
            
                <li>7 < 8</li>
            
        
            
                <li>9 > 8 && 7 < 8 || 9 < 1 && 7 > 1</li>
            
            
                <li>7 < 9</li>
            
        
            
                <li>0 > 8 && 7 < 8 || 0 < 1 && 7 > 1</li>
            
            
                <li>7 >= 0</li>
            
        
    
        
            
            
                <li>8 >= 1</li>
            
        
            
            
                <li>8 >= 2</li>
            
        
            
            
                <li>8 >= 3</li>
            
        
            
            
                <li>8 >= 4</li>
            
        
            
            
                <li>8 >= 5</li>
            
        
            
            
                <li>8 >= 6</li>
            
        
            
            
                <li>8 >= 7</li>
            
        
            
            
                <li>8 <= 8</li>
            
        
            
            
                <li>8 < 9</li>
            
        
            
                <li>0 > 8 && 8 < 8 || 0 < 1 && 8 > 1</li>
            
            
                <li>8 >= 0</li>
            
        
    
        
            
            
                <li>9 >= 1</li>
            
        
            
            
                <li>9 >= 2</li>
            
        
            
            
                <li>9 >= 3</li>
            
        
            
            
                <li>9 >= 4</li>
            
        
            
            
                <li>9 >= 5</li>
            
        
            
            
                <li>9 >= 6</li>
            
        
            
            
                <li>9 >= 7</li>
            
        
            
            
                <li>9 >= 8</li>
            
        
            
            
                <li>9 <= 9</li>
            
        
            
                <li>0 > 8 && 9 < 8 || 0 < 1 && 9 > 1</li>
            
            
                <li>9 >= 0</li>
            
        
    
        
            
            
                <li>0 < 1</li>
            
        
            
            
                <li>0 < 2</li>
            
        
            
            
                <li>0 < 3</li>
            
        
            
            
                <li>0 < 4</li>
            
        
            
            
                <li>0 < 5</li>
            
        
            
            
                <li>0 < 6</li>
            
        
            
            
                <li>0 < 7</li>
            
        
            
            
                <li>0 < 8</li>
            
        
            
                <li>9 > 8 && 0 < 8 || 9 < 1 && 0 > 1</li>
            
            
                <li>0 < 9</li>
            
        
            
            
                <li>0 <= 0</li>
            
        
    
    </ul>

    <h2>Математические операции</h2>
    <ul>
    
        
            <li>1 + 1 = 2</li>
            <li>1 - 1 = 0</li>
            <li>1 * 1 = 1</li>
            <li>1 / 1 = 1</li>
            <li>1 \ 1 = 1</li> 
            <li>1 % 1 = 0</li>
            <li>1 ^ 1 = 1</li>
        
            <li>1 + 2 = 3</li>
            <li>1 - 2 = -1</li>
            <li>1 * 2 = 2</li>
            <li>1 / 2 = 0.5</li>
            <li>1 \ 2 = 0</li> 
            <li>1 % 2 = 1</li>
            <li>1 ^ 2 = 1</li>
        
            <li>1 + 3 = 4</li>
            <li>1 - 3 = -2</li>
            <li>1 * 3 = 3</li>
            <li>1 / 3 = 0.3333333333333333</li>
            <li>1 \ 3 = 0</li> 
            <li>1 % 3 = 1</li>
            <li>1 ^ 3 = 1</li>
        
            <li>1 + 4 = 5</li>
            <li>1 - 4 = -3</li>
            <li>1 * 4 = 4</li>
            <li>1 / 4 = 0.25</li>
            <li>1 \ 4 = 0</li> 
            <li>1 % 4 = 1</li>
            <li>1 ^ 4 = 1</li>
        
            <li>1 + 5 = 6</li>
            <li>1 - 5 = -4</li>
            <li>1 * 5 = 5</li>
            <li>1 / 5 = 0.2</li>
            <li>1 \ 5 = 0</li> 
            <li>1 % 5 = 1</li>
            <li>1 ^ 5 = 1</li>
        
            <li>1 + 6 = 7</li>
            <li>1 - 6 = -5</li>
            <li>1 * 6 = 6</li>
            <li>1 / 6 = 0.16666666666666666</li>
            <li>1 \ 6 = 0</li> 
            <li>1 % 6 = 1</li>
            <li>1 ^ 6 = 1</li>
        
            <li>1 + 7 = 8</li>
            <li>1 - 7 = -6</li>
            <li>1 * 7 = 7</li>
            <li>1 / 7 = 0.14285714285714285</li>
            <li>1 \ 7 = 0</li> 
            <li>1 % 7 = 1</li>
            <li>1 ^ 7 = 1</li>
        
            <li>1 + 8 = 9</li>
            <li>1 - 8 = -7</li>
            <li>1 * 8 = 8</li>
            <li>1 / 8 = 0.125</li>
            <li>1 \ 8 = 0</li> 
            <li>1 % 8 = 1</li>
            <li>1 ^ 8 = 1</li>
        
            <li>1 + 9 = 10</li>
            <li>1 - 9 = -8</li>
            <li>1 * 9 = 9</li>
            <li>1 / 9 = 0.1111111111111111</li>
            <li>1 \ 9 = 0</li> 
            <li>1 % 9 = 1</li>
            <li>1 ^ 9 = 1</li>
        
            <li>1 + 0 = 1</li>
            <li>1 - 0 = 1</li>
            <li>1 * 0 = 0</li>
            <li>1 / 0 = Error: divide by zero</li>
            <li>1 \ 0 = Error: divide by zero</li> 
            <li>1 % 0 = Error: divide by zero</li>
            <li>1 ^ 0 = 1</li>
        
    
        
            <li>2 + 1 = 3</li>
            <li>2 - 1 = 1</li>
            <li>2 * 1 = 2</li>
            <li>2 / 1 = 2</li>
            <li>2 \ 1 = 2</li> 
            <li>2 % 1 = 0</li>
            <li>2 ^ 1 = 2</li>
        
            <li>2 + 2 = 4</li>
            <li>2 - 2 = 0</li>
            <li>2 * 2 = 4</li>
            <li>2 / 2 = 1</li>
            <li>2 \ 2 = 1</li> 
            <li>2 % 2 = 0</li>
            <li>2 ^ 2 = 4</li>
        
            <li>2 + 3 = 5</li>
            <li>2 - 3 = -1</li>
            <li>2 * 3 = 6</li>
            <li>2 / 3 = 0.6666666666666666</li>
            <li>2 \ 3 = 0</li> 
            <li>2 % 3 = 2</li>
            <li>2 ^ 3 = 8</li>
        
            <li>2 + 4 = 6</li>
            <li>2 - 4 = -2</li>
            <li>2 * 4 = 8</li>
            <li>2 / 4 = 0.5</li>
            <li>2 \ 4 = 0</li> 
            <li>2 % 4 = 2</li>
            <li>2 ^ 4 = 16</li>
        
            <li>2 + 5 = 7</li>
            <li>2 - 5 = -3</li>
            <li>2 * 5 = 10</li>
            <li>2 / 5 = 0.4</li>
            <li>2 \ 5 = 0</li> 
            <li>2 % 5 = 2</li>
            <li>2 ^ 5 = 32</li>
        
            <li>2 + 6 = 8</li>
            <li>2 - 6 = -4</li>
            <li>2 * 6 = 12</li>
            <li>2 / 6 = 0.3333333333333333</li>
            <li>2 \ 6 = 0</li> 
            <li>2 % 6 = 2</li>
            <li>2 ^ 6 = 64</li>
        
            <li>2 + 7 = 9</li>
            <li>2 - 7 = -5</li>
            <li>2 * 7 = 14</li>
            <li>2 / 7 = 0.2857142857142857</li>
            <li>2 \ 7 = 0</li> 
            <li>2 % 7 = 2</li>
            <li>2 ^ 7 = 128</li>
        
            <li>2 + 8 = 10</li>
            <li>2 - 8 = -6</li>
            <li>2 * 8 = 16</li>
            <li>2 / 8 = 0.25</li>
            <li>2 \ 8 = 0</li> 
            <li>2 % 8 = 2</li>
            <li>2 ^ 8 = 256</li>
        
            <li>2 + 9 = 11</li>
            <li>2 - 9 = -7</li>
            <li>2 * 9 = 18</li>
            <li>2 / 9 = 0.2222222222222222</li>
            <li>2 \ 9 = 0</li> 
            <li>2 % 9 = 2</li>
            <li>2 ^ 9 = 512</li>
        
            <li>2 + 0 = 2</li>
            <li>2 - 0 = 2</li>
            <li>2 * 0 = 0</li>
            <li>2 / 0 = Error: divide by zero</li>
            <li>2 \ 0 = Error: divide by zero</li> 
            <li>2 % 0 = Error: divide by zero</li>
            <li>2 ^ 0 = 1</li>
        
    
        
            <li>3 + 1 = 4</li>
            <li>3 - 1 = 2</li>
            <li>3 * 1 = 3</li>
            <li>3 / 1 = 3</li>
            <li>3 \ 1 = 3</li> 
            <li>3 % 1 = 0</li>
            <li>3 ^ 1 = 3</li>
        
            <li>3 + 2 = 5</li>
            <li>3 - 2 = 1</li>
            <li>3 * 2 = 6</li>
            <li>3 / 2 = 1.5</li>
            <li>3 \ 2 = 1</li> 
            <li>3 % 2 = 1</li>
            <li>3 ^ 2 = 9</li>
        
            <li>3 + 3 = 6</li>
            <li>3 - 3 = 0</li>
            <li>3 * 3 = 9</li>
            <li>3 / 3 = 1</li>
            <li>3 \ 3 = 1</li> 
            <li>3 % 3 = 0</li>
            <li>3 ^ 3 = 27</li>
        
            <li>3 + 4 = 7</li>
            <li>3 - 4 = -1</li>
            <li>3 * 4 = 12</li>
            <li>3 / 4 = 0.75</li>
            <li>3 \ 4 = 0</li> 
            <li>3 % 4 = 3</li>
            <li>3 ^ 4 = 81</li>
        
            <li>3 + 5 = 8</li>
            <li>3 - 5 = -2</li>
            <li>3 * 5 = 15</li>
            <li>3 / 5 = 0.6</li>
            <li>3 \ 5 = 0</li> 
            <li>3 % 5 = 3</li>
            <li>3 ^ 5 = 243</li>
        
            <li>3 + 6 = 9</li>
            <li>3 - 6 = -3</li>
            <li>3 * 6 = 18</li>
            <li>3 / 6 = 0.5</li>
            <li>3 \ 6 = 0</li> 
            <li>3 % 6 = 3</li>
            <li>3 ^ 6 = 729</li>
        
            <li>3 + 7 = 10</li>
            <li>3 - 7 = -4</li>
            <li>3 * 7 = 21</li>
            <li>3 / 7 = 0.42857142857142855</li>
            <li>3 \ 7 = 0</li> 
            <li>3 % 7 = 3</li>
            <li>3 ^ 7 = 2187</li>
        
            <li>3 + 8 = 11</li>
            <li>3 - 8 = -5</li>
            <li>3 * 8 = 24</li>
            <li>3 / 8 = 0.375</li>
            <li>3 \ 8 = 0</li> 
            <li>3 % 8 = 3</li>
            <li>3 ^ 8 = 6561</li>
        
            <li>3 + 9 = 12</li>
            <li>3 - 9 = -6</li>
            <li>3 * 9 = 27</li>
            <li>3 / 9 = 0.3333333333333333</li>
            <li>3 \ 9 = 0</li> 
            <li>3 % 9 = 3</li>
            <li>3 ^ 9 = 19683</li>
        
            <li>3 + 0 = 3</li>
            <li>3 - 0 = 3</li>
            <li>3 * 0 = 0</li>
            <li>3 / 0 = Error: divide by zero</li>
            <li>3 \ 0 = Error: divide by zero</li> 
            <li>3 % 0 = Error: divide by zero</li>
            <li>3 ^ 0 = 1</li>
        
    
        
            <li>4 + 1 = 5</li>
            <li>4 - 1 = 3</li>
            <li>4 * 1 = 4</li>
            <li>4 / 1 = 4</li>
            <li>4 \ 1 = 4</li> 
            <li>4 % 1 = 0</li>
            <li>4 ^ 1 = 4</li>
        
            <li>4 + 2 = 6</li>
            <li>4 - 2 = 2</li>
            <li>4 * 2 = 8</li>
            <li>4 / 2 = 2</li>
            <li>4 \ 2 = 2</li> 
            <li>4 % 2 = 0</li>
            <li>4 ^ 2 = 16</li>
        
            <li>4 + 3 = 7</li>
            <li>4 - 3 = 1</li>
            <li>4 * 3 = 12</li>
            <li>4 / 3 = 1.3333333333333333</li>
            <li>4 \ 3 = 1</li> 
            <li>4 % 3 = 1</li>
            <li>4 ^ 3 = 64</li>
        
            <li>4 + 4 = 8</li>
            <li>4 - 4 = 0</li>
            <li>4 * 4 = 16</li>
            <li>4 / 4 = 1</li>
            <li>4 \ 4 = 1</li> 
            <li>4 % 4 = 0</li>
            <li>4 ^ 4 = 256</li>
        
            <li>4 + 5 = 9</li>
            <li>4 - 5 = -1</li>
            <li>4 * 5 = 20</li>
            <li>4 / 5 = 0.8</li>
            <li>4 \ 5 = 0</li> 
            <li>4 % 5 = 4</li>
            <li>4 ^ 5 = 1024</li>
        
            <li>4 + 6 = 10</li>
            <li>4 - 6 = -2</li>
            <li>4 * 6 = 24</li>
            <li>4 / 6 = 0.6666666666666666</li>
            <li>4 \ 6 = 0</li> 
            <li>4 % 6 = 4</li>
            <li>4 ^ 6 = 4096</li>
        
            <li>4 + 7 = 11</li>
            <li>4 - 7 = -3</li>
            <li>4 * 7 = 28</li>
            <li>4 / 7 = 0.5714285714285714</li>
            <li>4 \ 7 = 0</li> 
            <li>4 % 7 = 4</li>
            <li>4 ^ 7 = 16384</li>
        
            <li>4 + 8 = 12</li>
            <li>4 - 8 = -4</li>
            <li>4 * 8 = 32</li>
            <li>4 / 8 = 0.5</li>
            <li>4 \ 8 = 0</li> 
            <li>4 % 8 = 4</li>
            <li>4 ^ 8 = 65536</li>
        
            <li>4 + 9 = 13</li>
            <li>4 - 9 = -5</li>
            <li>4 * 9 = 36</li>
            <li>4 / 9 = 0.4444444444444444</li>
            <li>4 \ 9 = 0</li> 
            <li>4 % 9 = 4</li>
            <li>4 ^ 9 = 262144</li>
        
            <li>4 + 0 = 4</li>
            <li>4 - 0 = 4</li>
            <li>4 * 0 = 0</li>
            <li>4 / 0 = Error: divide by zero</li>
            <li>4 \ 0 = Error: divide by zero</li> 
            <li>4 % 0 = Error: divide by zero</li>
            <li>4 ^ 0 = 1</li>
        
    
        
            <li>5 + 1 = 6</li>
            <li>5 - 1 = 4</li>
            <li>5 * 1 = 5</li>
            <li>5 / 1 = 5</li>
            <li>5 \ 1 = 5</li> 
            <li>5 % 1 = 0</li>
            <li>5 ^ 1 = 5</li>
        
            <li>5 + 2 = 7</li>
            <li>5 - 2 = 3</li>
            <li>5 * 2 = 10</li>
            <li>5 / 2 = 2.5</li>
            <li>5 \ 2 = 2</li> 
            <li>5 % 2 = 1</li>
            <li>5 ^ 2 = 25</li>
        
            <li>5 + 3 = 8</li>
            <li>5 - 3 = 2</li>
            <li>5 * 3 = 15</li>
            <li>5 / 3 = 1.6666666666666667</li>
            <li>5 \ 3 = 1</li> 
            <li>5 % 3 = 2</li>
            <li>5 ^ 3 = 125</li>
        
            <li>5 + 4 = 9</li>
            <li>5 - 4 = 1</li>
            <li>5 * 4 = 20</li>
            <li>5 / 4 = 1.25</li>
            <li>5 \ 4 = 1</li> 
            <li>5 % 4 = 1</li>
            <li>5 ^ 4 = 625</li>
        
            <li>5 + 5 = 10</li>
            <li>5 - 5 = 0</li>
            <li>5 * 5 = 25</li>
            <li>5 / 5 = 1</li>
            <li>5 \ 5 = 1</li> 
            <li>5 % 5 = 0</li>
            <li>5 ^ 5 = 3125</li>
        
            <li>5 + 6 = 11</li>
            <li>5 - 6 = -1</li>
            <li>5 * 6 = 30</li>
            <li>5 / 6 = 0.8333333333333334</li>
            <li>5 \ 6 = 0</li> 
            <li>5 % 6 = 5</li>
            <li>5 ^ 6 = 15625</li>
        
            <li>5 + 7 = 12</li>
            <li>5 - 7 = -2</li>
            <li>5 * 7 = 35</li>
            <li>5 / 7 = 0.7142857142857143</li>
            <li>5 \ 7 = 0</li> 
            <li>5 % 7 = 5</li>
            <li>5 ^ 7 = 78125</li>
        
            <li>5 + 8 = 13</li>
            <li>5 - 8 = -3</li>
            <li>5 * 8 = 40</li>
            <li>5 / 8 = 0.625</li>
            <li>5 \ 8 = 0</li> 
            <li>5 % 8 = 5</li>
            <li>5 ^ 8 = 390625</li>
        
            <li>5 + 9 = 14</li>
            <li>5 - 9 = -4</li>
            <li>5 * 9 = 45</li>
            <li>5 / 9 = 0.5555555555555556</li>
            <li>5 \ 9 = 0</li> 
            <li>5 % 9 = 5</li>
            <li>5 ^ 9 = 1953125</li>
        
            <li>5 + 0 = 5</li>
            <li>5 - 0 = 5</li>
            <li>5 * 0 = 0</li>
            <li>5 / 0 = Error: divide by zero</li>
            <li>5 \ 0 = Error: divide by zero</li> 
            <li>5 % 0 = Error: divide by zero</li>
            <li>5 ^ 0 = 1</li>
        
    
        
            <li>6 + 1 = 7</li>
            <li>6 - 1 = 5</li>
            <li>6 * 1 = 6</li>
            <li>6 / 1 = 6</li>
            <li>6 \ 1 = 6</li> 
            <li>6 % 1 = 0</li>
            <li>6 ^ 1 = 6</li>
        
            <li>6 + 2 = 8</li>
            <li>6 - 2 = 4</li>
            <li>6 * 2 = 12</li>
            <li>6 / 2 = 3</li>
            <li>6 \ 2 = 3</li> 
            <li>6 % 2 = 0</li>
            <li>6 ^ 2 = 36</li>
        
            <li>6 + 3 = 9</li>
            <li>6 - 3 = 3</li>
            <li>6 * 3 = 18</li>
            <li>6 / 3 = 2</li>
            <li>6 \ 3 = 2</li> 
            <li>6 % 3 = 0</li>
            <li>6 ^ 3 = 216</li>
        
            <li>6 + 4 = 10</li>
            <li>6 - 4 = 2</li>
            <li>6 * 4 = 24</li>
            <li>6 / 4 = 1.5</li>
            <li>6 \ 4 = 1</li> 
            <li>6 % 4 = 2</li>
            <li>6 ^ 4 = 1296</li>
        
            <li>6 + 5 = 11</li>
            <li>6 - 5 = 1</li>
            <li>6 * 5 = 30</li>
            <li>6 / 5 = 1.2</li>
            <li>6 \ 5 = 1</li> 
            <li>6 % 5 = 1</li>
            <li>6 ^ 5 = 7776</li>
        
            <li>6 + 6 = 12</li>
            <li>6 - 6 = 0</li>
            <li>6 * 6 = 36</li>
            <li>6 / 6 = 1</li>
            <li>6 \ 6 = 1</li> 
            <li>6 % 6 = 0</li>
            <li>6 ^ 6 = 46656</li>
        
            <li>6 + 7 = 13</li>
            <li>6 - 7 = -1</li>
            <li>6 * 7 = 42</li>
            <li>6 / 7 = 0.8571428571428571</li>
            <li>6 \ 7 = 0</li> 
            <li>6 % 7 = 6</li>
            <li>6 ^ 7 = 279936</li>
        
            <li>6 + 8 = 14</li>
            <li>6 - 8 = -2</li>
            <li>6 * 8 = 48</li>
            <li>6 / 8 = 0.75</li>
            <li>6 \ 8 = 0</li> 
            <li>6 % 8 = 6</li>
            <li>6 ^ 8 = 1679616</li>
        
            <li>6 + 9 = 15</li>
            <li>6 - 9 = -3</li>
            <li>6 * 9 = 54</li>
            <li>6 / 9 = 0.6666666666666666</li>
            <li>6 \ 9 = 0</li> 
            <li>6 % 9 = 6</li>
            <li>6 ^ 9 = 10077696</li>
        
            <li>6 + 0 = 6</li>
            <li>6 - 0 = 6</li>
            <li>6 * 0 = 0</li>
            <li>6 / 0 = Error: divide by zero</li>
            <li>6 \ 0 = Error: divide by zero</li> 
            <li>6 % 0 = Error: divide by zero</li>
            <li>6 ^ 0 = 1</li>
        
    
        
            <li>7 + 1 = 8</li>
            <li>7 - 1 = 6</li>
            <li>7 * 1 = 7</li>
            <li>7 / 1 = 7</li>
            <li>7 \ 1 = 7</li> 
            <li>7 % 1 = 0</li>
            <li>7 ^ 1 = 7</li>
        
            <li>7 + 2 = 9</li>
            <li>7 - 2 = 5</li>
            <li>7 * 2 = 14</li>
            <li>7 / 2 = 3.5</li>
            <li>7 \ 2 = 3</li> 
            <li>7 % 2 = 1</li>
            <li>7 ^ 2 = 49</li>
        
            <li>7 + 3 = 10</li>
            <li>7 - 3 = 4</li>
            <li>7 * 3 = 21</li>
            <li>7 / 3 = 2.3333333333333335</li>
            <li>7 \ 3 = 2</li> 
            <li>7 % 3 = 1</li>
            <li>7 ^ 3 = 343</li>
        
            <li>7 + 4 = 11</li>
            <li>7 - 4 = 3</li>
            <li>7 * 4 = 28</li>
            <li>7 / 4 = 1.75</li>
            <li>7 \ 4 = 1</li> 
            <li>7 % 4 = 3</li>
            <li>7 ^ 4 = 2401</li>
        
            <li>7 + 5 = 12</li>
            <li>7 - 5 = 2</li>
            <li>7 * 5 = 35</li>
            <li>7 / 5 = 1.4</li>
            <li>7 \ 5 = 1</li> 
            <li>7 % 5 = 2</li>
            <li>7 ^ 5 = 16807</li>
        
            <li>7 + 6 = 13</li>
            <li>7 - 6 = 1</li>
            <li>7 * 6 = 42</li>
            <li>7 / 6 = 1.1666666666666667</li>
            <li>7 \ 6 = 1</li> 
            <li>7 % 6 = 1</li>
            <li>7 ^ 6 = 117649</li>
        
            <li>7 + 7 = 14</li>
            <li>7 - 7 = 0</li>
            <li>7 * 7 = 49</li>
            <li>7 / 7 = 1</li>
            <li>7 \ 7 = 1</li> 
            <li>7 % 7 = 0</li>
            <li>7 ^ 7 = 823543</li>
        
            <li>7 + 8 = 15</li>
            <li>7 - 8 = -1</li>
            <li>7 * 8 = 56</li>
            <li>7 / 8 = 0.875</li>
            <li>7 \ 8 = 0</li> 
            <li>7 % 8 = 7</li>
            <li>7 ^ 8 = 5764801</li>
        
            <li>7 + 9 = 16</li>
            <li>7 - 9 = -2</li>
            <li>7 * 9 = 63</li>
            <li>7 / 9 = 0.7777777777777778</li>
            <li>7 \ 9 = 0</li> 
            <li>7 % 9 = 7</li>
            <li>7 ^ 9 = 40353607</li>
        
            <li>7 + 0 = 7</li>
            <li>7 - 0 = 7</li>
            <li>7 * 0 = 0</li>
            <li>7 / 0 = Error: divide by zero</li>
            <li>7 \ 0 = Error: divide by zero</li> 
            <li>7 % 0 = Error: divide by zero</li>
            <li>7 ^ 0 = 1</li>
        
    
        
            <li>8 + 1 = 9</li>
            <li>8 - 1 = 7</li>
            <li>8 * 1 = 8</li>
            <li>8 / 1 = 8</li>
            <li>8 \ 1 = 8</li> 
            <li>8 % 1 = 0</li>
            <li>8 ^ 1 = 8</li>
        
            <li>8 + 2 = 10</li>
            <li>8 - 2 = 6</li>
            <li>8 * 2 = 16</li>
            <li>8 / 2 = 4</li>
            <li>8 \ 2 = 4</li> 
            <li>8 % 2 = 0</li>
            <li>8 ^ 2 = 64</li>
        
            <li>8 + 3 = 11</li>
            <li>8 - 3 = 5</li>
            <li>8 * 3 = 24</li>
            <li>8 / 3 = 2.6666666666666665</li>
            <li>8 \ 3 = 2</li> 
            <li>8 % 3 = 2</li>
            <li>8 ^ 3 = 512</li>
        
            <li>8 + 4 = 12</li>
            <li>8 - 4 = 4</li>
            <li>8 * 4 = 32</li>
            <li>8 / 4 = 2</li>
            <li>8 \ 4 = 2</li> 
            <li>8 % 4 = 0</li>
            <li>8 ^ 4 = 4096</li>
        
            <li>8 + 5 = 13</li>
            <li>8 - 5 = 3</li>
            <li>8 * 5 = 40</li>
            <li>8 / 5 = 1.6</li>
            <li>8 \ 5 = 1</li> 
            <li>8 % 5 = 3</li>
            <li>8 ^ 5 = 32768</li>
        
            <li>8 + 6 = 14</li>
            <li>8 - 6 = 2</li>
            <li>8 * 6 = 48</li>
            <li>8 / 6 = 1.3333333333333333</li>
            <li>8 \ 6 = 1</li> 
            <li>8 % 6 = 2</li>
            <li>8 ^ 6 = 262144</li>
        
            <li>8 + 7 = 15</li>
            <li>8 - 7 = 1</li>
            <li>8 * 7 = 56</li>
            <li>8 / 7 = 1.1428571428571428</li>
            <li>8 \ 7 = 1</li> 
            <li>8 % 7 = 1</li>
            <li>8 ^ 7 = 2097152</li>
        
            <li>8 + 8 = 16</li>
            <li>8 - 8 = 0</li>
            <li>8 * 8 = 64</li>
            <li>8 / 8 = 1</li>
            <li>8 \ 8 = 1</li> 
            <li>8 % 8 = 0</li>
            <li>8 ^ 8 = 16777216</li>
        
            <li>8 + 9 = 17</li>
            <li>8 - 9 = -1</li>
            <li>8 * 9 = 72</li>
            <li>8 / 9 = 0.8888888888888888</li>
            <li>8 \ 9 = 0</li> 
            <li>8 % 9 = 8</li>
            <li>8 ^ 9 = 134217728</li>
        
            <li>8 + 0 = 8</li>
            <li>8 - 0 = 8</li>
            <li>8 * 0 = 0</li>
            <li>8 / 0 = Error: divide by zero</li>
            <li>8 \ 0 = Error: divide by zero</li> 
            <li>8 % 0 = Error: divide by zero</li>
            <li>8 ^ 0 = 1</li>
        
    
        
            <li>9 + 1 = 10</li>
            <li>9 - 1 = 8</li>
            <li>9 * 1 = 9</li>
            <li>9 / 1 = 9</li>
            <li>9 \ 1 = 9</li> 
            <li>9 % 1 = 0</li>
            <li>9 ^ 1 = 9</li>
        
            <li>9 + 2 = 11</li>
            <li>9 - 2 = 7</li>
            <li>9 * 2 = 18</li>
            <li>9 / 2 = 4.5</li>
            <li>9 \ 2 = 4</li> 
            <li>9 % 2 = 1</li>
            <li>9 ^ 2 = 81</li>
        
            <li>9 + 3 = 12</li>
            <li>9 - 3 = 6</li>
            <li>9 * 3 = 27</li>
            <li>9 / 3 = 3</li>
            <li>9 \ 3 = 3</li> 
            <li>9 % 3 = 0</li>
            <li>9 ^ 3 = 729</li>
        
            <li>9 + 4 = 13</li>
            <li>9 - 4 = 5</li>
            <li>9 * 4 = 36</li>
            <li>9 / 4 = 2.25</li>
            <li>9 \ 4 = 2</li> 
            <li>9 % 4 = 1</li>
            <li>9 ^ 4 = 6561</li>
        
            <li>9 + 5 = 14</li>
            <li>9 - 5 = 4</li>
            <li>9 * 5 = 45</li>
            <li>9 / 5 = 1.8</li>
            <li>9 \ 5 = 1</li> 
            <li>9 % 5 = 4</li>
            <li>9 ^ 5 = 59049</li>
        
            <li>9 + 6 = 15</li>
            <li>9 - 6 = 3</li>
            <li>9 * 6 = 54</li>
            <li>9 / 6 = 1.5</li>
            <li>9 \ 6 = 1</li> 
            <li>9 % 6 = 3</li>
            <li>9 ^ 6 = 531441</li>
        
            <li>9 + 7 = 16</li>
            <li>9 - 7 = 2</li>
            <li>9 * 7 = 63</li>
            <li>9 / 7 = 1.2857142857142858</li>
            <li>9 \ 7 = 1</li> 
            <li>9 % 7 = 2</li>
            <li>9 ^ 7 = 4782969</li>
        
            <li>9 + 8 = 17</li>
            <li>9 - 8 = 1</li>
            <li>9 * 8 = 72</li>
            <li>9 / 8 = 1.125</li>
            <li>9 \ 8 = 1</li> 
            <li>9 % 8 = 1</li>
            <li>9 ^ 8 = 43046721</li>
        
            <li>9 + 9 = 18</li>
            <li>9 - 9 = 0</li>
            <li>9 * 9 = 81</li>
            <li>9 / 9 = 1</li>
            <li>9 \ 9 = 1</li> 
            <li>9 % 9 = 0</li>
            <li>9 ^ 9 = 387420489</li>
        
            <li>9 + 0 = 9</li>
            <li>9 - 0 = 9</li>
            <li>9 * 0 = 0</li>
            <li>9 / 0 = Error: divide by zero</li>
            <li>9 \ 0 = Error: divide by zero</li> 
            <li>9 % 0 = Error: divide by zero</li>
            <li>9 ^ 0 = 1</li>
        
    
        
            <li>0 + 1 = 1</li>
            <li>0 - 1 = -1</li>
            <li>0 * 1 = 0</li>
            <li>0 / 1 = 0</li>
            <li>0 \ 1 = 0</li> 
            <li>0 % 1 = 0</li>
            <li>0 ^ 1 = 0</li>
        
            <li>0 + 2 = 2</li>
            <li>0 - 2 = -2</li>
            <li>0 * 2 = 0</li>
            <li>0 / 2 = 0</li>
            <li>0 \ 2 = 0</li> 
            <li>0 % 2 = 0</li>
            <li>0 ^ 2 = 0</li>
        
            <li>0 + 3 = 3</li>
            <li>0 - 3 = -3</li>
            <li>0 * 3 = 0</li>
            <li>0 / 3 = 0</li>
            <li>0 \ 3 = 0</li> 
            <li>0 % 3 = 0</li>
            <li>0 ^ 3 = 0</li>
        
            <li>0 + 4 = 4</li>
            <li>0 - 4 = -4</li>
            <li>0 * 4 = 0</li>
            <li>0 / 4 = 0</li>
            <li>0 \ 4 = 0</li> 
            <li>0 % 4 = 0</li>
            <li>0 ^ 4 = 0</li>
        
            <li>0 + 5 = 5</li>
            <li>0 - 5 = -5</li>
            <li>0 * 5 = 0</li>
            <li>0 / 5 = 0</li>
            <li>0 \ 5 = 0</li> 
            <li>0 % 5 = 0</li>
            <li>0 ^ 5 = 0</li>
        
            <li>0 + 6 = 6</li>
            <li>0 - 6 = -6</li>
            <li>0 * 6 = 0</li>
            <li>0 / 6 = 0</li>
            <li>0 \ 6 = 0</li> 
            <li>0 % 6 = 0</li>
            <li>0 ^ 6 = 0</li>
        
            <li>0 + 7 = 7</li>
            <li>0 - 7 = -7</li>
            <li>0 * 7 = 0</li>
            <li>0 / 7 = 0</li>
            <li>0 \ 7 = 0</li> 
            <li>0 % 7 = 0</li>
            <li>0 ^ 7 = 0</li>
        
            <li>0 + 8 = 8</li>
            <li>0 - 8 = -8</li>
            <li>0 * 8 = 0</li>
            <li>0 / 8 = 0</li>
            <li>0 \ 8 = 0</li> 
            <li>0 % 8 = 0</li>
            <li>0 ^ 8 = 0</li>
        
            <li>0 + 9 = 9</li>
            <li>0 - 9 = -9</li>
            <li>0 * 9 = 0</li>
            <li>0 / 9 = 0</li>
            <li>0 \ 9 = 0</li> 
            <li>0 % 9 = 0</li>
            <li>0 ^ 9 = 0</li>
        
            <li>0 + 0 = 0</li>
            <li>0 - 0 = 0</li>
            <li>0 * 0 = 0</li>
            <li>0 / 0 = Error: divide by zero</li>
            <li>0 \ 0 = Error: divide by zero</li> 
            <li>0 % 0 = Error: divide by zero</li>
            <li>0 ^ 0 = 1</li>
        
    
    </ul>

    <h2>Математические операции (краткая запись с присвоением)</h2>
    <ul>
    
        <li>
            $myVar = 1<br> 
            $myVar += 6<br> 
            $myVar = 7
        </li>
        <li>
            $myVar = 1<br> 
            $myVar -= 6<br> 
            $myVar = -5
        </li>
        <li>
            $myVar = 1<br> 
            $myVar *= 6<br> 
            $myVar = 6
        </li>
        <li>
            $myVar = 1<br> 
            $myVar /= 6<br> 
            $myVar = 0.16666666666666666
        </li>
        <li>
            $myVar = 1<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 1<br> 
            $myVar %= 6<br> 
            $myVar = 1
        </li>
        <li>
            $myVar = 1<br> 
            $myVar ^= 6<br> 
            $myVar = 1
        </li>
    
        <li>
            $myVar = 2<br> 
            $myVar += 6<br> 
            $myVar = 8
        </li>
        <li>
            $myVar = 2<br> 
            $myVar -= 6<br> 
            $myVar = -4
        </li>
        <li>
            $myVar = 2<br> 
            $myVar *= 6<br> 
            $myVar = 12
        </li>
        <li>
            $myVar = 2<br> 
            $myVar /= 6<br> 
            $myVar = 0.3333333333333333
        </li>
        <li>
            $myVar = 2<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 2<br> 
            $myVar %= 6<br> 
            $myVar = 2
        </li>
        <li>
            $myVar = 2<br> 
            $myVar ^= 6<br> 
            $myVar = 64
        </li>
    
        <li>
            $myVar = 3<br> 
            $myVar += 6<br> 
            $myVar = 9
        </li>
        <li>
            $myVar = 3<br> 
            $myVar -= 6<br> 
            $myVar = -3
        </li>
        <li>
            $myVar = 3<br> 
            $myVar *= 6<br> 
            $myVar = 18
        </li>
        <li>
            $myVar = 3<br> 
            $myVar /= 6<br> 
            $myVar = 0.5
        </li>
        <li>
            $myVar = 3<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 3<br> 
            $myVar %= 6<br> 
            $myVar = 3
        </li>
        <li>
            $myVar = 3<br> 
            $myVar ^= 6<br> 
            $myVar = 729
        </li>
    
        <li>
            $myVar = 4<br> 
            $myVar += 6<br> 
            $myVar = 10
        </li>
        <li>
            $myVar = 4<br> 
            $myVar -= 6<br> 
            $myVar = -2
        </li>
        <li>
            $myVar = 4<br> 
            $myVar *= 6<br> 
            $myVar = 24
        </li>
        <li>
            $myVar = 4<br> 
            $myVar /= 6<br> 
            $myVar = 0.6666666666666666
        </li>
        <li>
            $myVar = 4<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 4<br> 
            $myVar %= 6<br> 
            $myVar = 4
        </li>
        <li>
            $myVar = 4<br> 
            $myVar ^= 6<br> 
            $myVar = 4096
        </li>
    
        <li>
            $myVar = 5<br> 
            $myVar += 6<br> 
            $myVar = 11
        </li>
        <li>
            $myVar = 5<br> 
            $myVar -= 6<br> 
            $myVar = -1
        </li>
        <li>
            $myVar = 5<br> 
            $myVar *= 6<br> 
            $myVar = 30
        </li>
        <li>
            $myVar = 5<br> 
            $myVar /= 6<br> 
            $myVar = 0.8333333333333334
        </li>
        <li>
            $myVar = 5<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 5<br> 
            $myVar %= 6<br> 
            $myVar = 5
        </li>
        <li>
            $myVar = 5<br> 
            $myVar ^= 6<br> 
            $myVar = 15625
        </li>
    
    </ul>

    <h2>Конкатенация строк</h2>
    <ul>
    
        
            <li>Первый Первый</li>
        
            <li>Первый Второй</li>
        
            <li>Первый Третий</li>
        
            <li>Первый 048465450</li>
        
    
        
            <li>Второй Первый</li>
        
            <li>Второй Второй</li>
        
            <li>Второй Третий</li>
        
            <li>Второй 048465450</li>
        
    
        
            <li>Третий Первый</li>
        
            <li>Третий Второй</li>
        
            <li>Третий Третий</li>
        
            <li>Третий 048465450</li>
        
    
        
            <li>048465450 Первый</li>
        
            <li>048465450 Второй</li>
        
            <li>048465450 Третий</li>
        
            <li>048465450 048465450</li>
        
    
    </ul>

    <h2>Доступ к полям структур</h2>
    <ul>
    
        <li>
            <b>0:</b>
<pre>
struct {
    ID: 1,
    Name: my name,
    Value: my value,
    Function: Мама мыла раму 1 дней
}
</pre>
        </li>
    
        <li>
            <b>1:</b>
<pre>
struct {
    ID: 2,
    Name: my name 2,
    Value: my value 2,
    Function: Мама мыла раму 2 дней
}
</pre>
        </li>
    
    </ul>

    <h2>Вызов пользовательской функции</h2>
    <ul>
        <li>Присвоим результат выполнения функции переменной $myData </li>
        <li>В переменной $myData хранится значение 1945-05-09 00:43</li>
        <li>Выведем результат выполнения функции сразу в шаблон 1945-05-09 00:43 или так 1945-05-09 00:43</li>
    </ul>

    <h2>Комментарии</h2>
    
    

    <h2>Подключение шаблонов</h2>
    <div>layouts/include.tpl Подключен</div>

    <h1>Вывод кода шаблона без обработки</h2>
<pre>
    
    @foreach($numbers as $number1)
        @foreach($numbers as $number2)
            &lt;li&gt;{{$number1}} + {{$number2}} = {{$number1 + $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} - {{$number2}} = {{$number1 - $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} * {{$number2}} = {{$number1 * $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} / {{$number2}} = {{$number1 / $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} % {{$number2}} = {{$number1 % $number2}}&lt;/li&gt;
        @endforeach
    @endforeach
    
</pre>
</body>
</html>`

var result2 = `<!DOCTYPE html>
<html lang="ru">
<head>
    <title>Тестовая страница 2</title>
</head>
<body>
    <h1>Тестовая страница 2<h1/>

    <h2>Вывод переменных в шаблон</h2>
    <ul>
        <li>&lt;b&gt;Экранированный вывод&lt;/b&gt;</li>
        <li><b>НЕ экранированный вывод</b></li>
    </ul>

    <h2>Объявление переменных и и прочие операции, без вывода результата в шаблон</h2>
    <ul>
        <li>Тут объявляю переменную и присваиваю ей результат математической операции, в шаблоне не вывожу </li>
        <li>Тут выведу значение сохранённое в переменной 50</li>
    </ul>

    <h2>Цикл for</h2>
    <ul>
    
        <li>
            
                В переменной $strings, под индексом 0 содержится значение Первый 2
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 1 содержится значение Второй 2
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 2 содержится значение Третий 2
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 3 содержится значение 07869786
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 4 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 5 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 6 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 7 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 8 нет значения
            
        </li>
    
        <li>
            
                В переменной $strings, под индексом 9 нет значения
            
        </li>
    
    </ul>

    <h2>Цикл foreach</h2>
    <ul>
    
        <li>$strings[0] == Первый 2</li>
    
        <li>$strings[1] == Второй 2</li>
    
        <li>$strings[2] == Третий 2</li>
    
        <li>$strings[3] == 07869786</li>
    
    </ul>

    <h2>Условия</h2>
    <ul>
    
        
            
            
                <li>0 <= 0</li>
            
        
            
                <li>9 > 8 && 0 < 8 || 9 < 1 && 0 > 1</li>
            
            
                <li>0 < 9</li>
            
        
            
            
                <li>0 < 8</li>
            
        
            
            
                <li>0 < 7</li>
            
        
            
            
                <li>0 < 6</li>
            
        
            
            
                <li>0 < 5</li>
            
        
            
            
                <li>0 < 4</li>
            
        
            
            
                <li>0 < 3</li>
            
        
            
            
                <li>0 < 2</li>
            
        
            
            
                <li>0 < 1</li>
            
        
            
            
                <li>0 <= 0</li>
            
        
    
        
            
                <li>0 > 8 && 9 < 8 || 0 < 1 && 9 > 1</li>
            
            
                <li>9 >= 0</li>
            
        
            
            
                <li>9 <= 9</li>
            
        
            
            
                <li>9 >= 8</li>
            
        
            
            
                <li>9 >= 7</li>
            
        
            
            
                <li>9 >= 6</li>
            
        
            
            
                <li>9 >= 5</li>
            
        
            
            
                <li>9 >= 4</li>
            
        
            
            
                <li>9 >= 3</li>
            
        
            
            
                <li>9 >= 2</li>
            
        
            
            
                <li>9 >= 1</li>
            
        
            
                <li>0 > 8 && 9 < 8 || 0 < 1 && 9 > 1</li>
            
            
                <li>9 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 8 < 8 || 0 < 1 && 8 > 1</li>
            
            
                <li>8 >= 0</li>
            
        
            
            
                <li>8 < 9</li>
            
        
            
            
                <li>8 <= 8</li>
            
        
            
            
                <li>8 >= 7</li>
            
        
            
            
                <li>8 >= 6</li>
            
        
            
            
                <li>8 >= 5</li>
            
        
            
            
                <li>8 >= 4</li>
            
        
            
            
                <li>8 >= 3</li>
            
        
            
            
                <li>8 >= 2</li>
            
        
            
            
                <li>8 >= 1</li>
            
        
            
                <li>0 > 8 && 8 < 8 || 0 < 1 && 8 > 1</li>
            
            
                <li>8 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 7 < 8 || 0 < 1 && 7 > 1</li>
            
            
                <li>7 >= 0</li>
            
        
            
                <li>9 > 8 && 7 < 8 || 9 < 1 && 7 > 1</li>
            
            
                <li>7 < 9</li>
            
        
            
            
                <li>7 < 8</li>
            
        
            
            
                <li>7 <= 7</li>
            
        
            
            
                <li>7 >= 6</li>
            
        
            
            
                <li>7 >= 5</li>
            
        
            
            
                <li>7 >= 4</li>
            
        
            
            
                <li>7 >= 3</li>
            
        
            
            
                <li>7 >= 2</li>
            
        
            
            
                <li>7 >= 1</li>
            
        
            
                <li>0 > 8 && 7 < 8 || 0 < 1 && 7 > 1</li>
            
            
                <li>7 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 6 < 8 || 0 < 1 && 6 > 1</li>
            
            
                <li>6 >= 0</li>
            
        
            
                <li>9 > 8 && 6 < 8 || 9 < 1 && 6 > 1</li>
            
            
                <li>6 < 9</li>
            
        
            
            
                <li>6 < 8</li>
            
        
            
            
                <li>6 < 7</li>
            
        
            
            
                <li>6 <= 6</li>
            
        
            
            
                <li>6 >= 5</li>
            
        
            
            
                <li>6 >= 4</li>
            
        
            
            
                <li>6 >= 3</li>
            
        
            
            
                <li>6 >= 2</li>
            
        
            
            
                <li>6 >= 1</li>
            
        
            
                <li>0 > 8 && 6 < 8 || 0 < 1 && 6 > 1</li>
            
            
                <li>6 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 5 < 8 || 0 < 1 && 5 > 1</li>
            
            
                <li>5 >= 0</li>
            
        
            
                <li>9 > 8 && 5 < 8 || 9 < 1 && 5 > 1</li>
            
            
                <li>5 < 9</li>
            
        
            
            
                <li>5 < 8</li>
            
        
            
            
                <li>5 < 7</li>
            
        
            
            
                <li>5 < 6</li>
            
        
            
            
                <li>5 <= 5</li>
            
        
            
            
                <li>5 >= 4</li>
            
        
            
            
                <li>5 >= 3</li>
            
        
            
            
                <li>5 >= 2</li>
            
        
            
            
                <li>5 >= 1</li>
            
        
            
                <li>0 > 8 && 5 < 8 || 0 < 1 && 5 > 1</li>
            
            
                <li>5 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 4 < 8 || 0 < 1 && 4 > 1</li>
            
            
                <li>4 >= 0</li>
            
        
            
                <li>9 > 8 && 4 < 8 || 9 < 1 && 4 > 1</li>
            
            
                <li>4 < 9</li>
            
        
            
            
                <li>4 < 8</li>
            
        
            
            
                <li>4 < 7</li>
            
        
            
            
                <li>4 < 6</li>
            
        
            
            
                <li>4 < 5</li>
            
        
            
            
                <li>4 <= 4</li>
            
        
            
            
                <li>4 >= 3</li>
            
        
            
            
                <li>4 >= 2</li>
            
        
            
            
                <li>4 >= 1</li>
            
        
            
                <li>0 > 8 && 4 < 8 || 0 < 1 && 4 > 1</li>
            
            
                <li>4 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 3 < 8 || 0 < 1 && 3 > 1</li>
            
            
                <li>3 >= 0</li>
            
        
            
                <li>9 > 8 && 3 < 8 || 9 < 1 && 3 > 1</li>
            
            
                <li>3 < 9</li>
            
        
            
            
                <li>3 < 8</li>
            
        
            
            
                <li>3 < 7</li>
            
        
            
            
                <li>3 < 6</li>
            
        
            
            
                <li>3 < 5</li>
            
        
            
            
                <li>3 < 4</li>
            
        
            
            
                <li>3 <= 3</li>
            
        
            
            
                <li>3 >= 2</li>
            
        
            
            
                <li>3 >= 1</li>
            
        
            
                <li>0 > 8 && 3 < 8 || 0 < 1 && 3 > 1</li>
            
            
                <li>3 >= 0</li>
            
        
    
        
            
                <li>0 > 8 && 2 < 8 || 0 < 1 && 2 > 1</li>
            
            
                <li>2 >= 0</li>
            
        
            
                <li>9 > 8 && 2 < 8 || 9 < 1 && 2 > 1</li>
            
            
                <li>2 < 9</li>
            
        
            
            
                <li>2 < 8</li>
            
        
            
            
                <li>2 < 7</li>
            
        
            
            
                <li>2 < 6</li>
            
        
            
            
                <li>2 < 5</li>
            
        
            
            
                <li>2 < 4</li>
            
        
            
            
                <li>2 < 3</li>
            
        
            
            
                <li>2 <= 2</li>
            
        
            
            
                <li>2 >= 1</li>
            
        
            
                <li>0 > 8 && 2 < 8 || 0 < 1 && 2 > 1</li>
            
            
                <li>2 >= 0</li>
            
        
    
        
            
            
                <li>1 >= 0</li>
            
        
            
                <li>9 > 8 && 1 < 8 || 9 < 1 && 1 > 1</li>
            
            
                <li>1 < 9</li>
            
        
            
            
                <li>1 < 8</li>
            
        
            
            
                <li>1 < 7</li>
            
        
            
            
                <li>1 < 6</li>
            
        
            
            
                <li>1 < 5</li>
            
        
            
            
                <li>1 < 4</li>
            
        
            
            
                <li>1 < 3</li>
            
        
            
            
                <li>1 < 2</li>
            
        
            
            
                <li>1 <= 1</li>
            
        
            
            
                <li>1 >= 0</li>
            
        
    
        
            
            
                <li>0 <= 0</li>
            
        
            
                <li>9 > 8 && 0 < 8 || 9 < 1 && 0 > 1</li>
            
            
                <li>0 < 9</li>
            
        
            
            
                <li>0 < 8</li>
            
        
            
            
                <li>0 < 7</li>
            
        
            
            
                <li>0 < 6</li>
            
        
            
            
                <li>0 < 5</li>
            
        
            
            
                <li>0 < 4</li>
            
        
            
            
                <li>0 < 3</li>
            
        
            
            
                <li>0 < 2</li>
            
        
            
            
                <li>0 < 1</li>
            
        
            
            
                <li>0 <= 0</li>
            
        
    
    </ul>

    <h2>Математические операции</h2>
    <ul>
    
        
            <li>0 + 0 = 0</li>
            <li>0 - 0 = 0</li>
            <li>0 * 0 = 0</li>
            <li>0 / 0 = Error: divide by zero</li>
            <li>0 \ 0 = Error: divide by zero</li> 
            <li>0 % 0 = Error: divide by zero</li>
            <li>0 ^ 0 = 1</li>
        
            <li>0 + 9 = 9</li>
            <li>0 - 9 = -9</li>
            <li>0 * 9 = 0</li>
            <li>0 / 9 = 0</li>
            <li>0 \ 9 = 0</li> 
            <li>0 % 9 = 0</li>
            <li>0 ^ 9 = 0</li>
        
            <li>0 + 8 = 8</li>
            <li>0 - 8 = -8</li>
            <li>0 * 8 = 0</li>
            <li>0 / 8 = 0</li>
            <li>0 \ 8 = 0</li> 
            <li>0 % 8 = 0</li>
            <li>0 ^ 8 = 0</li>
        
            <li>0 + 7 = 7</li>
            <li>0 - 7 = -7</li>
            <li>0 * 7 = 0</li>
            <li>0 / 7 = 0</li>
            <li>0 \ 7 = 0</li> 
            <li>0 % 7 = 0</li>
            <li>0 ^ 7 = 0</li>
        
            <li>0 + 6 = 6</li>
            <li>0 - 6 = -6</li>
            <li>0 * 6 = 0</li>
            <li>0 / 6 = 0</li>
            <li>0 \ 6 = 0</li> 
            <li>0 % 6 = 0</li>
            <li>0 ^ 6 = 0</li>
        
            <li>0 + 5 = 5</li>
            <li>0 - 5 = -5</li>
            <li>0 * 5 = 0</li>
            <li>0 / 5 = 0</li>
            <li>0 \ 5 = 0</li> 
            <li>0 % 5 = 0</li>
            <li>0 ^ 5 = 0</li>
        
            <li>0 + 4 = 4</li>
            <li>0 - 4 = -4</li>
            <li>0 * 4 = 0</li>
            <li>0 / 4 = 0</li>
            <li>0 \ 4 = 0</li> 
            <li>0 % 4 = 0</li>
            <li>0 ^ 4 = 0</li>
        
            <li>0 + 3 = 3</li>
            <li>0 - 3 = -3</li>
            <li>0 * 3 = 0</li>
            <li>0 / 3 = 0</li>
            <li>0 \ 3 = 0</li> 
            <li>0 % 3 = 0</li>
            <li>0 ^ 3 = 0</li>
        
            <li>0 + 2 = 2</li>
            <li>0 - 2 = -2</li>
            <li>0 * 2 = 0</li>
            <li>0 / 2 = 0</li>
            <li>0 \ 2 = 0</li> 
            <li>0 % 2 = 0</li>
            <li>0 ^ 2 = 0</li>
        
            <li>0 + 1 = 1</li>
            <li>0 - 1 = -1</li>
            <li>0 * 1 = 0</li>
            <li>0 / 1 = 0</li>
            <li>0 \ 1 = 0</li> 
            <li>0 % 1 = 0</li>
            <li>0 ^ 1 = 0</li>
        
            <li>0 + 0 = 0</li>
            <li>0 - 0 = 0</li>
            <li>0 * 0 = 0</li>
            <li>0 / 0 = Error: divide by zero</li>
            <li>0 \ 0 = Error: divide by zero</li> 
            <li>0 % 0 = Error: divide by zero</li>
            <li>0 ^ 0 = 1</li>
        
    
        
            <li>9 + 0 = 9</li>
            <li>9 - 0 = 9</li>
            <li>9 * 0 = 0</li>
            <li>9 / 0 = Error: divide by zero</li>
            <li>9 \ 0 = Error: divide by zero</li> 
            <li>9 % 0 = Error: divide by zero</li>
            <li>9 ^ 0 = 1</li>
        
            <li>9 + 9 = 18</li>
            <li>9 - 9 = 0</li>
            <li>9 * 9 = 81</li>
            <li>9 / 9 = 1</li>
            <li>9 \ 9 = 1</li> 
            <li>9 % 9 = 0</li>
            <li>9 ^ 9 = 387420489</li>
        
            <li>9 + 8 = 17</li>
            <li>9 - 8 = 1</li>
            <li>9 * 8 = 72</li>
            <li>9 / 8 = 1.125</li>
            <li>9 \ 8 = 1</li> 
            <li>9 % 8 = 1</li>
            <li>9 ^ 8 = 43046721</li>
        
            <li>9 + 7 = 16</li>
            <li>9 - 7 = 2</li>
            <li>9 * 7 = 63</li>
            <li>9 / 7 = 1.2857142857142858</li>
            <li>9 \ 7 = 1</li> 
            <li>9 % 7 = 2</li>
            <li>9 ^ 7 = 4782969</li>
        
            <li>9 + 6 = 15</li>
            <li>9 - 6 = 3</li>
            <li>9 * 6 = 54</li>
            <li>9 / 6 = 1.5</li>
            <li>9 \ 6 = 1</li> 
            <li>9 % 6 = 3</li>
            <li>9 ^ 6 = 531441</li>
        
            <li>9 + 5 = 14</li>
            <li>9 - 5 = 4</li>
            <li>9 * 5 = 45</li>
            <li>9 / 5 = 1.8</li>
            <li>9 \ 5 = 1</li> 
            <li>9 % 5 = 4</li>
            <li>9 ^ 5 = 59049</li>
        
            <li>9 + 4 = 13</li>
            <li>9 - 4 = 5</li>
            <li>9 * 4 = 36</li>
            <li>9 / 4 = 2.25</li>
            <li>9 \ 4 = 2</li> 
            <li>9 % 4 = 1</li>
            <li>9 ^ 4 = 6561</li>
        
            <li>9 + 3 = 12</li>
            <li>9 - 3 = 6</li>
            <li>9 * 3 = 27</li>
            <li>9 / 3 = 3</li>
            <li>9 \ 3 = 3</li> 
            <li>9 % 3 = 0</li>
            <li>9 ^ 3 = 729</li>
        
            <li>9 + 2 = 11</li>
            <li>9 - 2 = 7</li>
            <li>9 * 2 = 18</li>
            <li>9 / 2 = 4.5</li>
            <li>9 \ 2 = 4</li> 
            <li>9 % 2 = 1</li>
            <li>9 ^ 2 = 81</li>
        
            <li>9 + 1 = 10</li>
            <li>9 - 1 = 8</li>
            <li>9 * 1 = 9</li>
            <li>9 / 1 = 9</li>
            <li>9 \ 1 = 9</li> 
            <li>9 % 1 = 0</li>
            <li>9 ^ 1 = 9</li>
        
            <li>9 + 0 = 9</li>
            <li>9 - 0 = 9</li>
            <li>9 * 0 = 0</li>
            <li>9 / 0 = Error: divide by zero</li>
            <li>9 \ 0 = Error: divide by zero</li> 
            <li>9 % 0 = Error: divide by zero</li>
            <li>9 ^ 0 = 1</li>
        
    
        
            <li>8 + 0 = 8</li>
            <li>8 - 0 = 8</li>
            <li>8 * 0 = 0</li>
            <li>8 / 0 = Error: divide by zero</li>
            <li>8 \ 0 = Error: divide by zero</li> 
            <li>8 % 0 = Error: divide by zero</li>
            <li>8 ^ 0 = 1</li>
        
            <li>8 + 9 = 17</li>
            <li>8 - 9 = -1</li>
            <li>8 * 9 = 72</li>
            <li>8 / 9 = 0.8888888888888888</li>
            <li>8 \ 9 = 0</li> 
            <li>8 % 9 = 8</li>
            <li>8 ^ 9 = 134217728</li>
        
            <li>8 + 8 = 16</li>
            <li>8 - 8 = 0</li>
            <li>8 * 8 = 64</li>
            <li>8 / 8 = 1</li>
            <li>8 \ 8 = 1</li> 
            <li>8 % 8 = 0</li>
            <li>8 ^ 8 = 16777216</li>
        
            <li>8 + 7 = 15</li>
            <li>8 - 7 = 1</li>
            <li>8 * 7 = 56</li>
            <li>8 / 7 = 1.1428571428571428</li>
            <li>8 \ 7 = 1</li> 
            <li>8 % 7 = 1</li>
            <li>8 ^ 7 = 2097152</li>
        
            <li>8 + 6 = 14</li>
            <li>8 - 6 = 2</li>
            <li>8 * 6 = 48</li>
            <li>8 / 6 = 1.3333333333333333</li>
            <li>8 \ 6 = 1</li> 
            <li>8 % 6 = 2</li>
            <li>8 ^ 6 = 262144</li>
        
            <li>8 + 5 = 13</li>
            <li>8 - 5 = 3</li>
            <li>8 * 5 = 40</li>
            <li>8 / 5 = 1.6</li>
            <li>8 \ 5 = 1</li> 
            <li>8 % 5 = 3</li>
            <li>8 ^ 5 = 32768</li>
        
            <li>8 + 4 = 12</li>
            <li>8 - 4 = 4</li>
            <li>8 * 4 = 32</li>
            <li>8 / 4 = 2</li>
            <li>8 \ 4 = 2</li> 
            <li>8 % 4 = 0</li>
            <li>8 ^ 4 = 4096</li>
        
            <li>8 + 3 = 11</li>
            <li>8 - 3 = 5</li>
            <li>8 * 3 = 24</li>
            <li>8 / 3 = 2.6666666666666665</li>
            <li>8 \ 3 = 2</li> 
            <li>8 % 3 = 2</li>
            <li>8 ^ 3 = 512</li>
        
            <li>8 + 2 = 10</li>
            <li>8 - 2 = 6</li>
            <li>8 * 2 = 16</li>
            <li>8 / 2 = 4</li>
            <li>8 \ 2 = 4</li> 
            <li>8 % 2 = 0</li>
            <li>8 ^ 2 = 64</li>
        
            <li>8 + 1 = 9</li>
            <li>8 - 1 = 7</li>
            <li>8 * 1 = 8</li>
            <li>8 / 1 = 8</li>
            <li>8 \ 1 = 8</li> 
            <li>8 % 1 = 0</li>
            <li>8 ^ 1 = 8</li>
        
            <li>8 + 0 = 8</li>
            <li>8 - 0 = 8</li>
            <li>8 * 0 = 0</li>
            <li>8 / 0 = Error: divide by zero</li>
            <li>8 \ 0 = Error: divide by zero</li> 
            <li>8 % 0 = Error: divide by zero</li>
            <li>8 ^ 0 = 1</li>
        
    
        
            <li>7 + 0 = 7</li>
            <li>7 - 0 = 7</li>
            <li>7 * 0 = 0</li>
            <li>7 / 0 = Error: divide by zero</li>
            <li>7 \ 0 = Error: divide by zero</li> 
            <li>7 % 0 = Error: divide by zero</li>
            <li>7 ^ 0 = 1</li>
        
            <li>7 + 9 = 16</li>
            <li>7 - 9 = -2</li>
            <li>7 * 9 = 63</li>
            <li>7 / 9 = 0.7777777777777778</li>
            <li>7 \ 9 = 0</li> 
            <li>7 % 9 = 7</li>
            <li>7 ^ 9 = 40353607</li>
        
            <li>7 + 8 = 15</li>
            <li>7 - 8 = -1</li>
            <li>7 * 8 = 56</li>
            <li>7 / 8 = 0.875</li>
            <li>7 \ 8 = 0</li> 
            <li>7 % 8 = 7</li>
            <li>7 ^ 8 = 5764801</li>
        
            <li>7 + 7 = 14</li>
            <li>7 - 7 = 0</li>
            <li>7 * 7 = 49</li>
            <li>7 / 7 = 1</li>
            <li>7 \ 7 = 1</li> 
            <li>7 % 7 = 0</li>
            <li>7 ^ 7 = 823543</li>
        
            <li>7 + 6 = 13</li>
            <li>7 - 6 = 1</li>
            <li>7 * 6 = 42</li>
            <li>7 / 6 = 1.1666666666666667</li>
            <li>7 \ 6 = 1</li> 
            <li>7 % 6 = 1</li>
            <li>7 ^ 6 = 117649</li>
        
            <li>7 + 5 = 12</li>
            <li>7 - 5 = 2</li>
            <li>7 * 5 = 35</li>
            <li>7 / 5 = 1.4</li>
            <li>7 \ 5 = 1</li> 
            <li>7 % 5 = 2</li>
            <li>7 ^ 5 = 16807</li>
        
            <li>7 + 4 = 11</li>
            <li>7 - 4 = 3</li>
            <li>7 * 4 = 28</li>
            <li>7 / 4 = 1.75</li>
            <li>7 \ 4 = 1</li> 
            <li>7 % 4 = 3</li>
            <li>7 ^ 4 = 2401</li>
        
            <li>7 + 3 = 10</li>
            <li>7 - 3 = 4</li>
            <li>7 * 3 = 21</li>
            <li>7 / 3 = 2.3333333333333335</li>
            <li>7 \ 3 = 2</li> 
            <li>7 % 3 = 1</li>
            <li>7 ^ 3 = 343</li>
        
            <li>7 + 2 = 9</li>
            <li>7 - 2 = 5</li>
            <li>7 * 2 = 14</li>
            <li>7 / 2 = 3.5</li>
            <li>7 \ 2 = 3</li> 
            <li>7 % 2 = 1</li>
            <li>7 ^ 2 = 49</li>
        
            <li>7 + 1 = 8</li>
            <li>7 - 1 = 6</li>
            <li>7 * 1 = 7</li>
            <li>7 / 1 = 7</li>
            <li>7 \ 1 = 7</li> 
            <li>7 % 1 = 0</li>
            <li>7 ^ 1 = 7</li>
        
            <li>7 + 0 = 7</li>
            <li>7 - 0 = 7</li>
            <li>7 * 0 = 0</li>
            <li>7 / 0 = Error: divide by zero</li>
            <li>7 \ 0 = Error: divide by zero</li> 
            <li>7 % 0 = Error: divide by zero</li>
            <li>7 ^ 0 = 1</li>
        
    
        
            <li>6 + 0 = 6</li>
            <li>6 - 0 = 6</li>
            <li>6 * 0 = 0</li>
            <li>6 / 0 = Error: divide by zero</li>
            <li>6 \ 0 = Error: divide by zero</li> 
            <li>6 % 0 = Error: divide by zero</li>
            <li>6 ^ 0 = 1</li>
        
            <li>6 + 9 = 15</li>
            <li>6 - 9 = -3</li>
            <li>6 * 9 = 54</li>
            <li>6 / 9 = 0.6666666666666666</li>
            <li>6 \ 9 = 0</li> 
            <li>6 % 9 = 6</li>
            <li>6 ^ 9 = 10077696</li>
        
            <li>6 + 8 = 14</li>
            <li>6 - 8 = -2</li>
            <li>6 * 8 = 48</li>
            <li>6 / 8 = 0.75</li>
            <li>6 \ 8 = 0</li> 
            <li>6 % 8 = 6</li>
            <li>6 ^ 8 = 1679616</li>
        
            <li>6 + 7 = 13</li>
            <li>6 - 7 = -1</li>
            <li>6 * 7 = 42</li>
            <li>6 / 7 = 0.8571428571428571</li>
            <li>6 \ 7 = 0</li> 
            <li>6 % 7 = 6</li>
            <li>6 ^ 7 = 279936</li>
        
            <li>6 + 6 = 12</li>
            <li>6 - 6 = 0</li>
            <li>6 * 6 = 36</li>
            <li>6 / 6 = 1</li>
            <li>6 \ 6 = 1</li> 
            <li>6 % 6 = 0</li>
            <li>6 ^ 6 = 46656</li>
        
            <li>6 + 5 = 11</li>
            <li>6 - 5 = 1</li>
            <li>6 * 5 = 30</li>
            <li>6 / 5 = 1.2</li>
            <li>6 \ 5 = 1</li> 
            <li>6 % 5 = 1</li>
            <li>6 ^ 5 = 7776</li>
        
            <li>6 + 4 = 10</li>
            <li>6 - 4 = 2</li>
            <li>6 * 4 = 24</li>
            <li>6 / 4 = 1.5</li>
            <li>6 \ 4 = 1</li> 
            <li>6 % 4 = 2</li>
            <li>6 ^ 4 = 1296</li>
        
            <li>6 + 3 = 9</li>
            <li>6 - 3 = 3</li>
            <li>6 * 3 = 18</li>
            <li>6 / 3 = 2</li>
            <li>6 \ 3 = 2</li> 
            <li>6 % 3 = 0</li>
            <li>6 ^ 3 = 216</li>
        
            <li>6 + 2 = 8</li>
            <li>6 - 2 = 4</li>
            <li>6 * 2 = 12</li>
            <li>6 / 2 = 3</li>
            <li>6 \ 2 = 3</li> 
            <li>6 % 2 = 0</li>
            <li>6 ^ 2 = 36</li>
        
            <li>6 + 1 = 7</li>
            <li>6 - 1 = 5</li>
            <li>6 * 1 = 6</li>
            <li>6 / 1 = 6</li>
            <li>6 \ 1 = 6</li> 
            <li>6 % 1 = 0</li>
            <li>6 ^ 1 = 6</li>
        
            <li>6 + 0 = 6</li>
            <li>6 - 0 = 6</li>
            <li>6 * 0 = 0</li>
            <li>6 / 0 = Error: divide by zero</li>
            <li>6 \ 0 = Error: divide by zero</li> 
            <li>6 % 0 = Error: divide by zero</li>
            <li>6 ^ 0 = 1</li>
        
    
        
            <li>5 + 0 = 5</li>
            <li>5 - 0 = 5</li>
            <li>5 * 0 = 0</li>
            <li>5 / 0 = Error: divide by zero</li>
            <li>5 \ 0 = Error: divide by zero</li> 
            <li>5 % 0 = Error: divide by zero</li>
            <li>5 ^ 0 = 1</li>
        
            <li>5 + 9 = 14</li>
            <li>5 - 9 = -4</li>
            <li>5 * 9 = 45</li>
            <li>5 / 9 = 0.5555555555555556</li>
            <li>5 \ 9 = 0</li> 
            <li>5 % 9 = 5</li>
            <li>5 ^ 9 = 1953125</li>
        
            <li>5 + 8 = 13</li>
            <li>5 - 8 = -3</li>
            <li>5 * 8 = 40</li>
            <li>5 / 8 = 0.625</li>
            <li>5 \ 8 = 0</li> 
            <li>5 % 8 = 5</li>
            <li>5 ^ 8 = 390625</li>
        
            <li>5 + 7 = 12</li>
            <li>5 - 7 = -2</li>
            <li>5 * 7 = 35</li>
            <li>5 / 7 = 0.7142857142857143</li>
            <li>5 \ 7 = 0</li> 
            <li>5 % 7 = 5</li>
            <li>5 ^ 7 = 78125</li>
        
            <li>5 + 6 = 11</li>
            <li>5 - 6 = -1</li>
            <li>5 * 6 = 30</li>
            <li>5 / 6 = 0.8333333333333334</li>
            <li>5 \ 6 = 0</li> 
            <li>5 % 6 = 5</li>
            <li>5 ^ 6 = 15625</li>
        
            <li>5 + 5 = 10</li>
            <li>5 - 5 = 0</li>
            <li>5 * 5 = 25</li>
            <li>5 / 5 = 1</li>
            <li>5 \ 5 = 1</li> 
            <li>5 % 5 = 0</li>
            <li>5 ^ 5 = 3125</li>
        
            <li>5 + 4 = 9</li>
            <li>5 - 4 = 1</li>
            <li>5 * 4 = 20</li>
            <li>5 / 4 = 1.25</li>
            <li>5 \ 4 = 1</li> 
            <li>5 % 4 = 1</li>
            <li>5 ^ 4 = 625</li>
        
            <li>5 + 3 = 8</li>
            <li>5 - 3 = 2</li>
            <li>5 * 3 = 15</li>
            <li>5 / 3 = 1.6666666666666667</li>
            <li>5 \ 3 = 1</li> 
            <li>5 % 3 = 2</li>
            <li>5 ^ 3 = 125</li>
        
            <li>5 + 2 = 7</li>
            <li>5 - 2 = 3</li>
            <li>5 * 2 = 10</li>
            <li>5 / 2 = 2.5</li>
            <li>5 \ 2 = 2</li> 
            <li>5 % 2 = 1</li>
            <li>5 ^ 2 = 25</li>
        
            <li>5 + 1 = 6</li>
            <li>5 - 1 = 4</li>
            <li>5 * 1 = 5</li>
            <li>5 / 1 = 5</li>
            <li>5 \ 1 = 5</li> 
            <li>5 % 1 = 0</li>
            <li>5 ^ 1 = 5</li>
        
            <li>5 + 0 = 5</li>
            <li>5 - 0 = 5</li>
            <li>5 * 0 = 0</li>
            <li>5 / 0 = Error: divide by zero</li>
            <li>5 \ 0 = Error: divide by zero</li> 
            <li>5 % 0 = Error: divide by zero</li>
            <li>5 ^ 0 = 1</li>
        
    
        
            <li>4 + 0 = 4</li>
            <li>4 - 0 = 4</li>
            <li>4 * 0 = 0</li>
            <li>4 / 0 = Error: divide by zero</li>
            <li>4 \ 0 = Error: divide by zero</li> 
            <li>4 % 0 = Error: divide by zero</li>
            <li>4 ^ 0 = 1</li>
        
            <li>4 + 9 = 13</li>
            <li>4 - 9 = -5</li>
            <li>4 * 9 = 36</li>
            <li>4 / 9 = 0.4444444444444444</li>
            <li>4 \ 9 = 0</li> 
            <li>4 % 9 = 4</li>
            <li>4 ^ 9 = 262144</li>
        
            <li>4 + 8 = 12</li>
            <li>4 - 8 = -4</li>
            <li>4 * 8 = 32</li>
            <li>4 / 8 = 0.5</li>
            <li>4 \ 8 = 0</li> 
            <li>4 % 8 = 4</li>
            <li>4 ^ 8 = 65536</li>
        
            <li>4 + 7 = 11</li>
            <li>4 - 7 = -3</li>
            <li>4 * 7 = 28</li>
            <li>4 / 7 = 0.5714285714285714</li>
            <li>4 \ 7 = 0</li> 
            <li>4 % 7 = 4</li>
            <li>4 ^ 7 = 16384</li>
        
            <li>4 + 6 = 10</li>
            <li>4 - 6 = -2</li>
            <li>4 * 6 = 24</li>
            <li>4 / 6 = 0.6666666666666666</li>
            <li>4 \ 6 = 0</li> 
            <li>4 % 6 = 4</li>
            <li>4 ^ 6 = 4096</li>
        
            <li>4 + 5 = 9</li>
            <li>4 - 5 = -1</li>
            <li>4 * 5 = 20</li>
            <li>4 / 5 = 0.8</li>
            <li>4 \ 5 = 0</li> 
            <li>4 % 5 = 4</li>
            <li>4 ^ 5 = 1024</li>
        
            <li>4 + 4 = 8</li>
            <li>4 - 4 = 0</li>
            <li>4 * 4 = 16</li>
            <li>4 / 4 = 1</li>
            <li>4 \ 4 = 1</li> 
            <li>4 % 4 = 0</li>
            <li>4 ^ 4 = 256</li>
        
            <li>4 + 3 = 7</li>
            <li>4 - 3 = 1</li>
            <li>4 * 3 = 12</li>
            <li>4 / 3 = 1.3333333333333333</li>
            <li>4 \ 3 = 1</li> 
            <li>4 % 3 = 1</li>
            <li>4 ^ 3 = 64</li>
        
            <li>4 + 2 = 6</li>
            <li>4 - 2 = 2</li>
            <li>4 * 2 = 8</li>
            <li>4 / 2 = 2</li>
            <li>4 \ 2 = 2</li> 
            <li>4 % 2 = 0</li>
            <li>4 ^ 2 = 16</li>
        
            <li>4 + 1 = 5</li>
            <li>4 - 1 = 3</li>
            <li>4 * 1 = 4</li>
            <li>4 / 1 = 4</li>
            <li>4 \ 1 = 4</li> 
            <li>4 % 1 = 0</li>
            <li>4 ^ 1 = 4</li>
        
            <li>4 + 0 = 4</li>
            <li>4 - 0 = 4</li>
            <li>4 * 0 = 0</li>
            <li>4 / 0 = Error: divide by zero</li>
            <li>4 \ 0 = Error: divide by zero</li> 
            <li>4 % 0 = Error: divide by zero</li>
            <li>4 ^ 0 = 1</li>
        
    
        
            <li>3 + 0 = 3</li>
            <li>3 - 0 = 3</li>
            <li>3 * 0 = 0</li>
            <li>3 / 0 = Error: divide by zero</li>
            <li>3 \ 0 = Error: divide by zero</li> 
            <li>3 % 0 = Error: divide by zero</li>
            <li>3 ^ 0 = 1</li>
        
            <li>3 + 9 = 12</li>
            <li>3 - 9 = -6</li>
            <li>3 * 9 = 27</li>
            <li>3 / 9 = 0.3333333333333333</li>
            <li>3 \ 9 = 0</li> 
            <li>3 % 9 = 3</li>
            <li>3 ^ 9 = 19683</li>
        
            <li>3 + 8 = 11</li>
            <li>3 - 8 = -5</li>
            <li>3 * 8 = 24</li>
            <li>3 / 8 = 0.375</li>
            <li>3 \ 8 = 0</li> 
            <li>3 % 8 = 3</li>
            <li>3 ^ 8 = 6561</li>
        
            <li>3 + 7 = 10</li>
            <li>3 - 7 = -4</li>
            <li>3 * 7 = 21</li>
            <li>3 / 7 = 0.42857142857142855</li>
            <li>3 \ 7 = 0</li> 
            <li>3 % 7 = 3</li>
            <li>3 ^ 7 = 2187</li>
        
            <li>3 + 6 = 9</li>
            <li>3 - 6 = -3</li>
            <li>3 * 6 = 18</li>
            <li>3 / 6 = 0.5</li>
            <li>3 \ 6 = 0</li> 
            <li>3 % 6 = 3</li>
            <li>3 ^ 6 = 729</li>
        
            <li>3 + 5 = 8</li>
            <li>3 - 5 = -2</li>
            <li>3 * 5 = 15</li>
            <li>3 / 5 = 0.6</li>
            <li>3 \ 5 = 0</li> 
            <li>3 % 5 = 3</li>
            <li>3 ^ 5 = 243</li>
        
            <li>3 + 4 = 7</li>
            <li>3 - 4 = -1</li>
            <li>3 * 4 = 12</li>
            <li>3 / 4 = 0.75</li>
            <li>3 \ 4 = 0</li> 
            <li>3 % 4 = 3</li>
            <li>3 ^ 4 = 81</li>
        
            <li>3 + 3 = 6</li>
            <li>3 - 3 = 0</li>
            <li>3 * 3 = 9</li>
            <li>3 / 3 = 1</li>
            <li>3 \ 3 = 1</li> 
            <li>3 % 3 = 0</li>
            <li>3 ^ 3 = 27</li>
        
            <li>3 + 2 = 5</li>
            <li>3 - 2 = 1</li>
            <li>3 * 2 = 6</li>
            <li>3 / 2 = 1.5</li>
            <li>3 \ 2 = 1</li> 
            <li>3 % 2 = 1</li>
            <li>3 ^ 2 = 9</li>
        
            <li>3 + 1 = 4</li>
            <li>3 - 1 = 2</li>
            <li>3 * 1 = 3</li>
            <li>3 / 1 = 3</li>
            <li>3 \ 1 = 3</li> 
            <li>3 % 1 = 0</li>
            <li>3 ^ 1 = 3</li>
        
            <li>3 + 0 = 3</li>
            <li>3 - 0 = 3</li>
            <li>3 * 0 = 0</li>
            <li>3 / 0 = Error: divide by zero</li>
            <li>3 \ 0 = Error: divide by zero</li> 
            <li>3 % 0 = Error: divide by zero</li>
            <li>3 ^ 0 = 1</li>
        
    
        
            <li>2 + 0 = 2</li>
            <li>2 - 0 = 2</li>
            <li>2 * 0 = 0</li>
            <li>2 / 0 = Error: divide by zero</li>
            <li>2 \ 0 = Error: divide by zero</li> 
            <li>2 % 0 = Error: divide by zero</li>
            <li>2 ^ 0 = 1</li>
        
            <li>2 + 9 = 11</li>
            <li>2 - 9 = -7</li>
            <li>2 * 9 = 18</li>
            <li>2 / 9 = 0.2222222222222222</li>
            <li>2 \ 9 = 0</li> 
            <li>2 % 9 = 2</li>
            <li>2 ^ 9 = 512</li>
        
            <li>2 + 8 = 10</li>
            <li>2 - 8 = -6</li>
            <li>2 * 8 = 16</li>
            <li>2 / 8 = 0.25</li>
            <li>2 \ 8 = 0</li> 
            <li>2 % 8 = 2</li>
            <li>2 ^ 8 = 256</li>
        
            <li>2 + 7 = 9</li>
            <li>2 - 7 = -5</li>
            <li>2 * 7 = 14</li>
            <li>2 / 7 = 0.2857142857142857</li>
            <li>2 \ 7 = 0</li> 
            <li>2 % 7 = 2</li>
            <li>2 ^ 7 = 128</li>
        
            <li>2 + 6 = 8</li>
            <li>2 - 6 = -4</li>
            <li>2 * 6 = 12</li>
            <li>2 / 6 = 0.3333333333333333</li>
            <li>2 \ 6 = 0</li> 
            <li>2 % 6 = 2</li>
            <li>2 ^ 6 = 64</li>
        
            <li>2 + 5 = 7</li>
            <li>2 - 5 = -3</li>
            <li>2 * 5 = 10</li>
            <li>2 / 5 = 0.4</li>
            <li>2 \ 5 = 0</li> 
            <li>2 % 5 = 2</li>
            <li>2 ^ 5 = 32</li>
        
            <li>2 + 4 = 6</li>
            <li>2 - 4 = -2</li>
            <li>2 * 4 = 8</li>
            <li>2 / 4 = 0.5</li>
            <li>2 \ 4 = 0</li> 
            <li>2 % 4 = 2</li>
            <li>2 ^ 4 = 16</li>
        
            <li>2 + 3 = 5</li>
            <li>2 - 3 = -1</li>
            <li>2 * 3 = 6</li>
            <li>2 / 3 = 0.6666666666666666</li>
            <li>2 \ 3 = 0</li> 
            <li>2 % 3 = 2</li>
            <li>2 ^ 3 = 8</li>
        
            <li>2 + 2 = 4</li>
            <li>2 - 2 = 0</li>
            <li>2 * 2 = 4</li>
            <li>2 / 2 = 1</li>
            <li>2 \ 2 = 1</li> 
            <li>2 % 2 = 0</li>
            <li>2 ^ 2 = 4</li>
        
            <li>2 + 1 = 3</li>
            <li>2 - 1 = 1</li>
            <li>2 * 1 = 2</li>
            <li>2 / 1 = 2</li>
            <li>2 \ 1 = 2</li> 
            <li>2 % 1 = 0</li>
            <li>2 ^ 1 = 2</li>
        
            <li>2 + 0 = 2</li>
            <li>2 - 0 = 2</li>
            <li>2 * 0 = 0</li>
            <li>2 / 0 = Error: divide by zero</li>
            <li>2 \ 0 = Error: divide by zero</li> 
            <li>2 % 0 = Error: divide by zero</li>
            <li>2 ^ 0 = 1</li>
        
    
        
            <li>1 + 0 = 1</li>
            <li>1 - 0 = 1</li>
            <li>1 * 0 = 0</li>
            <li>1 / 0 = Error: divide by zero</li>
            <li>1 \ 0 = Error: divide by zero</li> 
            <li>1 % 0 = Error: divide by zero</li>
            <li>1 ^ 0 = 1</li>
        
            <li>1 + 9 = 10</li>
            <li>1 - 9 = -8</li>
            <li>1 * 9 = 9</li>
            <li>1 / 9 = 0.1111111111111111</li>
            <li>1 \ 9 = 0</li> 
            <li>1 % 9 = 1</li>
            <li>1 ^ 9 = 1</li>
        
            <li>1 + 8 = 9</li>
            <li>1 - 8 = -7</li>
            <li>1 * 8 = 8</li>
            <li>1 / 8 = 0.125</li>
            <li>1 \ 8 = 0</li> 
            <li>1 % 8 = 1</li>
            <li>1 ^ 8 = 1</li>
        
            <li>1 + 7 = 8</li>
            <li>1 - 7 = -6</li>
            <li>1 * 7 = 7</li>
            <li>1 / 7 = 0.14285714285714285</li>
            <li>1 \ 7 = 0</li> 
            <li>1 % 7 = 1</li>
            <li>1 ^ 7 = 1</li>
        
            <li>1 + 6 = 7</li>
            <li>1 - 6 = -5</li>
            <li>1 * 6 = 6</li>
            <li>1 / 6 = 0.16666666666666666</li>
            <li>1 \ 6 = 0</li> 
            <li>1 % 6 = 1</li>
            <li>1 ^ 6 = 1</li>
        
            <li>1 + 5 = 6</li>
            <li>1 - 5 = -4</li>
            <li>1 * 5 = 5</li>
            <li>1 / 5 = 0.2</li>
            <li>1 \ 5 = 0</li> 
            <li>1 % 5 = 1</li>
            <li>1 ^ 5 = 1</li>
        
            <li>1 + 4 = 5</li>
            <li>1 - 4 = -3</li>
            <li>1 * 4 = 4</li>
            <li>1 / 4 = 0.25</li>
            <li>1 \ 4 = 0</li> 
            <li>1 % 4 = 1</li>
            <li>1 ^ 4 = 1</li>
        
            <li>1 + 3 = 4</li>
            <li>1 - 3 = -2</li>
            <li>1 * 3 = 3</li>
            <li>1 / 3 = 0.3333333333333333</li>
            <li>1 \ 3 = 0</li> 
            <li>1 % 3 = 1</li>
            <li>1 ^ 3 = 1</li>
        
            <li>1 + 2 = 3</li>
            <li>1 - 2 = -1</li>
            <li>1 * 2 = 2</li>
            <li>1 / 2 = 0.5</li>
            <li>1 \ 2 = 0</li> 
            <li>1 % 2 = 1</li>
            <li>1 ^ 2 = 1</li>
        
            <li>1 + 1 = 2</li>
            <li>1 - 1 = 0</li>
            <li>1 * 1 = 1</li>
            <li>1 / 1 = 1</li>
            <li>1 \ 1 = 1</li> 
            <li>1 % 1 = 0</li>
            <li>1 ^ 1 = 1</li>
        
            <li>1 + 0 = 1</li>
            <li>1 - 0 = 1</li>
            <li>1 * 0 = 0</li>
            <li>1 / 0 = Error: divide by zero</li>
            <li>1 \ 0 = Error: divide by zero</li> 
            <li>1 % 0 = Error: divide by zero</li>
            <li>1 ^ 0 = 1</li>
        
    
        
            <li>0 + 0 = 0</li>
            <li>0 - 0 = 0</li>
            <li>0 * 0 = 0</li>
            <li>0 / 0 = Error: divide by zero</li>
            <li>0 \ 0 = Error: divide by zero</li> 
            <li>0 % 0 = Error: divide by zero</li>
            <li>0 ^ 0 = 1</li>
        
            <li>0 + 9 = 9</li>
            <li>0 - 9 = -9</li>
            <li>0 * 9 = 0</li>
            <li>0 / 9 = 0</li>
            <li>0 \ 9 = 0</li> 
            <li>0 % 9 = 0</li>
            <li>0 ^ 9 = 0</li>
        
            <li>0 + 8 = 8</li>
            <li>0 - 8 = -8</li>
            <li>0 * 8 = 0</li>
            <li>0 / 8 = 0</li>
            <li>0 \ 8 = 0</li> 
            <li>0 % 8 = 0</li>
            <li>0 ^ 8 = 0</li>
        
            <li>0 + 7 = 7</li>
            <li>0 - 7 = -7</li>
            <li>0 * 7 = 0</li>
            <li>0 / 7 = 0</li>
            <li>0 \ 7 = 0</li> 
            <li>0 % 7 = 0</li>
            <li>0 ^ 7 = 0</li>
        
            <li>0 + 6 = 6</li>
            <li>0 - 6 = -6</li>
            <li>0 * 6 = 0</li>
            <li>0 / 6 = 0</li>
            <li>0 \ 6 = 0</li> 
            <li>0 % 6 = 0</li>
            <li>0 ^ 6 = 0</li>
        
            <li>0 + 5 = 5</li>
            <li>0 - 5 = -5</li>
            <li>0 * 5 = 0</li>
            <li>0 / 5 = 0</li>
            <li>0 \ 5 = 0</li> 
            <li>0 % 5 = 0</li>
            <li>0 ^ 5 = 0</li>
        
            <li>0 + 4 = 4</li>
            <li>0 - 4 = -4</li>
            <li>0 * 4 = 0</li>
            <li>0 / 4 = 0</li>
            <li>0 \ 4 = 0</li> 
            <li>0 % 4 = 0</li>
            <li>0 ^ 4 = 0</li>
        
            <li>0 + 3 = 3</li>
            <li>0 - 3 = -3</li>
            <li>0 * 3 = 0</li>
            <li>0 / 3 = 0</li>
            <li>0 \ 3 = 0</li> 
            <li>0 % 3 = 0</li>
            <li>0 ^ 3 = 0</li>
        
            <li>0 + 2 = 2</li>
            <li>0 - 2 = -2</li>
            <li>0 * 2 = 0</li>
            <li>0 / 2 = 0</li>
            <li>0 \ 2 = 0</li> 
            <li>0 % 2 = 0</li>
            <li>0 ^ 2 = 0</li>
        
            <li>0 + 1 = 1</li>
            <li>0 - 1 = -1</li>
            <li>0 * 1 = 0</li>
            <li>0 / 1 = 0</li>
            <li>0 \ 1 = 0</li> 
            <li>0 % 1 = 0</li>
            <li>0 ^ 1 = 0</li>
        
            <li>0 + 0 = 0</li>
            <li>0 - 0 = 0</li>
            <li>0 * 0 = 0</li>
            <li>0 / 0 = Error: divide by zero</li>
            <li>0 \ 0 = Error: divide by zero</li> 
            <li>0 % 0 = Error: divide by zero</li>
            <li>0 ^ 0 = 1</li>
        
    
    </ul>

    <h2>Математические операции (краткая запись с присвоением)</h2>
    <ul>
    
        <li>
            $myVar = 1<br> 
            $myVar += 6<br> 
            $myVar = 7
        </li>
        <li>
            $myVar = 1<br> 
            $myVar -= 6<br> 
            $myVar = -5
        </li>
        <li>
            $myVar = 1<br> 
            $myVar *= 6<br> 
            $myVar = 6
        </li>
        <li>
            $myVar = 1<br> 
            $myVar /= 6<br> 
            $myVar = 0.16666666666666666
        </li>
        <li>
            $myVar = 1<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 1<br> 
            $myVar %= 6<br> 
            $myVar = 1
        </li>
        <li>
            $myVar = 1<br> 
            $myVar ^= 6<br> 
            $myVar = 1
        </li>
    
        <li>
            $myVar = 2<br> 
            $myVar += 6<br> 
            $myVar = 8
        </li>
        <li>
            $myVar = 2<br> 
            $myVar -= 6<br> 
            $myVar = -4
        </li>
        <li>
            $myVar = 2<br> 
            $myVar *= 6<br> 
            $myVar = 12
        </li>
        <li>
            $myVar = 2<br> 
            $myVar /= 6<br> 
            $myVar = 0.3333333333333333
        </li>
        <li>
            $myVar = 2<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 2<br> 
            $myVar %= 6<br> 
            $myVar = 2
        </li>
        <li>
            $myVar = 2<br> 
            $myVar ^= 6<br> 
            $myVar = 64
        </li>
    
        <li>
            $myVar = 3<br> 
            $myVar += 6<br> 
            $myVar = 9
        </li>
        <li>
            $myVar = 3<br> 
            $myVar -= 6<br> 
            $myVar = -3
        </li>
        <li>
            $myVar = 3<br> 
            $myVar *= 6<br> 
            $myVar = 18
        </li>
        <li>
            $myVar = 3<br> 
            $myVar /= 6<br> 
            $myVar = 0.5
        </li>
        <li>
            $myVar = 3<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 3<br> 
            $myVar %= 6<br> 
            $myVar = 3
        </li>
        <li>
            $myVar = 3<br> 
            $myVar ^= 6<br> 
            $myVar = 729
        </li>
    
        <li>
            $myVar = 4<br> 
            $myVar += 6<br> 
            $myVar = 10
        </li>
        <li>
            $myVar = 4<br> 
            $myVar -= 6<br> 
            $myVar = -2
        </li>
        <li>
            $myVar = 4<br> 
            $myVar *= 6<br> 
            $myVar = 24
        </li>
        <li>
            $myVar = 4<br> 
            $myVar /= 6<br> 
            $myVar = 0.6666666666666666
        </li>
        <li>
            $myVar = 4<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 4<br> 
            $myVar %= 6<br> 
            $myVar = 4
        </li>
        <li>
            $myVar = 4<br> 
            $myVar ^= 6<br> 
            $myVar = 4096
        </li>
    
        <li>
            $myVar = 5<br> 
            $myVar += 6<br> 
            $myVar = 11
        </li>
        <li>
            $myVar = 5<br> 
            $myVar -= 6<br> 
            $myVar = -1
        </li>
        <li>
            $myVar = 5<br> 
            $myVar *= 6<br> 
            $myVar = 30
        </li>
        <li>
            $myVar = 5<br> 
            $myVar /= 6<br> 
            $myVar = 0.8333333333333334
        </li>
        <li>
            $myVar = 5<br> 
            $myVar \= 6<br> 
            $myVar = 0
        </li>
        <li>
            $myVar = 5<br> 
            $myVar %= 6<br> 
            $myVar = 5
        </li>
        <li>
            $myVar = 5<br> 
            $myVar ^= 6<br> 
            $myVar = 15625
        </li>
    
    </ul>

    <h2>Конкатенация строк</h2>
    <ul>
    
        
            <li>Первый 2 Первый 2</li>
        
            <li>Первый 2 Второй 2</li>
        
            <li>Первый 2 Третий 2</li>
        
            <li>Первый 2 07869786</li>
        
    
        
            <li>Второй 2 Первый 2</li>
        
            <li>Второй 2 Второй 2</li>
        
            <li>Второй 2 Третий 2</li>
        
            <li>Второй 2 07869786</li>
        
    
        
            <li>Третий 2 Первый 2</li>
        
            <li>Третий 2 Второй 2</li>
        
            <li>Третий 2 Третий 2</li>
        
            <li>Третий 2 07869786</li>
        
    
        
            <li>07869786 Первый 2</li>
        
            <li>07869786 Второй 2</li>
        
            <li>07869786 Третий 2</li>
        
            <li>07869786 07869786</li>
        
    
    </ul>

    <h2>Доступ к полям структур</h2>
    <ul>
    
        <li>
            <b>0:</b>
<pre>
struct {
    ID: 3,
    Name: my name 2,
    Value: my value 2,
    Function: Мама мыла раму 3 часов
}
</pre>
        </li>
    
        <li>
            <b>1:</b>
<pre>
struct {
    ID: 4,
    Name: my name 2 2,
    Value: my value 2 2,
    Function: Мама мыла раму 4 часов
}
</pre>
        </li>
    
    </ul>

    <h2>Вызов пользовательской функции</h2>
    <ul>
        <li>Присвоим результат выполнения функции переменной $myData </li>
        <li>В переменной $myData хранится значение 1945-05-09 00:43</li>
        <li>Выведем результат выполнения функции сразу в шаблон 1945-05-09 00:43 или так 1945-05-09 00:43</li>
    </ul>

    <h2>Комментарии</h2>
    
    

    <h2>Подключение шаблонов</h2>
    <div>layouts/include.tpl Подключен</div>

    <h1>Вывод кода шаблона без обработки</h2>
<pre>
    
    @foreach($numbers as $number1)
        @foreach($numbers as $number2)
            &lt;li&gt;{{$number1}} + {{$number2}} = {{$number1 + $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} - {{$number2}} = {{$number1 - $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} * {{$number2}} = {{$number1 * $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} / {{$number2}} = {{$number1 / $number2}}&lt;/li&gt;
            &lt;li&gt;{{$number1}} % {{$number2}} = {{$number1 % $number2}}&lt;/li&gt;
        @endforeach
    @endforeach
    
</pre>
</body>
</html>`

func TestXtplCollection_View(t *testing.T) {
	xtpl.ViewsPath("./templates_test")
	//xtpl.ViewExtension("tpl")
	xtpl.CycleLimit(100)
	xtpl.Debug(false)
	xtpl.Functions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	var wg = sync.WaitGroup{}
	wg.Add(200)

	for i := 0; i < 100; i++ {
		go func() {
			var buff = &bytes.Buffer{}
			if err := xtpl.View("index", testData1, buff); err != nil {
				t.Error(err)
			} else if result1 != buff.String() {
				t.Error("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
			}
			wg.Done()
		}()
		go func() {
			var buff = &bytes.Buffer{}
			if err := xtpl.View("index", testData2, buff); err != nil {
				t.Error(err)
			} else if result2 != buff.String() {
				t.Error("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestXtplCollection_ParseString(t *testing.T) {
	xtpl.ViewsPath("./templates_test")
	//xtpl.ViewExtension("tpl")
	xtpl.CycleLimit(100)
	xtpl.Debug(false)
	xtpl.Functions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	var wg = sync.WaitGroup{}
	wg.Add(200)

	var source, err = ioutil.ReadFile("./templates_test/index.tpl")
	if err != nil {
		t.Error("Не удалось прочитать файл шаблона")
		return
	}

	for i := 0; i < 100; i++ {
		go func() {
			var buff = &bytes.Buffer{}
			if err := xtpl.ParseString(string(source), testData1, buff); err != nil {
				t.Error(err)
			} else if result1 != buff.String() {
				t.Error("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
			}
			wg.Done()
		}()
		go func() {
			var buff = &bytes.Buffer{}
			if err := xtpl.ParseString(string(source), testData2, buff); err != nil {
				t.Error(err)
			} else if result2 != buff.String() {
				t.Error("Возможно, что-то пошло не так, результат обработки шаблона не совпадает с образцом")
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkView(b *testing.B) {
	xtpl.ViewsPath("./templates_test")
	xtpl.ViewExtension("tpl")
	xtpl.CycleLimit(100)
	xtpl.Debug(false)
	xtpl.Functions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	for i := 0; i < b.N; i++ {
		var buff = &bytes.Buffer{}
		if err := xtpl.View("index", testData1, buff); err != nil {
			b.Error(err)
		}
	}
}
