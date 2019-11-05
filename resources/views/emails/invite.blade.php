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
        <link rel="stylesheet" href="https://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">
        <link rel="stylesheet" href="assets/css/style.css">
        <!-- HTML5 shim, for IE6-8 support of HTML5 elements -->
        <!--[if lt IE 9]>
            <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
        <![endif]-->
        @show
    </head>

    <body>
            <h3>Secure Access Lease Invited</h3>
            <p>You have been invited for a secure access lease as follows:</p>
            <table class="table table-hover table-bordered">
                <thead>
                    <tr>
                    <th>Creator</th>
                    <th>Protocol</th>
                    <th>Port(s)</th>
                    <th>Expiry</th>
                    </tr>
                </thead>
                <tbody>
                   <tr>
                   <td>{{{User::find($invite['user_id'])->username}}}({{{User::find($invite['user_id'])->name}}})</td>
                   <td>{{{$invite['protocol']}}}</td>
                   <td>{{{$invite['port_from']}}}-{{{$invite['port_to']}}}</td>
                   <td>
                   <?php
                    //Calculating time to expiry in hours & minutes
                    $hours=intval(floor($invite['expiry']/3600));
            	      $minutes=intval(floor(($invite['expiry']-$hours*3600)/60));
            	      echo "$hours hours $minutes minutes";
            	     ?>
            	     </td>
                   </tr>
                </tbody>
            </table>
            <h4>Click <a href="{{url('/invite')}}/{{$invite['token']}}" target="_blank">here</a> to access the lease.</h4>
            <p>
            Alternatively, you may open the following link in your browser:<br/>
            {{url('/invite')}}/{{$invite['token']}}
            </p>
            <p>Information: The lease will only be valid on the security groups chosen by the creator.</p>
    </body>
</html>
