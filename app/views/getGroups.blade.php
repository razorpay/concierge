@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
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
		     	@foreach($leases as $lease)
		      	<tr>
		       		<td>{{{$lease->user->username}}}</td>
		       		<td>{{{$lease->lease_ip}}}</td>
		       		<td><a href="/manage/{{{$lease->group_id}}}">{{{$lease->group_id}}}</a></td>
		       		<td>{{{$lease->protocol}}}</td>
		       		<td>{{{$lease->port_from}}}-{{{$lease->port_to}}}</td>
		       		<td>
		       		<?php
		       			$time_left=strtotime($lease->created_at)+$lease->expiry-time(); 
	    				$hours=intval(floor($time_left/3600)); 
	    				$minutes=intval(floor(($time_left-$hours*3600)/60));
	    				echo "$hours hours $minutes minutes";
	    			?>
	    			</td>	
		       	</tr>
		       	@endforeach
		    </tbody>
   		    </table>
   		    <br/>
 
			@foreach($security_groups as $security_group)
	        <a href="/manage/{{{$security_group['GroupId']}}}">{{{$security_group['GroupName']}}} - {{{$security_group['GroupId']}}}</a><br/>{{{$security_group['Description']}}}<br/>
			@endforeach
		</div>
	</div>
@stop