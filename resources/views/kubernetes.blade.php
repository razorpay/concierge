@extends('layouts.master')

@section('content')
     <div class="row">
        <div class="col-md-6 col-md-offset-3 modal-outer noPad">
            <h2>Kubernetes Services</h2>
            <table class="table table-hover table-bordered">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Namespace</th>
                    <th>Host(s)</th>
                </tr>
            </thead>
            <tbody>
            @foreach($ingresses as $ingress)
            <tr>
                {{-- <td><pre>{{var_dump($ingress)}}</pre> --}}
                <td>
                    <a href="/kubernetes/{{$ingress['metadata']['namespace']}}/{{$ingress['metadata']['name']}}" title="click here to access">
                    {{$ingress['metadata']['name']}}
                    </a>
                </td>
                <td>{{$ingress['metadata']['namespace']}}</td>
                <td>
                    @foreach($ingress['hosts'] as $host)
                    <a target=_blank refl="noopener noreferrer" href="https://{{$host}}">{{$host}}</a>
                    @endforeach
                </td>
            </tr>
            @endforeach
            </tbody>
            </table>
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
