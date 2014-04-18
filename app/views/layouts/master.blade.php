<!DOCTYPE html>
<html lang="en">

    <head>
        <meta charset="utf-8" />
        <title>
            @section('title')
                Laravel Duo Security
            @show
        </title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0 minimum-scale=1">
        <meta http-equiv='X-UA-Compatible' content='IE=edge,chrome=1'>
        {{ HTML::style('//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css') }}
        {{ HTML::style('assets/css/style.css') }}
        <!-- HTML5 shim, for IE6-8 support of HTML5 elements -->
        <!--[if lt IE 9]>
            <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
        <![endif]-->
    </head>

    <body>

        @if(Session::has('message'))
        <div id="sessionMessage" class="alert alert-warning fade in">
            <button type="button" class="close" data-dismiss="alert" aria-hidden="true">x</button>
            <p class="center">{{ Session::get('message') }}</p>
        </div>
        @endif

        @if(Auth::check())
            <div class="container">
                <p class="center">
                    <a href="{{ URL::to('/logout') }}">Logout</a>
                </p>
            </div>
        @else
            <div class="container">
                @yield('content')
            </div>
        @endif

        @section('footer_scripts')
            {{ HTML::script('//code.jquery.com/jquery-1.10.2.min.js') }}
            {{ HTML::script('//netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js') }}
        @show

    </body>

</html>