'use strict';

const previousEmail = document.getElementById('previousEmail');
const codeField = document.getElementById('codeField');
const sendBtn = document.getElementById('sendBtn');
const infoField = document.getElementById('infoField');
const timerTextField = document.getElementById('timerTextField');
const timerField = document.getElementById('timerField');
const timerLimit = 30;
const backBtn = document.getElementById('backBtn');

const hostUrl = 'https://password.berizaryad.ru';
const sendCodeUrlProd = hostUrl + '/api/sendCode';
const validateCodeUrlProd = hostUrl +'/api/validateCode';

window.onload = () => {
    if (localStorage.getItem('email') == null) {
        location.href = "index.html";
    }

    sendCode();
}

document.addEventListener("DOMContentLoaded", startTimer);
previousEmail.innerHTML = localStorage.getItem('email');

backBtn.addEventListener('click', pageBack);
sendBtn.addEventListener('click', sendCode);

codeField.addEventListener('input', function (event) {
    infoField.innerHTML = "";
    
    if (codeField.checkValidity()) {
        validateCode();
        codeField.disabled = true;
    }
});

function errorSendCodeMsg() {
    infoField.innerHTML = "при отправке кода что-то пошло не так. давай еще";
}

function errorWrongCodeMsg() {
    infoField.innerHTML = "ой-ёй. а код не подходит";
}

function errorAttemptsLimitMsg() {
    infoField.innerHTML = "надо получить новый код";
}

function errorSendCodeLimitMsg() {
    infoField.innerHTML = "получить код больше 5 раз в 5 минут нельзя. подождем";
}

function startTimer() {
    timerTextField.innerHTML = "Повторно запросить код можно через: &nbsp";
    var seconds = countdownTimer(timerLimit);
    var interval = setInterval(() => {
        seconds = countdownTimer(seconds);
        if (seconds == 0) {
            clearInterval(interval);
            timerTextField.innerHTML = "Код просрочен";
            timerField.innerHTML = "";
            sendBtn.disabled = false;
            codeField.disabled = false;
        }
    }, 1000);
}

function countdownTimer(seconds) {
    timerField.innerHTML = seconds;
    return (seconds = seconds - 1);
}

function pageBack() {
    location.href = "index.html";
}

async function sendCode() {
    var requestBody = {
        email : localStorage.getItem('email')
    }
    var jsonBody = JSON.stringify(requestBody);

    startTimer();
    infoField.innerHTML = "";
    codeField.value = "";

    fetch(sendCodeUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "same-origin",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '400' || response.status == '500') {
            console.log("Error while sending code");
            errorSendCodeMsg();
            return;
        } else if (response.status == '429') {
            console.log("SendCode limit error");
            errorSendCodeLimitMsg();
            sendBtn.disabled = true;
            return;
        } else if (response.status == '200') {
            console.log("ok");
            sendBtn.disabled = true;
            return;
        } else {
	    return Promise.reject(response);
        }
    });
}

async function validateCode() {
    var requestBody = {
        email : localStorage.getItem('email'),
        code : codeField.value
    }
    var jsonBody = JSON.stringify(requestBody);

    var redirect = sessionStorage.getItem('redirect');
    fetch(validateCodeUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json',
            'X-Redirect-To': redirect
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        codeField.disabled = false;

        if (response.status == '500') {
            errorSendCodeMsg();
            return;
        } else if (response.status == '400') {
            errorWrongCodeMsg();
            return;
        } else if (response.status == '429') {
            errorAttemptsLimitMsg();
            return;
        } else if (response.status == '200') {
            var token = response.headers.get('X-Green-Token');
            if (token != null) {
                sessionStorage.removeItem('redirect');
	        let link = new URL(redirect);
		link.searchParams.set('greenToken', token);
		window.location.replace(link);
		return;
	    } else {
                location.href = "home.html";
                return;
	    }
        } else {
	        return Promise.reject(response);
        }
    });
}
