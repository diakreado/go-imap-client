
function removeErrors() {
    document.getElementById("login").classList.remove('error')
    document.getElementById("label-login").classList.remove('error')
    document.getElementById("password").classList.remove('error')
    document.getElementById("label-password").classList.remove('error')
}

document.getElementById('login-block').addEventListener("click", event => {
    removeErrors();
});

if (document.getElementById('logout-button')) {
    document.getElementById('logout-button').addEventListener("click", event => {
        event.preventDefault();
        window.location.replace("/auth");
    });
}

if (document.getElementById('login-button')) {
    document.getElementById('login-button').addEventListener("click", event => {
        event.preventDefault();
        const server = document.getElementById("server").value,
              login = document.getElementById("login").value,
              password = document.getElementById("password").value;
        
        if (!server || !login || !password) {
            if (!login) {
                document.getElementById("login").classList.add('error')
                document.getElementById("label-login").classList.add('error')
            }
            if (!password) {
                document.getElementById("password").classList.add('error')
                document.getElementById("label-password").classList.add('error')
            }
            return;
        }
            
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function() {
            if (xhr.readyState == XMLHttpRequest.DONE) {
                window.location.reload();
            }
        }
    
        const body = 'server=' + encodeURIComponent(server) +
        '&login=' + encodeURIComponent(login) +
        '&password=' + encodeURIComponent(password);
    
        xhr.open("POST", '/auth', true);
        xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
        xhr.send(body);
    });
} 
