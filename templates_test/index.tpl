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
            @if ($strings[$i] != "")
                В переменной $strings, под индексом {{$i}} содержится значение {{$strings[$i]}}
            @else
                В переменной $strings, под индексом {{$i}} нет значения
            @endif
        </li>
    @endfor
    </ul>

    <h2>Цикл foreach</h2>
    <ul>
    @foreach ($strings as $i => $value)
        <li>$strings[{{$i}}] == {{$value}}</li>
    @endforeach
    </ul>

    <h2>Условия</h2>
    <ul>
    @foreach($numbers as $number1)
        @foreach($numbers as $number2)
            @if ($number1 < $number2)
                <li>{{$number1}} < {{$number2}}</li>
            @endif
            @if ($number1 <= $number2)
                <li>{{$number1}} <= {{$number2}}</li>
            @endif
            @if ($number1 == $number2)
                <li>{{$number1}} == {{$number2}}</li>
            @endif
            @if ($number1 >= $number2)
                <li>{{$number1}} >= {{$number2}}</li>
            @endif
            @if ($number1 > $number2)
                <li>{{$number1}} > {{$number2}}</li>
            @endif
            @if ($number1 != $number2)
                <li>{{$number1}} != {{$number2}}</li>
            @endif
            @if ($number1 <> $number2)
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
    Value: {{$struct["Value"]}}
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

    <h2>Комментарии</h2>
    {* Этой строки не будет видно *}
    {{-- И этой тоже --}}

    <h2>Подключение шаблонов</h2>
    @include("layouts/include")

    <h1>Вывод кода шаблона без обработки</h2>
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

