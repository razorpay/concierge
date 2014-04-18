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
	<script type="text/javascript">
	function displayform(){
	 document.getElementById("custom_form").style.visibility="visible";
	}
	</script>
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
</body>
</html>
