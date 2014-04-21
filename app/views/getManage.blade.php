@extends('layouts.master')
@section('headincludes')
	 @parent
	 <script typpe="text/javascript">
	 function displayform(form_id)
	 {
	 	document.getElementById(form_id).style.visibility="visible";
	 }
	 </script>
@stop
@section('content')
	<div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
          	Name: {{{$security_group['GroupName']}}}<br/>
			Id: {{{$security_group['GroupId']}}}<br/>
			Description: {{{$security_group['Description']}}}<br/>
			VPC-Id: {{{$security_group['VpcId']}}} <br/>
			Name Tag: 
			@if(isset($security_group['Tags']['0']['Value']))
			{{{$security_group['Tags']['0']['Value']}}}
			@endif
			<br/>
			<h2>Active Leases:</h2>
			<table class="table table-hover table-bordered">
	        <thead>
	          <tr>
	            <th>Creator</th>
	            <th>Leased IP</th>
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
	        		<td>{{{$lease->protocol}}}</td>
	        		<td>{{{$lease->port_from}}}-{{{$lease->port_to}}}</td>
	        		<td>
	        		<?php
	        		    //Calculating Time to Expiry in Hours and minutes
	        			$time_left=strtotime($lease->created_at)+$lease->expiry-time(); 
    					$hours=intval(floor($time_left/3600)); 
    					$minutes=intval(floor(($time_left-$hours*3600)/60));
    					echo "$hours hours $minutes minutes";
    				?>
    				</td>
    				<td><form method="post" action="">
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
		       	<tr><td colspan="6" style="text-align:center">No Active Leases</td></tr>
		       	@endif
	        </tbody>
	        </table>

			<h2>Security Group Rules:</h2>
            Inbound Rules: 
	        <table class="table table-hover table-bordered">
	        <thead>
	          <tr>
	            <th>Protocol</th>
	            <th>Port</th>
	            <th>Source</th>
	          </tr>
	        </thead>
	        <tbody>
			@foreach($security_group['IpPermissions'] as $rule)
			   @foreach($rule['UserIdGroupPairs'] as $rule_group)
			    <tr>
			     	{{-- Checking for all traffic rule --}}
					@if("-1"!=$rule['IpProtocol'])
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
					@else
					<td>All</td>
					<td>All</td>
					@endif
					<td>Security Group: <a href="/manage/{{$rule_group['GroupId']}}">{{$rule_group['GroupId']}}</a></td>
			    </tr>
				@endforeach
			    @foreach($rule['IpRanges'] as $rule_ip)
			    <tr>
			    	{{-- Checking for all traffic rule --}}
			    	@if("-1"!=$rule['IpProtocol'])
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
					@else
					<td>All</td>
					<td>All</td>
					@endif
					<td>CIDR IP(s): {{$rule_ip['CidrIp']}}</td>
			    </tr>
				@endforeach
		    @endforeach 
		    </tbody>
		    </table>
		    Outbound Rules: <br/>
		    <table class="table table-hover table-bordered">
	        <thead>
	          <tr>
	            <th>Protocol</th>
	            <th>Port</th>
	            <th>Destination</th>
	          </tr>
	        </thead>
	        <tbody>
			@foreach($security_group['IpPermissionsEgress'] as $rule)
			    @foreach($rule['UserIdGroupPairs'] as $rule_group)
			    <tr>
			    	{{-- Checking for all traffic rule --}}
					@if("-1"!=$rule['IpProtocol'])
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
					@else
					<td>All</td>
					<td>All</td>
					@endif
					<td>Security Group: <a href="/manage/{{$rule_group['GroupId']}}">{{$rule_group['GroupId']}}</a></td>
				</tr>
				@endforeach
			    @foreach($rule['IpRanges'] as $rule_ip)
			    <tr>
			    	{{-- Checking for all traffic rule --}}
					@if("-1"!=$rule['IpProtocol'])
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
					@else
					<td>All</td>
					<td>All</td>
					@endif
					<td>CIDR IP(s): {{$rule_ip['CidrIp']}}</td>
				</tr>
				@endforeach
		    @endforeach
		    </tbody>
		    </table>
		    <br/>

		    <button onclick="javascript: displayform('ssh_form')">Get SSH Access on this Group</button><br/>
		    <button onclick="javascript: displayform('https_form')">Get HTTPS Access on this Group</button><br/>
		    <button onclick="javascript: displayform('custom_form')">Get Custom Access on this Group</button><br/>
		    
		    <form id="ssh_form" style="visibility:hidden" action="" method="POST">
			    <input type="hidden" name="rule_type" value="ssh" />
			    <input type="hidden" name="_token" value="{{{csrf_token()}}}">
			   SSH Access Expiry: <select name="expiry" required>
				  <option value="3600" selected>1 hour</option>
				  <option value="14400">4 hours</option>
				  <option value="43200">12 hours</option>
				  <option value="86400">1 Day</option>
				</select>
			    <input type="submit" value="Get Access" />
		    </form>

		    <form id="https_form" style="visibility:hidden" action="" method="POST">
			    <input type="hidden" name="rule_type" value="https" />
			    <input type="hidden" name="_token" value="{{{csrf_token()}}}">
			    HTTPS Access Expiry: <select name="expiry" required>
				  <option value="3600" selected>1 hour</option>
				  <option value="14400">4 hours</option>
				  <option value="43200">12 hours</option>
				  <option value="86400">1 Day</option>
				</select>
			    <input type="submit" value="Get Access" />
		    </form>
		    
		    <form id="custom_form" style="visibility:hidden" action="" method="POST">
			    <input type="hidden" name="rule_type" value="custom" />
			    <input type="hidden" name="_token" value="{{{csrf_token()}}}">
			    Protocol: <input type="text" name="protocol" value="" placeholder="tcp/udp" required/><br/>
			    From Port:<input type="text" name="port_from" placeholder="From Port:" required/>
			    To Port:<input type="text" name="port_to" placeholder="To port:" required/> <br/> Keep From & To Port same for single port access.<br/>
			    Custom Access Expiry: <select name="expiry" required>
				  <option value="3600" selected>1 hour</option>
				  <option value="14400">4 hours</option>
				  <option value="43200">12 hours</option>
				  <option value="86400">1 Day</option>
				</select>
			    <input type="submit" value="Get Access" />
		    </form>
		</div>    
	</div>
@stop