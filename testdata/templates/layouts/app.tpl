<!DOCTYPE html>
<html lang="ru">
<head>
    <title>{{$page_title}}</title>
</head>
<body>
    @yield("content")
    @yield("section1")
    @yield("section2")
    @yield("section3")
</body>
</html>

@section("section3")
    <div>Секция из файла app.tpl</div>
@endsection