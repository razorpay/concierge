@extends('layouts.master')
@section('headincludes')
	 @parent
	 <script typpe="text/javascript">
	 function displayform()
	 {
	 	document.getElementById('custom_form').style.visibility="visible";
	 }
	 </script>
@stop
@section('content')
	<div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
          	Name: {{$security_group['GroupName']}}<br/>
			Id: {{$security_group['GroupId']}}<br/>
			Description: {{$security_group['Description']}}<br/>
			Active Leases:
			<table class="table-bordered table-striped">
	        <thead>
	          <tr>
	            <th>Created By:</th>
	            <th>Leased Ip:</th>
	            <th>Protocol</th>
	            <th>Port(s)</th>
	            <th>Time Left</th>
	          </tr>
	        </thead>
	        <tbody>
	        	@foreach($leases as $lease)
	        	<tr>
	        		<td>{{{User::find($lease->user_id)->username}}}</td>
	        		<td>{{{$lease->lease_ip}}}</td>
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

			Security Group Rules:<br/>
            Inbound Rules: 
	        <table class="table-bordered table-striped">
	        <thead>
	          <tr>
	            <th>Protocol</th>
	            <th>Port</th>
	            <th>Source:</th>
	          </tr>
	        </thead>
	        <tbody>
			@foreach($security_group['IpPermissions'] as $rule)
			   @foreach($rule['UserIdGroupPairs'] as $rule_group)
			    <tr>
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
		    <table class="table-bordered table-striped">
	        <thead>
	          <tr>
	            <th>Protocol</th>
	            <th>Port</th>
	            <th>Destination:</th>
	          </tr>
	        </thead>
	        <tbody>
			@foreach($security_group['IpPermissionsEgress'] as $rule)
			    @foreach($rule['UserIdGroupPairs'] as $rule_group)
			    <tr>
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

		    <form action="" method="POST"><input type="hidden" name="rule_type" value="ssh" /><input type="submit" value="Get SSH Access on this Group" /></form>
		    <form action="" method="POST"><input type="hidden" name="rule_type" value="https" /><input type="submit" value="Get HTTPS Access on this Group" /></form>
		    <button onclick="javascript: displayform()">Get Custom Access</button>
		    <form id="custom_form" style="visibility:hidden" action="" method="POST">
		    <input type="hidden" name="rule_type" value="custom" />
		    Protocol: <input type="text" name="protocol" value="" placeholder="tcp/udp"/><br/>
		    From Port:<input type="text" name="port_from" placeholder="From Port:"/>
		    To Port:<input type="text" name="port_to" placeholder="To port:"/> <br/> Keep From & To Port samefor single port access.<br/>
		    <input type="submit" value="Get Access" />
		    </form>
		</div>    
	</div>
@stop