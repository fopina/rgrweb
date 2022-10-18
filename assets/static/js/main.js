function checkStatus($led) {
    $.get('api/check', function(jd) {
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
        alert("Error " + r.status + ": " + r.responseText);
        $led.removeClass('active');
    });
}

$(function(){
    var $btn = $('.button-open');
    var $led = $('.led-circle');

    var token = window.location.hash.substring(1)
    if (token != "") {
        $.ajaxSetup({
            headers:{
               'X-Token': token
            }
        });
    }

    checkStatus($led);
    $btn.click(function(event) {
        event.preventDefault();
        $btn.addClass('disabled');
        $.get('api/open', function(jd) {
            if (jd == "ok") {
                checkStatus($led);
            } else {
                alert("error: " + jd);
            }
            
        }).fail(function(r) {
            alert("Error " + r.status + ": " + r.responseText);
        }).always(function() {
            $btn.removeClass('disabled');   
        })
    })
})
