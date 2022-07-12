'use strict';

const textContent = document.getElementById('textContent');
const getFileBtn = document.getElementById('getFileBtn');
const addPassBtn = document.getElementById('addPassBtn');
const infoField = document.getElementById('infoField');
const infoField2 = document.getElementById('infoField2');
const genPassBtn = document.getElementById('genPassBtn');
const multiAddPassBtn = document.getElementById('multiAddPassBtn');
const multiGenPassBtn = document.getElementById('multiGenPassBtn');

const swaggerBtn = document.getElementById('swaggerBtn');
const swaggerRawBtn = document.getElementById('swaggerRawBtn');

const loader = document.getElementById('loader');
const loginInputMsg = document.getElementById('loginInputMsg');
const loginInput = document.getElementById('loginInput');
const loginValidateMsg = document.getElementById('loginValidateMsg');

const passModal = document.getElementById('passModal');
const closeModal = document.getElementById('closeModal');
const pass = document.getElementById('pass');
const passes = document.getElementById('passes');
const modalMsg = document.getElementById('modalMsg');
const clipPassBtn = document.getElementById('clipPassBtn');
const clipPassesBtn = document.getElementById('clipPassesBtn');

const getFileUrlProd = 'https://green.auth.ktrn.com/api/getFile';
const addPassUrlProd = 'https://green.auth.ktrn.com/api/addPass';
const genPassUrlProd = 'https://green.auth.ktrn.com/api/generatePass';
const multiGenPassUrlProd = 'https://green.auth.ktrn.com/api/multiGeneratePass';
const swaggerUrlProd = 'https://green.auth.ktrn.com/swagger';
const swaggerRawLink = 'https://green.auth.ktrn.com/swagger/swagger.yaml';

getFileBtn.addEventListener('click', preGetFile);
addPassBtn.addEventListener('click', preAddPass);
multiAddPassBtn.addEventListener('click', preMultiAddPass);
genPassBtn.addEventListener('click', generatePass);
multiGenPassBtn.addEventListener('click', multiGeneratePass);

clipPassesBtn.addEventListener('click', clipPasses);
clipPassBtn.addEventListener('click', clipPass);
backBtn.addEventListener('click', pageBack);

swaggerBtn.addEventListener('click', swaggerUI);
swaggerRawBtn.addEventListener('click', swaggerRawUI);

window.onload = () => {
    textContent.style.display = "none";
    loginInput.style.display = "none";
    loginInputMsg.style.display = "none";
    loginValidateMsg.style.display = "none";
    genPassBtn.style.display = "none";
    multiGenPassBtn.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == passModal) {
    passModal.style.display = "none";
  }
}

closeModal.onclick = function() {
  passModal.style.display = "none";
}

loginInput.addEventListener('keyup', function (event) {
    loginValidateMsg.innerHTML = "";
});

function pageBack() {
    location.href = "index.html";
}

function swaggerUI() {
    location.href = swaggerUrlProd;
}

function swaggerRawUI() {
    location.href = swaggerRawLink;
}

function errorReloadPageMsg(id) {
    id.innerHTML = "yay, something's wrong. reload the page";
}

function errorNoAdminAccessMsg(id) {
    id.innerHTML = "no touching, that's admin stuff";
}

function preGetFile() {
    loader.innerHTML = "Loading...";
    getFile();
}

function preAddPass() {
    loader.innerHTML = "Loading...";
    addPass();
}

function preMultiAddPass() {
    loader.innerHTML = "Loading...";
    multiAddPass();
}

function clipPasses() {
    passes.select();
    document.execCommand("copy");
    alert("Data copied to your Clipboard");
}

function clipPass() {
  var TempText = document.createElement("input");
  TempText.value = pass.innerHTML;
  document.body.appendChild(TempText);
  TempText.select();
  
  document.execCommand("copy");
  document.body.removeChild(TempText);
  alert("Data copied to your Clipboard");
}

async function getFile() {
    var requestBody = {
        file : "service-map"
    }
    var jsonBody = JSON.stringify(requestBody);
  
    fetch(getFileUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '400' || response.status == '500') {
            console.log("Error while getting file");
    	    loader.innerHTML = "error";
            return Promise.reject(response);
	} else if (response.status == '200') {
    	    loader.innerHTML = "";
            return (response.json());
        } else {
    	    loader.innerHTML = "error";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        textContent.innerHTML = json.response;
        textContent.style.display = "block";
        loginInputMsg.style.display = "none";
        loginValidateMsg.style.display = "none";
        loginInput.style.display = "none";
        multiGenPassBtn.style.display = "none";
        genPassBtn.style.display = "none";
    });
}

