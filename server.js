var express = require("express");
var path = require("path")
var app = express();
var bodyParser = require('body-parser')

var auth = require("./src/auth.js");

var port = 8080;

var HTML_PATH = __dirname + "/public/html/";

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
  extended: true
})); 
app.use(express.static('public'));

app.get('/', function(req, res) {
    if(auth.isAuthenticated("test")) {
        res.redirect("/home");
    } else {
        res.redirect("/login");
    }
});

app.get("/login", function(req, res) {
    if(auth.isAuthenticated("test")) {
        res.redirect("/home");
    } else {
        res.sendFile(HTML_PATH + "login.html")
    }
});

app.get("/:path", function(req, res) {
    path = req.params.path;

    if(auth.isAuthenticated("test")) {
        switch(path) {
            case "home":
                renderHome(req, res);
                break;
            default:
                res.send("404 page not found", 404);
        }
    } else {
        res.redirect("/login");
    }
});

app.post("/webservice/login", function(req, res) {
    console.log("Login Post");
    console.log("username: " + req.body.username);
    auth.authenticateCredentials(req.body.username, req.body.password, function(sessionKey) {
        if(sessionKey) {
            console.log("Authenticated");
        } else {
            console.log("Not Authenticated");
        }
    });
    res.status(200).send();
});

console.log("Server listening on port " + port);
app.listen(port);