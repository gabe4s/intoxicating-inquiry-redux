var exports = module.exports = {};
var mysql = require("mysql");

var connection = mysql.createConnection({
    host: "",
    user: "",
    password: "",
    database: ""
});

connection.connect(function(err) {
    if(err) {
        console.log(err);
    } else {
        console.log("Connected!");
    }
});

exports.isAuthenticated = function(sessionKey) {
    return false;
}

exports.authenticateCredentials = function(username, password, cb) {
    var sql = "SELECT * FROM user WHERE username=? LIMIT 1";
    connection.query(sql, [username], function(err, results) {
        if(err) throw err;
        if(results.length == 1) {
            console.log(results);
        } else {

        }
    });
}