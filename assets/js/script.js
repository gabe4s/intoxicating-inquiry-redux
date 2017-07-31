function login() {
  var xhttp = new XMLHttpRequest();
  xhttp.open("POST", "/webservice/login", true);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  var username = document.getElementById("usernameBox").value;
  var password = document.getElementById("passwordBox").value;
  xhttp.send("username=" + username + "&password=" + password);
}
