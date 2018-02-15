package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "io"
  "log"
  "net/http"
  "strconv"
  "strings"
)

var db []User
var nextID int

type User struct {
  ID    string `json:"id,omitempty"`
  First string `json:"first,omitempty"`
  Last  string `json:"last,omitempty"`
  Email string `json:"email,omitempty"`
  Phone string `json:"phone,omitempty"`
}

func main() {
  router := AddressRouter()
  log.Fatal(http.ListenAndServe(":8000", router))
}

func AddressRouter() *mux.Router{
  router := mux.NewRouter()
  router.HandleFunc("/users", GetUsers).Methods("GET")
  router.HandleFunc("/users/{id}", GetUser).Methods("GET")
  router.HandleFunc("/users", CreateUser).Methods("POST")
  router.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
  router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
  router.HandleFunc("/users.csv", ExportUsers).Methods("GET")
  router.HandleFunc("/users.csv", ImportUsers).Methods("POST")
  return router
}

// Return all users in address book
func GetUsers(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(db)
}

// Return specific user
func GetUser(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  for _, i := range db {
    if i.ID == params["id"] {
      json.NewEncoder(w).Encode(i)
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")
}

// Create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
  var newUser User
  _ = json.NewDecoder(r.Body).Decode(&newUser)
  newUser.ID = strconv.Itoa(nextID)
  nextID++

  db = append(db, newUser)

  w.Header().Set("Location",fmt.Sprintf("%s/users/%s", r.Host, newUser.ID))
  json.NewEncoder(w).Encode(newUser)
}

// Update a particular user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  for n := range db {
    if db[n].ID == params["id"] {
      var newData User
      _ = json.NewDecoder(r.Body).Decode(&newData)

      if (newData.First != "") {
        db[n].First = newData.First
      }
      if (newData.Last != "") {
        db[n].Last = newData.Last
      }
      if (newData.Email != "") {
        db[n].Email = newData.Email
      }
      if (newData.Phone != "") {
        db[n].Phone = newData.Phone
      }

      json.NewEncoder(w).Encode(db[n])
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")
}

// Delete a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  for n := range db {
    if db[n].ID == params["id"] {
      db[n] = db[len(db)-1]
      db = db[:len(db)-1]
      json.NewEncoder(w).Encode("User removed")
      return
    }
  }
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode("User not found")
}

func ImportUsers(w http.ResponseWriter, r *http.Request) {
  var Buf bytes.Buffer
  file, _, err := r.FormFile("csv")
  if err != nil {
      json.NewEncoder(w).Encode("Invalid file")
      return
  }
  defer file.Close()

  io.Copy(&Buf, file)
  contents := Buf.String()
  for _,line := range strings.Split(contents,"\n") {
    if (strings.Count(line, ",") == 3) {
      newData := strings.Split(strings.TrimSpace(line),",")
      var newUser User
      newUser.ID = strconv.Itoa(nextID)
      newUser.First = newData[0]
      newUser.Last = newData[1]
      newUser.Email = newData[2]
      newUser.Phone = newData[3]
      nextID++

      db = append(db,newUser)
    }
  }
  Buf.Reset()
  return
}

func ExportUsers(w http.ResponseWriter, r *http.Request) {
  output := ""
  for _,user := range db {
    output = output + fmt.Sprintf("%s,%s,%s,%s,%s\n",
                      user.ID,user.First,user.Last,user.Email,user.Phone)
  }
  readOut := strings.NewReader(output)
  io.Copy(w, readOut)
  return
}
