
document.getElementById('login-button').addEventListener("click", event => {
    event.preventDefault();
    const server = document.getElementById("server").value,
          login = document.getElementById("login").value,
          password = document.getElementById("password").value;
          
    var xhr = new XMLHttpRequest();

    var body = 'server=' + encodeURIComponent(server) +
    '&login=' + encodeURIComponent(login) +
    '&password=' + encodeURIComponent(password);
    
    xhr.open("POST", '/auth', true);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.send(body);    
});
