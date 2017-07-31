package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
  "crypto/sha256"
  "encoding/hex"
)

func main() {
  port := "8080";

  http.HandleFunc("/img/", serveResource)
  http.HandleFunc("/css/", serveResource)
  http.HandleFunc("/js/", serveResource)
  http.HandleFunc("/", serveHtml)
  http.HandleFunc("/webservice/login", loginService)

  fmt.Println("Server listening on port " + port);
  log.Fatal(http.ListenAndServe(":" + port, nil))
}

func serveHtml(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "assets/html/" + r.URL.Path[1:] + ".html")
}

func serveResource(w http.ResponseWriter, r *http.Request) {
  path := "assets" + r.URL.Path
  http.ServeFile(w, r, path)
}

func loginService(w http.ResponseWriter, r *http.Request) {
  cookie := http.Cookie{Name: "intoxicating_inquiry_session_cookie", Value: "23fwfffw5tye", Expires: time.Now().Add(1 * 24 * time.Hour), Path: "/"}
  http.SetCookie(w, &cookie)
  r.ParseForm()
  username := r.Form.Get("username")
  password := r.Form.Get("password")
  hasher := sha256.New()
  hasher.Write([]byte(password))
  hashedPassword := hex.EncodeToString(hasher.Sum(nil))
  http.Redirect(w, r, "/home", http.StatusFound)
}
