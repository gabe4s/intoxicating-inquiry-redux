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
  _ "github.com/go-sql-driver/mysql"
)

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

func getUserIdFromDb(username string, passwordHash string) int {
  user_id := -1
  getDatabaseConnection().QueryRow("SELECT user_id FROM user WHERE username=? AND password=? LIMIT 1", username, passwordHash).Scan(&user_id)

  return user_id
}

func isValidSession(cookie *http.Cookie) bool {
  returnValue := false

  if(cookie != nil) {
    sessionId := cookie.Value
    fmt.Println("SessionId: " + sessionId)
    user_id := -1
    getDatabaseConnection().QueryRow("SELECT user_id FROM session WHERE session_id=? LIMIT 1", sessionId).Scan(&user_id)
    fmt.Println("UserId: ", user_id)
    if(user_id != -1) {
      returnValue = true
    }
  }

  fmt.Println("ValidSession: ", returnValue)

  return returnValue
}

func handlePath(w http.ResponseWriter, r *http.Request) {
  requestedUrlPath := r.URL.Path[1:]
  fmt.Println(requestedUrlPath)

  cookie, err := r.Cookie("intoxicating_inquiry_session_cookie")
  if(err != nil) {
    logError(err)
  }

  if(requestedUrlPath == "login") {
    if(isValidSession(cookie)) {
      http.Redirect(w, r, "/home", http.StatusFound)
    } else {
      serveHtml(w, r, "login")
    }
  } else {
    if(isValidSession(cookie)) {
      serveHtml(w, r, requestedUrlPath)
    } else {
      http.Redirect(w, r, "/login", http.StatusFound)
    }
  }

}

func serveHtml(w http.ResponseWriter, r *http.Request, path string) {
  if(path == "login") {
    http.ServeFile(w, r, "assets/html/login.html")
  } else if (path == "home"){
    http.ServeFile(w, r, "assets/html/template.html")
  }
}

func serveResource(w http.ResponseWriter, r *http.Request) {
  path := "assets" + r.URL.Path
  http.ServeFile(w, r, path)
}

func loginService(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  username := r.Form.Get("username")
  password := r.Form.Get("password")
  hasher := sha256.New()
  hasher.Write([]byte(password))
  passwordHash := hex.EncodeToString(hasher.Sum(nil))

  t := time.Now()

  sessionIdString := username + ":" + passwordHash + ":" + t.String()
  fmt.Println("SessIDString: " + sessionIdString)

  hasher = sha256.New()
  hasher.Write([]byte(sessionIdString))
  sessionIdHash := hex.EncodeToString(hasher.Sum(nil))
  fmt.Println("SessIdHash: " + sessionIdHash)

  responseCode := http.StatusUnauthorized

  userId := getUserIdFromDb(username, passwordHash)
  if(userId != -1) {
    stmt, err := db.Prepare("INSERT INTO session (`session_id`, `user_id`) VALUES (?, ?)")
    fmt.Println("Inserting SessionId: " + sessionIdHash)
    if(err != nil) {
      logError(err)
    } else {
      _, err = stmt.Exec(sessionIdHash, userId)
      if(err != nil) {
        logError(err)
      } else {
        fmt.Println("Logged in user: ", userId)
        cookie := http.Cookie{Name: "intoxicating_inquiry_session_cookie", Value: sessionIdHash, Expires: time.Now().Add(1 * 24 * time.Hour), Path: "/"}
        http.SetCookie(w, &cookie)
        fmt.Println("Sending redirect")
        responseCode = http.StatusOK
      }
    }
  }

  w.WriteHeader(responseCode)

}

func logError(err error) {
  log.Println("Error: ", err)
}
