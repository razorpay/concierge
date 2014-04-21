@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
	        <h1>Change Password</h1>
	        <form class="form-horizontal" role="form" action="" method="POST">
			  <div class="form-group">
			    <label for="old_password" class="col-sm-4 control-label">Old Password</label>
			    <div class="col-sm-6">
			      <input type="password" class="form-control" name="old_password" placeholder="Old Password" required>
			    </div>
			  </div>
			  <div class="form-group">
			    <label for="password" class="col-sm-4 control-label">New Password</label>
			    <div class="col-sm-6">
			      <input type="password" class="form-control" name="password" placeholder="New Password" required>
			    </div>
			  </div>
			  <div class="form-group">
			    <label for="password_confirmation" class="col-sm-4 control-label">Confirm New Password</label>
			    <div class="col-sm-6">
			      <input type="password" class="form-control" name="password_confirmation" placeholder="New Password" required>
			    </div>
			  </div>
			  <input type="hidden" name="_token" value="{{{csrf_token()}}}">
			  <div class="form-group">
			    <div class="col-sm-offset-4 col-sm-6">
			      <button type="submit" class="btn btn-default">Change Password</button>
			    </div>
			  </div>
			</form>
	     </div>
	</div>
@stop