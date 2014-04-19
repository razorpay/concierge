@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
			@foreach($security_groups as $security_group)
	        <a href="/manage/{{{$security_group['GroupId']}}}">{{{$security_group['GroupName']}}}</a><br/>{{{$security_group['Description']}}}<br/>
			@endforeach
		</div>
	</div>
@stop