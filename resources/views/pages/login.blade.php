@extends('layouts.master')

@section('content')

    <div class="container">

        <div class="row">

            <div class="col-md-6 col-md-offset-3 modal-outer noPad">

                <div class="modal-body">
                    {{ Form::open(array('url'=>'/signin', 'class'=>'form-signup')) }}
                    <h3 class="form-signup-heading">Login</h3>

                    {{ Form::text('username', null, array('class'=>'form-control', 'placeholder'=>'Username')) }}
                    {{ Form::password('password', array('class'=>'form-control', 'placeholder'=>'Password')) }}


                </div>
                <div class="modal-footer">
                    {{ Form::submit('Login', array('class'=>'btn btn-large btn-success btn-block'))}}
                </div>
                {{ Form::close() }}

            </div>

        </div>

    </div>

@stop