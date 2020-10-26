function checkStatus($led) {
    $.get('/api/check', function(jd) {
        console.log(jd)
        if (jd == "true") {
            $led.addClass('active');
            setTimeout(function() {
                checkStatus($led);
            }, 1000);
        } else if (jd == "false") {
            $led.removeClass('active');
        } else {
            alert("error: " + jd);
            $led.removeClass('active');
        }
        
    }).fail(function(r) {
        alert("unexpected error... (" + r.status + ")");
        $led.removeClass('active');
    });
}

$(function(){
    var $btn = $('.button-open');
    var $led = $('.led-circle');
    checkStatus($led);
    $btn.click(function() {
        $btn.addClass('disabled');
        $.get('/api/open', function(jd) {
            if (jd == "ok") {
                checkStatus($led);
            } else {
                alert("error: " + jd);
            }
            
        }).fail(function(r) {
            alert("unexpected error... (" + r.status + ")");
        }).always(function() {
            $btn.removeClass('disabled');   
        })
    })
})