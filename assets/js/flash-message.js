var data = $(".custom-flash").text();
if(data!='')
{
    $(".custom-flash").show();
    $(".custom-flash").animate({
        top: "80px"},
        '2000',
        'linear'
    );
    $(".custom-flash").delay(4000).animate({
        top : "-50px"},
        '2000',
        'linear',
        function() {
            $(".custom-flash").hide();
        }
    );
}
