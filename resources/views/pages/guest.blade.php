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
                   <div class="time" id="{{{$lease->id}}}">
                    {{{strtotime($lease->created_at)+$lease->expiry-time()}}}
                    </div>
                   </td>
                   </tr>
                </tbody>
            </table>
            <h4>Tip: Keep this window open to keep a track of time left while you are using this lease.</h4>
            @endif
         </div>
    </div>
@stop
@section('footer_scripts')
@parent
<script type="text/javascript">
   function component(x, v) {
        return Math.floor(x / v);
    }
    $('.time').each(function(i, obj) {   
      var div=$(this)                   
      var timestamp = div.text()
      
      if(timestamp>0)
      {
        setInterval(function() { // execute code each second

            timestamp--; // decrement timestamp with one second each second

            var hours   = component(timestamp,      60 * 60), // hours
                minutes = component(timestamp,           60) % 60, // minutes
                seconds = component(timestamp,            1) % 60; // seconds

            div.text(("0" + hours).slice(-2) + ":" + ("0" + minutes).slice(-2) + ":" + ("0" + seconds).slice(-2) ); // display
            //alert($(this).text());

        }, 1000); // interval each second = 1000 ms
      }
      else
      {
        div.text("Expired");
      }
  });

</script>
@stop