@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
	        <h2>Site Users</h2>
	        <table class="table table-hover table-bordered">
		    <thead>
		        <tr>
		        <th>Username</th>
		        <th>Name</th>
		        <th>Role</th>
		        <th>Active Self Leases</th>
		        <th>Active URL Leases</th>
		        <th>Active Email Leases</th>
		        <th>Total Leases Created <span title="Includes Terminated & Expired Leases" class="glyphicon glyphicon-question-sign"></span></th>
		        <th>Active URL Invites</th>
		        <th>Active Email Invites</th>
		        <th>Delete?</th>
		        </tr>
		    </thead>
		    <tbody>
		     	@foreach($users as $user)
		      	<tr>
		       		<td>{{{$user->username}}}</td>
		       		<td>{{{$user->name}}}</td>
		       		<td>
		       		@if($user->admin)
		       		Admin
		       		@else
		       		Standard User
		       		@endif
		       		</td>
		       		<td>{{{$user->getActiveLeases("self")->count()}}}</td>
		       		<td>{{{$user->getActiveLeases("url")->count()}}}</td>
		       		<td>{{{$user->getActiveLeases("email")->count()}}}</td>
		       		<td>{{{$user->leases()->withTrashed()->count()}}}</td>
		       		<td>{{{$user->getActiveInvites("url")->count()}}}</td>
		       		<td>{{{$user->getActiveInvites("email")->count()}}}</td>
	    			<td>
	    				@if($user->id != Auth::user()->id)
		    			<form method="post" action="">
		    			<input type="hidden" name="user_id" value="{{{$user->id}}}" />
		    			<input type="hidden" name="_token" value="{{{csrf_token()}}}">
		    			<a href="" style="color: #ff0000;" onclick="if(confirm('Are you sure you want to delete this user (All his active leases/invites will be terminated?')) {parentNode.submit();} return false;">
		    			<span title="Delete User" class="glyphicon glyphicon-minus-sign"></span>
		    			</a>
		    			</form>
		    			@endif
	    			</td>
		       	</tr>
		       	@endforeach
		    </tbody>
   		    </table>
		</div>
	</div>
@stop