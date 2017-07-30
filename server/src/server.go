package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "encoding/json"
)

func main() {
  port := "8080";

  fmt.Println("Server listening on port " + port);
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/api/test", Test).Methods("GET")

  log.Fatal(http.ListenAndServe(":" + port, router))
}

func Test(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  j, _ := json.Marshal("TEST TEXT")
  w.Write(j)
}
