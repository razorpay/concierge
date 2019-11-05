<!DOCTYPE html>
<html lang="en">

    <head>
        <meta charset="utf-8" />
        <title>
            @section('title')
                Concierge - AWS Access Control
            @show
        </title>
        @section('headincludes')
        <link rel="shortcut icon" href="{{asset('favicon.png')}}">
        <meta name="viewport" content="width=device-width, initial-scale=1.0 minimum-scale=1">
        <meta http-equiv='X-UA-Compatible' content='IE=edge,chrome=1'>
        <link rel="stylesheet" href="https://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">
        <link rel="stylesheet" href="assets/css/style.css">
        <!-- HTML5 shim, for IE6-8 support of HTML5 elements -->
        <!--[if lt IE 9]>
            <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
        <![endif]-->
        @show
        <link href="{{ asset('assets/css/flash-message.css') }}" rel="stylesheet">
    </head>

    <body>
        @if(Session::has('message'))
        <div class="custom-flash {{ Session::get('class') }} ">{{ Session::get('message') }}</div>
        @endif

        @if(Auth::check())
            <nav class="navbar navbar-inverse" role="navigation">
              <div class="container-fluid">
                <!-- Brand and toggle get grouped for better mobile display -->
                <div class="navbar-header">
                  <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
                    <span class="sr-only">Toggle navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                  </button>
                  <a class="navbar-brand" href="/">Concierge Home</a>
                </div>

                <!-- Collect the nav links, forms, and other content for toggling -->
                <div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
                  <ul class="nav navbar-nav navbar-right">
                    <li><p class="navbar-text">Signed in as {{{Auth::user()->name}}}</p></li>
                    @if(Auth::user()->admin)
                    <li><a href="/users">Manage Users</a></li>
                    @endif
                    <li><a href="/logout">Logout</a></li>
                  </ul>
                </div><!-- /.navbar-collapse -->
              </div><!-- /.container-fluid -->
            </nav>
            <div class="container">
                @yield('content')
            </div>
        @else
            <div class="container">
                @yield('content')
            </div>
        @endif

        @section('footer_scripts')
            <script src='https://code.jquery.com/jquery-2.1.0.min.js'></script>
            <script type="text/javascript" src="{{ asset('assets/js/flash-message.js') }}"></script>
            <script src='https://netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js'></script>
            <script type="text/javascript">
            setTimeout(function() {
                $('#sessionMessage').fadeOut('slow');
            }, 4000);
            </script>
        @show

    </body>

</html>
