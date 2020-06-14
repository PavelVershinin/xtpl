@extends("layouts/app")

@section("content")
    <h1>{{$page_title}}<h1/>

    <h2>Вывод переменных в шаблон</h2>
    <ul>
        <li>{{"<b>Экранированный вывод</b>"}}</li>
        <li>{!! "<b>НЕ экранированный вывод</b>" !!}</li>
    </ul>

    <h2>Объявление переменных и и прочие операции, без вывода результата в шаблон</h2>
    <ul>
        <li>Тут объявляю переменную и присваиваю ей результат математической операции, в шаблоне не вывожу @exec($myVar = 5 + 45)</li>
        <li>Тут выведу значение сохранённое в переменной {{$myVar}}</li>
    </ul>

    <h2>Цикл for</h2>
    <ul>
    @for ($i = 0; $i < 10; $i++)
        <li>
            @if ($strings.$i != "")
                В переменной $strings, под индексом {{$i}} содержится значение {{$strings.$i}}
            @else
                В переменной $strings, под индексом {{$i}} нет значения
            @endif
        </li>
    @endfor
    </ul>

    <h2>Цикл foreach</h2>
    <ul>
    @foreach ($strings as $i => $value)
        <li>$strings.{{$i}} == {{$value}}</li>
    @endforeach
    </ul>

    <h2>Условия</h2>
    <ul>
    @foreach($numbers as $number1)
        @foreach($numbers as $number2)
            @if($number2 > 8 && $number1 < 8 || $number2 < 1 && $number1 > 1)
                <li>{{$number2}} > 8 && {{$number1}} < 8 || {{$number2}} < 1 && {{$number1}} > 1</li>
            @endif
            @if ($number1 < $number2)
                <li>{{$number1}} < {{$number2}}</li>
            @elseif ($number1 <= $number2)
                <li>{{$number1}} <= {{$number2}}</li>
            @elseif ($number1 == $number2)
                <li>{{$number1}} == {{$number2}}</li>
            @elseif ($number1 >= $number2)
                <li>{{$number1}} >= {{$number2}}</li>
            @elseif ($number1 > $number2)
                <li>{{$number1}} > {{$number2}}</li>
            @elseif ($number1 != $number2)
                <li>{{$number1}} != {{$number2}}</li>
            @elseif ($number1 <> $number2)
                <li>{{$number1}} <> {{$number2}}</li>
            @endif
        @endforeach
    @endforeach
    </ul>

    <h2>Математические операции</h2>
    <ul>
    @foreach($numbers as $number1)
        @foreach($numbers as $number2)
            <li>{{$number1}} + {{$number2}} = {{$number1 + $number2}}</li>
            <li>{{$number1}} - {{$number2}} = {{$number1 - $number2}}</li>
            <li>{{$number1}} * {{$number2}} = {{$number1 * $number2}}</li>
            <li>{{$number1}} / {{$number2}} = {{$number1 / $number2}}</li>
            <li>{{$number1}} \ {{$number2}} = {{$number1 \ $number2}}</li> {* Деление без остатка, с отбрасыванием дробной части *}
            <li>{{$number1}} % {{$number2}} = {{$number1 % $number2}}</li>
            <li>{{$number1}} ^ {{$number2}} = {{$number1 ^ $number2}}</li>
        @endforeach
    @endforeach
    </ul>

    <h2>Математические операции (краткая запись с присвоением)</h2>
    <ul>
    @for($i = 1; $i <= 5; $i++)
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar += 6<br> @exec($myVar += 6)
            $myVar = {{$myVar}}
        </li>
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar -= 6<br> @exec($myVar -= 6)
            $myVar = {{$myVar}}
        </li>
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar *= 6<br> @exec($myVar *= 6)
            $myVar = {{$myVar}}
        </li>
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar /= 6<br> @exec($myVar /= 6)
            $myVar = {{$myVar}}
        </li>
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar \= 6<br> @exec($myVar \= 6)
            $myVar = {{$myVar}}
        </li>
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar %= 6<br> @exec($myVar %= 6)
            $myVar = {{$myVar}}
        </li>
        <li>
            $myVar = {{$i}}<br> @exec($myVar = $i)
            $myVar ^= 6<br> @exec($myVar ^= 6)
            $myVar = {{$myVar}}
        </li>
    @endfor
    </ul>

    <h2>Конкатенация строк</h2>
    <ul>
    @foreach($strings as $string1)
        @foreach($strings as $string2)
            <li>{{$string1 + " " + $string2}}</li>
        @endforeach
    @endforeach
    </ul>

    <h2>Доступ к полям структур</h2>
    <ul>
    @foreach($structs as $i => $struct)
        <li>
            <b>{{$i}}:</b>
<pre>
struct {
    ID: {{$struct.ID}},
    Name: {{$struct.Name}},
    Value: {{$struct.Value}},
    Function: {{$struct.Function("Мама мыла раму", $struct.ID)}}
}
</pre>
        </li>
    @endforeach
    </ul>

    <h2>Вызов пользовательской функции</h2>
    <ul>
        <li>Присвоим результат выполнения функции переменной $myData @exec($myData = date(-777856620, "2006-01-02 15:04"))</li>
        <li>В переменной $myData хранится значение {{$myData}}</li>
        <li>Выведем результат выполнения функции сразу в шаблон @date(-777856620, "2006-01-02 15:04") или так {{date(-777856620, "2006-01-02 15:04")}}</li>
    </ul>

    <h2>Создание map[string]interface{}</h2>
    @exec($map = [
        "key1" => "value 1",
        "key2" => "value 2",
        "key3" => "value 3",
        "key4" => "value 4",
        "key5" => "value 5",
        "key6" => "value 6",
        "key7" => "value 7",
        "my_data" => $myData
    ])

    <ul>
        <li>$map.key1 = {{$map.key1}}</li>
        <li>$map.key2 = {{$map.key2}}</li>
        <li>$map.key3 = {{$map.key3}}</li>
        <li>$map.key4 = {{$map.key4}}</li>
        <li>$map.key5 = {{$map.key5}}</li>
        <li>$map.key6 = {{$map.key6}}</li>
        <li>$map.key7 = {{$map.key7}}</li>
        <li>$map.my_data = {{$map.my_data}}</li>
    </ul>

    <h2>Комментарии</h2>
    {* Этой строки не будет видно *}
    {{-- И этой тоже --}}

    <h2>Подключение шаблонов</h2>
    @include("layouts/include")

    <h2>Вывод кода шаблона без обработки</h2>
<pre>
    {{"
    @foreach($numbers as $number1)
        @foreach($numbers as $number2)
            <li>{{$number1}} + {{$number2}} = {{$number1 + $number2}}</li>
            <li>{{$number1}} - {{$number2}} = {{$number1 - $number2}}</li>
            <li>{{$number1}} * {{$number2}} = {{$number1 * $number2}}</li>
            <li>{{$number1}} / {{$number2}} = {{$number1 / $number2}}</li>
            <li>{{$number1}} % {{$number2}} = {{$number1 % $number2}}</li>
        @endforeach
    @endforeach
    "}}
</pre>


@endsection

