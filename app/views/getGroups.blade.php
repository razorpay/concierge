@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
	        <h1>Active Leases</h1>
	        <table class="table-bordered table-striped">
		    <thead>
		        <tr>
		        <th>Created By:</th>
		        <th>Leased Ip:</th>
		        <th>On Group:</th>
		        <th>Protocol</th>
		        <th>Port(s)</th>
		        <th>Time Left</th>
		        <th>Terminate?</th>
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
	    			<td>
	    			<form method="post" action="/manage/{{{$lease->group_id}}}">
	    			<input type="hidden" name="lease_id" value="{{{$lease->id}}}" />
	    			<button type="submit" onclick="return confirm('Are you sure you want to terminate this lease?');">
	    			<span title="Terminate Lease" class="glyphicon glyphicon-minus-sign"></span>
	    			</button>
	    			</form>
	    			</td>
		       	</tr>
		       	@endforeach
		       	@if(!$leases->count())
		       	<tr><td colspan="7">No Active Leases</td></tr>
		       	@endif
		    </tbody>
   		    </table>
   		    <br/>
 			<h1>Security Groups</h1>
 			<table class="table-bordered table-striped">
		    <thead>
		        <tr>
		        <th>Group Name:</th>
		        <th>Group Id:</th>
		        <th>Group Description:</th>
		        <th>VPC:</th>
		        <th>Name Tag:</th>
		        </tr>
		    </thead>
		    <tbody>
			@foreach($security_groups as $security_group)
			<tr>
			<td><a href="/manage/{{{$security_group['GroupId']}}}">{{{$security_group['GroupName']}}}</a></td>
			<td>{{{$security_group['GroupId']}}}</td>
			<td>{{{$security_group['Description']}}}</td>
			<td>{{{$security_group['VpcId']}}}</td>
			<td>{{{$security_group['Tags']['0']['Value']}}}</td>
			</tr>
			@endforeach
			</tbody>
			</table>
		</div>
	</div>
@stop