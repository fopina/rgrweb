var token = window.location.hash.substring(1)
var fetchOptions = token != "" ? {
    headers: {
        'X-Token': token
    },
} : {}

function handleErrors(error) {
    if (error instanceof Array) {
        if (error[0] == 1) {
            // bad response status
            if (error[1] == 403) {
                alert('Error: invalid token');
            } else {
                alert('Error: unexpected (' + error[1] + ')');
            }
        } else if (error[0] == 2) {
            // error in response text
            alert('Error: ' + error[1]);
        }
    } else {
        alert(error);
    }
}

function checkStatus($led) {
    fetch('api/check', fetchOptions)
    .then((response) => {
        if (response.status !== 200) {
            throw [1, response.status];
        }
        return response.text()
    })
    .then((jd) => {
        console.log("server reply:", jd);
        if (jd === "true") {
            $led.addClass('active');
            setTimeout(function() {
                checkStatus($led);
            }, 1000);
        } else if (jd === "false") {
            $led.removeClass('active');
        } else {
            throw [2, jd];
        };
    })
    .catch((error) => {
        handleErrors(error);
        $led.removeClass('active');
    })
}

var $btn = u('.button-open');
var $led = u('.led-circle');

checkStatus($led);

$btn.on('click', function(event) {
    event.preventDefault();
    $btn.addClass('disabled');
    fetch('api/open', fetchOptions)
    .then((response) => {
        if (response.status !== 200) {
            throw [1, response.status];
        }
        return response.text()
    })
    .then((jd) => {
        if (jd == "ok") {
            checkStatus($led);
        } else {
            throw [2, jd];
        }
    })
    .catch((error) => {handleErrors(error)})
    .finally(() => {$btn.removeClass('disabled')});
})
