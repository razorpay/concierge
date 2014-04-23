@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
	        <h2>Active Leases</h2>
	        <table class="table table-hover table-bordered">
		    <thead>
		        <tr>
		        <th>Creator</th>
		        <th>Leased IP</th>
		        <th>Security Group</th>
		        <th>Protocol</th>
		        <th>Port(s)</th>
		        <th>Time Left</th>
		        <th>Type</th>
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
		       		<div class="time" id="{{{$lease->id}}}">
	        		{{{strtotime($lease->created_at)+$lease->expiry-time()}}}
    				</div>
	    			</td>
	    			<td>
	    			@if($lease->invite_email)
	    			 	@if("URL"==$lease->invite_email)
	    			 		URL Invite
	    			 	@else
	    			 		Email Invite: {{{$lease->invite_email}}}
	    			 	@endif
	    			@else
	    					Self Access
	    			@endif
	    			</td>
	    			<td>
		    			<form method="post" action="/manage/{{{$lease->group_id}}}/terminate">
		    			<input type="hidden" name="lease_id" value="{{{$lease->id}}}" />
		    			<input type="hidden" name="_token" value="{{{csrf_token()}}}">
		    			<a href="" style="color: #ff0000;" onclick="if(confirm('Are you sure you want to terminate this lease?')) {parentNode.submit();} return false;">
		    			<span title="Terminate Lease" class="glyphicon glyphicon-minus-sign"></span>
		    			</a>
		    			</form>
	    			</td>
		       	</tr>
		       	@endforeach

		       	@if(!$leases->count())
		       	<tr><td colspan="8" style="text-align:center">No Active Leases</td></tr>
		       	@endif
		    </tbody>
   		    </table>

   		    <h2>Active Invites</h2>
			<table class="table table-hover table-bordered">
	        <thead>
	          <tr>
	            <th>Creator</th>
	            <th>Security Group</th>
	            <th>Protocol</th>
	            <th>Port(s)</th>
	            <th>Expiry</th>
	            <th>Type</th>
	            <th>Terminate?</th>
	          </tr>
	        </thead>
	        <tbody>
	        	@foreach($invites as $invite)
	        	<tr>
	        		<td>{{{$invite->user->username}}}</td>
	        		<td><a href="/manage/{{{$invite->group_id}}}">{{{$invite->group_id}}}</a></td>
	        		<td>{{{$invite->protocol}}}</td>
	        		<td>{{{$invite->port_from}}}-{{{$invite->port_to}}}</td>
	        		<td>
	        		<?php
	        		    //Calculating Time to Expiry in Hours and minutes
    					$hours=intval(floor($invite->expiry/3600)); 
    					$minutes=intval(floor(($invite->expiry-$hours*3600)/60));
    					echo "$hours hours $minutes minutes";
    				?>
    				</td>
    				<td>
    				@if($invite->email)
    					Email: {{{$invite->email}}}
    				@else
    					URL Invite
    				@endif
    				</td>
    				<td>
	    				<form method="post" action="/manage/{{{$invite->group_id}}}/terminate">
	    				<input type="hidden" name="invite_id" value="{{{$invite->id}}}" />
	    				<input type="hidden" name="_token" value="{{{csrf_token()}}}">
	    				<a href="" style="color: #ff0000;" onclick="if(confirm('Are you sure you want to terminate this invite?')) {parentNode.submit();} return false;">
	    					<span title="Terminate Invite" class="glyphicon glyphicon-minus-sign"></span>
	    				</a>
	    				</form>
    				</td>	
	        	</tr>
	        	@endforeach
	        	@if(!$invites->count())
		       	<tr><td colspan="7" style="text-align:center">No Active Invites</td></tr>
		       	@endif
	        </tbody>
	        </table>

 			<h2>Security Groups</h2>
 			<table class="table table-hover table-bordered">
		    <thead>
		        <tr>
		        <th>Name</th>
		        <th>ID</th>
		        <th>Description</th>
		        <th>VPC</th>
		        <th>Name Tag</th>
		        </tr>
		    </thead>
		    <tbody>
			@foreach($security_groups as $security_group)
			<tr>
			<td><a href="/manage/{{{$security_group['GroupId']}}}">{{{$security_group['GroupName']}}}</a></td>
			<td>{{{$security_group['GroupId']}}}</td>
			<td>{{{$security_group['Description']}}}</td>
			<td>{{{$security_group['VpcId']}}}</td>
			<td>
			{{-- Display The Name Tag if it exists --}}
			@if(isset($security_group['Tags']['0']['Value']))
			{{{$security_group['Tags']['0']['Value']}}}
			@endif
			</td>
			</tr>
			@endforeach
			</tbody>
			</table>
		</div>
	</div>
@stop
@section('footer_scripts')
@parent
<script type="text/javascript">
   function component(x, v) {
        return Math.floor(x / v);
    }
    $('.time').each(function(i, obj) {   
	    var div=$(this)                   
	    var timestamp = div.text()
	    
	    if(timestamp>0)
	    {
		    setInterval(function() { // execute code each second

		        timestamp--; // decrement timestamp with one second each second

		        var hours   = component(timestamp,      60 * 60), // hours
		            minutes = component(timestamp,           60) % 60, // minutes
		            seconds = component(timestamp,            1) % 60; // seconds

		        div.text(hours + " hours " + minutes + " minutes " + seconds + " seconds"); // display
		        //alert($(this).text());

		    }, 1000); // interval each second = 1000 ms
	    }
	    else
	    {
	    	div.text("Expired");
	    }
	});
</script>
@stop