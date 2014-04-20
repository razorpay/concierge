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
        <h3>Secure Access Lease Created</h3>
        <p>A new lease has been created as follows:</p>
        <table class="table-bordered table-striped">
            <thead>
                <tr>
                <th>Created By:</th>
                <th>Leased Ip:</th>
                <th>On Group:</th>
                <th>Protocol</th>
                <th>Port(s)</th>
                <th>Time Left</th>
                </tr>
            </thead>
            <tbody>
               <tr>
               <td>{{{User::find($lease['user_id'])->username}}}</td>
               <td>{{{$lease['lease_ip']}}}</td>
               <td>{{{$lease['group_id']}}}</td>
               <td>{{{$lease['protocol']}}}</td>
               <td>{{{$lease['port_from']}}}-{{{$lease['port_to']}}}</td>
               <td>
               <?php
               	$time_left=strtotime($lease['created_at'])+$lease['expiry']-time(); 
        	    $hours=intval(floor($time_left/3600)); 
        	    $minutes=intval(floor(($time_left-$hours*3600)/60));
        	    echo "$hours hours $minutes minutes";
        	   ?>
        	   </td>
               </tr>
            </tbody>
        </table>
    </body>
</html>