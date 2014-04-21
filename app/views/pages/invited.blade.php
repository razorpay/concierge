@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
            <h2>Guest Access Lease Created</h2>
            <table class="table table-hover table-bordered">
                <thead>
                    <tr>
                    <th>Creator</th>
                    <th>Security Group</th>
                    <th>Protocol</th>
                    <th>Port(s)</th>
                    <th>Expiry</th>
                    </tr>
                </thead>
                <tbody>
                   <tr>
                   <td>{{{User::find($invite['user_id'])->username}}}</td>
                   <td>{{{$invite['group_id']}}}</td>
                   <td>{{{$invite['protocol']}}}</td>
                   <td>{{{$invite['port_from']}}}-{{{$invite['port_to']}}}</td>
                   <td>
                   <?php
                    //Calculating time to expiry in hours & minutes
                    $hours=intval(floor($invite['expiry']/3600)); 
                    $minutes=intval(floor(($invite['expiry']-$hours*3600)/60));
                    echo "$hours hours $minutes minutes";
                   ?>
                    </td>
                   </tr>
                </tbody>
            </table>
            <h4>Share the following link to invite someone for this lease</h4>
            <input type="text" onclick="this.focus();this.select()" class="form-control" readonly="readonly" value="{{URL::to('/invite/')}}/{{$invite['token']}}">
            
           <h4></h4>
            <p>Information: The link will only work once.</p>
         </div>
    </div>
@stop