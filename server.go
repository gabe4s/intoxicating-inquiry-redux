package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
)

type user struct {
  name string
}

func main() {
  port := "8080";

  http.HandleFunc("/api/test", Test)
  http.HandleFunc("/login", Login)
  http.HandleFunc("/img/", serveResource)
  http.HandleFunc("/css/", serveResource)

  fmt.Println("Server listening on port " + port);
  log.Fatal(http.ListenAndServe(":" + port, nil))
}

func Test(w http.ResponseWriter, r *http.Request) {
  usera := user{name: "testname"}
  w.Header().Set("Content-Type", "application/json")
  j, _ := json.Marshal(usera)
  w.Write(j)
}

func Login(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "assets/html/login.html")
}

func serveResource(w http.ResponseWriter, req *http.Request) {
  path := "assets" + req.URL.Path
  http.ServeFile(w, req, path)
}
