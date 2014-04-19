@extends('layouts.master')

@section('content')
	<div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
          	Name: {{$security_group['GroupName']}}<br/>
			Id: {{$security_group['GroupId']}}<br/>
			Description: {{$security_group['Description']}}<br/>
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
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
					<td>Security Group: <a href="/manage/{{$rule_group['GroupId']}}">{{$rule_group['GroupId']}}</a></td>
			    </tr>
				@endforeach
			    @foreach($rule['IpRanges'] as $rule_ip)
			    <tr>
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
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
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
					<td>Security Group: <a href="/manage/{{$rule_group['GroupId']}}">{{$rule_group['GroupId']}}</a></td>
				</tr>
				@endforeach
			    @foreach($rule['IpRanges'] as $rule_ip)
			    <tr>
					<td>{{$rule['IpProtocol']}}</td>
					<td>{{$rule['FromPort']}}-{{$rule['ToPort']}}</td>
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
		    <input type="text" name="protocol" value="" placeholder="tcp/udp"/>
		    <input type="text" name="port_from" placeholder="From Port:"/>
		    <input type="text" name="port_to" placeholder="To port:"/>
		    <input type="submit" value="Get Access" />
		    </form>
		</div>    
	</div>
@stop