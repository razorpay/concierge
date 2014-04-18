<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>AWS Security Group Manager</title>
	<style>
		@import url(//fonts.googleapis.com/css?family=Lato:700);

		body {
			margin:0;
			font-family:'Lato', sans-serif;
			text-align:center;
			color: #999;
		}

		.content {
			width: 300px;
			height: 200px;
			position: absolute;
			left: 50%;
			top: 50%;
			margin-left: -150px;
			margin-top: -100px;
		}

		a, a:visited {
			text-decoration:none;
		}

		h1 {
			font-size: 32px;
			margin: 16px 0 0 0;
		}
	</style>
</head>
<body>
<a href="{{action('HomeController@getLogout')}}">Logout</a>
	<div class="content">
		@foreach($security_groups as $security_group)
        <a href="/manage/{{{$security_group['GroupId']}}}">{{{$security_group['GroupName']}}}</a><br/>{{{$security_group['Description']}}}<br/>
		@endforeach
	</div>
</body>
</html>
