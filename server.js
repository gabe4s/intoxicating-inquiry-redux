var express = require("express");
var path = require("path")
var app = express();

var port = 8080;

app.use(express.static('public'));

app.get('/', function (req, res) {
    res.sendFile(__dirname + "/public/html/login.html")
});

app.get('/login', function (req, res) {
    res.sendFile(__dirname + "/public/html/login.html")
});

app.listen(port);
console.log("Server listening on port " + port);