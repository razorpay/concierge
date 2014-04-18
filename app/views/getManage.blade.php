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
			color: #999;
		}

		.content {
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
	<div class="content">
	Name: {{$security_group['GroupName']}}<br/>
	Id: {{$security_group['GroupId']}}<br/>
	Description: {{$security_group['Description']}}<br/>
	Inbound Rules: <br/>
	Internal Rules: <br/>
	@foreach($security_group['IpPermissions'] as $rule)
	    @foreach($rule['UserIdGroupPairs'] as $rule_group)
			Protocol: {{$rule['IpProtocol']}}<br/>
			Port(s): {{$rule['FromPort']}}-{{$rule['ToPort']}}<br/>
			From Security Group: {{$rule_group['GroupId']}}<br/>
		@endforeach
	@endforeach
	<br/>
	External Rules: <br/>
	@foreach($security_group['IpPermissions'] as $rule)
	    @foreach($rule['IpRanges'] as $rule_ip)
			Protocol: {{$rule['IpProtocol']}}<br/>
			Port(s): {{$rule['FromPort']}}-{{$rule['ToPort']}}<br/>
			From CIDR IP(s): {{$rule_ip['CidrIp']}}<br/>
		@endforeach
    @endforeach 
    <br/>
    Outbound Rules: <br/>
    Internal Rules: <br/>
	@foreach($security_group['IpPermissionsEgress'] as $rule)
	    @foreach($rule['UserIdGroupPairs'] as $rule_group)
			Protocol: {{$rule['IpProtocol']}}<br/>
			Port(s): {{$rule['FromPort']}}-{{$rule['ToPort']}}<br/>
			To Security Group: {{$rule_group['GroupId']}}<br/>
		@endforeach
	@endforeach
	<br/>
	External Rules: <br/>
	@foreach($security_group['IpPermissionsEgress'] as $rule)
	    @foreach($rule['IpRanges'] as $rule_ip)
			Protocol: {{$rule['IpProtocol']}}<br/>
			Port(s): {{$rule['FromPort']}}-{{$rule['ToPort']}}<br/>
			To CIDR IP(s): {{$rule_ip['CidrIp']}}<br/>
		@endforeach
    @endforeach
    <br/>
	</div>
</body>
</html>
