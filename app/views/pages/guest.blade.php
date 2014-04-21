@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
            <h2>Guest Access Lease</h2>
            @if(isset($failure))
            <h3>Lease Creation Failed</h3>
            <p>{{{$failure}}}</p>
            @else
            <h3>Lease acquired Successfully</h3>
            <table class="table table-hover table-bordered">
                <thead>
                    <tr>
                    <th>Creator</th>
                    <th>Leased IP</th>
                    <th>Security Group</th>
                    <th>Protocol</th>
                    <th>Port(s)</th>
                    <th>Time Left</th>
                    </tr>
                </thead>
                <tbody>
                   <tr>
                   <td>{{{$lease->user->username}}}</td>
                   <td>{{{$lease->lease_ip}}}</td>
                   <td>{{{$lease->group_id}}}</td>
                   <td>{{{$lease->protocol}}}</td>
                   <td>{{{$lease->port_from}}}-{{{$lease->port_to}}}</td>
                   <td>
                   <?php
                    //Calculating time to expiry in hours & minutes
                     $time_left=strtotime($lease->created_at)+$lease->expiry-time(); 
                     $hours=intval(floor($time_left/3600)); 
                     $minutes=intval(floor(($time_left-$hours*3600)/60));
                     echo "$hours hours $minutes minutes";
                   ?>
                    </td>
                   </tr>
                </tbody>
            </table>
            @endif
         </div>
    </div>
@stop