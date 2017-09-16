function login() {

  var usernameBox = document.getElementById("usernameBox");
  var passwordBox = document.getElementById("passwordBox");

  var xhttp = new XMLHttpRequest();
  xhttp.open("POST", "/webservice/login", false);
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhttp.send("username=" + usernameBox.value + "&password=" + passwordBox.value);
  
  var statusCode = xhttp.status;
  var loginNotificationText = "";
  if(statusCode == 200) {
    location.href = "/home";
  } else if (statusCode == 401) {
    loginNotificationText = "Invalid Login Details";
    usernameBox.value = "";
    passwordBox.value = "";
  } else {
    loginNotificationText = "Unknown Error - Please Try Again Later";
  }
  document.getElementById("loginNotification").innerHTML = loginNotificationText;

}