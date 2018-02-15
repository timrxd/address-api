package main

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestAPI(t *testing.T) {

	cases := []struct {
	  method    string
    path      string
    body      User

    outCode   int
    outBody   string
	}{
    // Add first user to address book
		{"POST", "/users",User{First: "Tim", Last: "Dowd", Email: "tim@gmail.com", Phone:"123-456-7890"},
    200,`{"id":"0","first":"Tim","last":"Dowd","email":"tim@gmail.com","phone":"123-456-7890"}`},

    // Add second user to address book
    {"POST", "/users",User{First: "Jeff", Last: "Dodge", Email: "jeff@gmail.com", Phone:"111-222-3333"},
    200,`{"id":"1","first":"Jeff","last":"Dodge","email":"jeff@gmail.com","phone":"111-222-3333"}`},

    // Check if users are there
    {"GET", "/users",User{},
    200,`[{"id":"0","first":"Tim","last":"Dowd","email":"tim@gmail.com","phone":"123-456-7890"},`+
    `{"id":"1","first":"Jeff","last":"Dodge","email":"jeff@gmail.com","phone":"111-222-3333"}]`},

    // Check a specific user
    {"GET", "/users/0",User{},
    200,`{"id":"0","first":"Tim","last":"Dowd","email":"tim@gmail.com","phone":"123-456-7890"}`},

    // Check a user that doesn't exist
    {"GET", "/users/8",User{},
    404,`"User not found"`},

    // Update the name and email of a user
    {"PUT", "/users/0",User{First: "Tom", Email: "tom@gmail.com"},
    200,`{"id":"0","first":"Tom","last":"Dowd","email":"tom@gmail.com","phone":"123-456-7890"}`},

    // Delete a user
    {"DELETE", "/users/1",User{},
    200,`"User removed"`},

    // Make sure user was deleted
    {"GET", "/users",User{},
    200,`[{"id":"0","first":"Tom","last":"Dowd","email":"tom@gmail.com","phone":"123-456-7890"}]`},

	}
	for _, c := range cases {
    response := httptest.NewRecorder()
	  body,_ := json.Marshal(c.body)
    request,_ := http.NewRequest(c.method,c.path,bytes.NewBuffer(body))
    AddressRouter().ServeHTTP(response, request)

    b := response.Body.String()
    if (b[:len(b)-1] != c.outBody) {
      t.Errorf("Body Mismatch: %s %s\nResponse\t%s\nExpected\t%s",c.method,c.path,b[:len(b)-1],c.outBody)
    } else if (response.Code != c.outCode) {
      t.Errorf("Response Code Mismatch: %s %s\nResponse\t%s\nExpected\t%s",c.method,c.path,response.Code,c.outCode)
    }
	}
}
