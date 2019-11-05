@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
            @if(isset($user->id))
                <h1>Edit User</h1>
            @else
                <h1>Add New User</h1>
            @endif
            <form class="form-horizontal" role="form" action="" method="POST">

            <div class="form-group{{ $errors->has('email') ? ' has-error' : '' }}">
                <label for="email" class="col-sm-4 control-label">Email Id</label>
                <div class="col-sm-6">
                  <input type="text" class="form-control" name="email" placeholder="Email" required value="{{$user->email}}">
                </div>

                @if ($errors->has('email'))
                <div class="col-sm-offset-4 col-sm-6">
                    <span class="help-block">
                        <strong>{{ $errors->first('email') }}</strong>
                    </span>
                </div>
                @endif
            </div>

            <div class="form-group{{ $errors->has('name') ? ' has-error' : '' }}">
                <label for="Full Name" class="col-sm-4 control-label">Full Name</label>
                <div class="col-sm-6">
                  <input type="text" class="form-control" name="name" placeholder="Name" required value="{{$user->name}}">
                </div>

                @if ($errors->has('name'))
                <div class="col-sm-offset-4 col-sm-6">
                    <span class="help-block">
                        <strong>{{ $errors->first('name') }}</strong>
                    </span>
                </div>
                @endif
            </div>

              <div class="form-group">
                <label for="admin" class="col-sm-4 control-label">Role</label>
                <div class="col-sm-6">
                <label class="radio-inline">
                  <input type="radio" name="admin" value="0" {{ $user->admin ? '' : 'checked' }}>Standard User
                </label>
                <label class="radio-inline">
                  <input type="radio" name="admin" value="1" {{ $user->admin ? 'checked' : '' }}>Admin User
                </label>
                </div>
              </div>
              <input type="hidden" name="_token" value="{{{csrf_token()}}}">
              <div class="form-group">
                <div class="col-sm-offset-4 col-sm-6">
                  <button type="submit" class="btn btn-default">Save</button>
                </div>
              </div>
            </form>
            <h4>Notes:</h4>
            <ul>
            <li>Admin Users can add/remove other users (including other admin users)</li>
            </ul>
         </div>
    </div>
@stop
