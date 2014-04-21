@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
	        <h1>Change Password</h1>
	        <form action="" method="POST">
		     	Old Password: <input type="password" name="old_password" placeholder="Old Password" required /><br/>
		     	New Password: <input type="password" name="password" placeholder="New Password" required /><br/>
		     	Confirm New Password: <input type="password" name="password_confirmation" placeholder="New Password" required /><br/>
		     	<input type="hidden" name="_token" value="{{{csrf_token()}}}">
		     	<input type="submit" value="Change Password" />
		    </form>
	     </div>
	</div>
@stop