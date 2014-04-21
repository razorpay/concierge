<!DOCTYPE html>
<html lang="en">

    <head>
        <meta charset="utf-8" />
        <title>
            @section('title')
                Aws Access Manager
            @show
        </title>
        @section('headincludes')
        <meta name="viewport" content="width=device-width, initial-scale=1.0 minimum-scale=1">
        <meta http-equiv='X-UA-Compatible' content='IE=edge,chrome=1'>
        {{ HTML::style('//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css') }}
        {{ HTML::style('assets/css/style.css') }}
        <!-- HTML5 shim, for IE6-8 support of HTML5 elements -->
        <!--[if lt IE 9]>
            <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
        <![endif]-->
        @show
    </head>

    <body>

        @if(Session::has('message'))
        <div id="sessionMessage" class="alert alert-warning fade in">
            <button type="button" class="close" data-dismiss="alert" aria-hidden="true">x</button>
            <p class="center">{{ Session::get('message') }}</p>
        </div>
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
                  <a class="navbar-brand" href="{{ URL::to('/')}}">Home</a>
                </div>

                <!-- Collect the nav links, forms, and other content for toggling -->
                <div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
                  <ul class="nav navbar-nav navbar-right">
                    <li><p class="navbar-text">Signed in as {{{Auth::user()->username}}}</p></li>
                    <li><a href="{{ URL::to('/password')}}"> Change Password</a></li>
                    <li><a href="{{ URL::to('/logout') }}">Logout</a></li>
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
            {{ HTML::script('//code.jquery.com/jquery-1.10.2.min.js') }}
            {{ HTML::script('//netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js') }}
            <script type="text/javascript">
            setTimeout(function() {
                $('#sessionMessage').fadeOut('slow');
            }, 4000);
            </script>
        @show

    </body>

</html>