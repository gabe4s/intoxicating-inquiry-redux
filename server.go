package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
  "crypto/sha256"
  "encoding/hex"
  "os"
  "database/sql"
  "io/ioutil"
  "encoding/json"
  "text/template"
  _ "github.com/go-sql-driver/mysql"
)

type User struct {
  user_id int
  first_name string
  last_name string
  username string
  password string
}

type TemplateVariables struct {
  Title string
  Name string
  Content string
}

type Config struct {
  Db_Username string  `json:"db_username"`
  Db_Password string  `json:"db_password"`
  Db_Schema string  `json:"db_schema"`
}

var db *sql.DB
var config Config

func main() {
  port := "8080"

  if len(os.Args) > 1 {
    port = os.Args[1]
  }

  http.HandleFunc("/img/", serveResource)
  http.HandleFunc("/css/", serveResource)
  http.HandleFunc("/js/", serveResource)
  http.HandleFunc("/", handlePath)
  http.HandleFunc("/webservice/login", loginService)

  fmt.Println("Intoxicating Inquiry server listening on port " + port)
  log.Fatal(http.ListenAndServe(":" + port, nil))
}

func getConfig() Config {
  if(config == Config{}) {
    raw, err := ioutil.ReadFile("./config/config.json")
    if(err != nil) {
      log.Fatal("Could not load config, exiting...")
    }
    json.Unmarshal(raw, &config)
  }

  return config
}

func getDatabaseConnection() *sql.DB {
  if(db == nil) {
    config := getConfig()

    var err error
    db, err = sql.Open("mysql", config.Db_Username + ":" + config.Db_Password + "@/" + config.Db_Schema + "?charset=utf8")
    if(err != nil) {
      log.Fatal("Could not connect to database, exiting...")
    }
  }

  return db
}

func getUserFromDbByCredentials(username string, passwordHash string) User {
  var user User
  err := getDatabaseConnection().QueryRow("SELECT * FROM user WHERE username=? AND password=? LIMIT 1", username, passwordHash).Scan(&user.user_id, &user.first_name, &user.last_name, &user.username, &user.password)

  if(err != nil) {
    fmt.Println(err)
    return User{}
  } else {
    return user
  }
}

func getUserFromDbBySessionKey(sessionKey string) User {
  var user User
  err := getDatabaseConnection().QueryRow("SELECT u.* FROM user u JOIN session s ON u.user_id=s.user_id WHERE s.session_key=?", sessionKey).Scan(&user.user_id, &user.first_name, &user.last_name, &user.username, &user.password)

  if(err != nil) {
    fmt.Println(err)
    return User{}
  } else {
    return user
  }
}

func handlePath(res http.ResponseWriter, req *http.Request) {
  requestedUrlPath := req.URL.Path[1:]
  fmt.Println(requestedUrlPath)

  user := User{}

  cookie, err := req.Cookie("intoxicating_inquiry_session_cookie")
  if(err != nil) {
    logError(err)
  }

  if(cookie != nil) {
    sessionKey := cookie.Value
    user = getUserFromDbBySessionKey(sessionKey)
  }

  isValidSession := false
  if(user != (User{})) {
    isValidSession = true
  }

  if(requestedUrlPath == "login") {
    if(isValidSession) {
      http.Redirect(res, req, "/home", http.StatusFound)
    } else {
      serveHtml(res, req, user, "login")
    }
  } else {
    if(isValidSession) {
      serveHtml(res, req, user, requestedUrlPath)
    } else {
      http.Redirect(res, req, "/login", http.StatusFound)
    }
  }

}

func serveHtml(res http.ResponseWriter, req *http.Request, user User, path string) {
  if(path == "login") {
    http.ServeFile(res, req, "assets/html/login.html")
  } else {
    templateFileName := "assets/html/template.html"
    template := template.Must(template.ParseFiles(templateFileName))
    if (path == "home") {
      templateVariables := TemplateVariables{"Home", user.first_name + " " + user.last_name, ""}
      template.Execute(res, templateVariables)
    }
  }
}

func serveResource(res http.ResponseWriter, req *http.Request) {
  path := "assets" + req.URL.Path
  http.ServeFile(res, req, path)
}

func loginService(res http.ResponseWriter, req *http.Request) {
  req.ParseForm()
  username := req.Form.Get("username")
  password := req.Form.Get("password")
  hasher := sha256.New()
  hasher.Write([]byte(password))
  passwordHash := hex.EncodeToString(hasher.Sum(nil))

  t := time.Now()

  sessionKeyString := username + ":" + passwordHash + ":" + t.String()
  fmt.Println("SessIDString: " + sessionKeyString)

  hasher = sha256.New()
  hasher.Write([]byte(sessionKeyString))
  sessionKeyHash := hex.EncodeToString(hasher.Sum(nil))
  fmt.Println("SessIdHash: " + sessionKeyHash)

  responseCode := http.StatusUnauthorized

  user := getUserFromDbByCredentials(username, passwordHash)
  if(user != (User{})) {
    stmt, err := db.Prepare("INSERT INTO session (`session_key`, `user_id`) VALUES (?, ?)")
    fmt.Println("Inserting sessionKey: " + sessionKeyHash)
    if(err != nil) {
      logError(err)
    } else {
      _, err = stmt.Exec(sessionKeyHash, user.user_id)
      if(err != nil) {
        logError(err)
      } else {
        fmt.Println("User Logged In: ", user.username)
        cookie := http.Cookie{Name: "intoxicating_inquiry_session_cookie", Value: sessionKeyHash, Expires: time.Now().Add(1 * 24 * time.Hour), Path: "/"}
        http.SetCookie(res, &cookie)
        responseCode = http.StatusOK
      }
    }
  }

  res.WriteHeader(responseCode)

}

func logError(err error) {
  log.Println("Error: ", err)
}