async function addPass() {
   loader.style.display = "block";

   fetch(addPassUrlProd, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '500') {
    	    loader.innerHTML = "";
            errorReloadPageMsg(infoField);
            return Promise.reject(response);
        } else if (response.status == '403') {
    	    loader.innerHTML = "";
            errorNoAdminAccessMsg(infoField);
            return Promise.reject(response);
        } else if (response.status == '200') {
    	    loader.innerHTML = "";
            return (response.json());
        } else {
    	    loader.innerHTML = "error";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        multiGenPassBtn.style.display = "none";
        textContent.style.display = "none";
        loginInputMsg.innerHTML = "Input new login";
        loginInputMsg.style.display = "block";
        loginInput.value = json.login;
        loginInput.style.display = "block";
        genPassBtn.style.display = "block";
    });
}

async function generatePass() {
    var requestBody = {
        login : loginInput.value
    }
    var jsonBody = JSON.stringify(requestBody);

    loginValidateMsg.style.display = "block";
    loginValidateMsg.innerHTML = "Loading...";

    fetch(genPassUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '204') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "suggesting empty logins? that's a bad thing to do";
            return Promise.reject(response);
        } else if (response.status == '400') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "login must not consist spaces";
            return Promise.reject(response);
        } else if (response.status == '409') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "there's a login like that already. try another";
            return Promise.reject(response); 
        } else if (response.status == '200') {
            loginValidateMsg.innerHTML = "";
            return (response.json());
        } else {
            loginValidateMsg.style.display = "block";
	    loginValidateMsg.innerHTML = "yay, something's wrong. reload the page";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        pass.innerHTML = json.password;
        modalMsg.innerHTML = "Save that password! It will be shown just once";
	passModal.style.display = "block";
	pass.style.display = "block";
	passes.style.display = "none";
	clipPassBtn.style.display = "inline-block";
	clipPassesBtn.style.display = "none";
    });
}

async function multiGeneratePass() {
    var requestBody = {
        quantity : loginInput.value
    }
    var jsonBody = JSON.stringify(requestBody);
        
    loginValidateMsg.style.display = "block";
    loginValidateMsg.innerHTML = "Loading...";

    fetch(multiGenPassUrlProd, {
        method: 'POST',
        body: jsonBody,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '400') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "wrong quantity format";
            return Promise.reject(response);
        } else if (response.status == '400') {
            loginValidateMsg.style.display = "block";
            loginValidateMsg.innerHTML = "you can generate from 2 to 100 accounts at once";
            return Promise.reject(response); 
        } else if (response.status == '200') {
            loginValidateMsg.innerHTML = "";
            return (response.json());
        } else {
            loginValidateMsg.style.display = "block";
	    loginValidateMsg.innerHTML = "yay, something's wrong. reload the page";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        passes.innerHTML = json.response;
        modalMsg.innerHTML = "Save that list! It will shown just once";
        passModal.style.display = "block";
	passes.style.display = "block";
	pass.style.display = "none";
	clipPassBtn.style.display = "none";
	clipPassesBtn.style.display = "inline-block";
    });
}

async function multiAddPass() {
   loader.style.display = "block";

   fetch(addPassUrlProd, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: "include",
        withCredentials: true
    })
    .then(function (response) {
        if (response.status == '500') {
    	    loader.innerHTML = "";
            errorReloadPageMsg(infoField2);
            return Promise.reject(response);
        } else if (response.status == '403') {
    	    loader.innerHTML = "";
            errorNoAdminAccessMsg(infoField2);
            return Promise.reject(response);
        } else if (response.status == '200') {
    	    loader.innerHTML = "";
            return (response.json());
        } else {
    	    loader.innerHTML = "error";
	    return Promise.reject(response);
        }
    })
    .then(function (json) {
        textContent.style.display = "none";
        loginInputMsg.innerHTML = "Enter the number of new accounts";
        loginInputMsg.style.display = "block";
        loginInput.value = "2";
        loginInput.style.display = "block";
        genPassBtn.style.display = "none";
        multiGenPassBtn.style.display = "block";
    });
}
