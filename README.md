# address-api
Online address book exposed as a REST API

#### To run:
* Clone this repository to $GOPATH/src/timrxd/address-API
* cd address-api
* run `go build` and `.\address-api`.  
The api will be located at localhost:8000/users.

#### Current endpoints:

##### GET /users
* Returns list of all users

##### GET /users/{id}
* Returns a single user based on {id}

##### POST /users
* Creates a user
* Body of POST is formatted like: `{"first":"Firstname","last":"Lastname","email":"e@mail.com","phone":"xxx-xxx-xxxx"}`
* User will be assigned an ID on creation

##### PUT /users/{id}
* Updates a user
* Body of PUT is formatted like: `{"first":"Firstname","last":"Lastname","email":"e@mail.com","phone":"xxx-xxx-xxxx"}`
* All fields are optional

##### DELETE /users/{id}
* Deletes a user based on {id}

##### POST /users.csv
* Imports a set of users from a CSV file
* CSV file should be under the "csv" key in form-data body of the POST
* No spaces in between fields
* Example POST:
`POST /users/csv HTTP/1.1`  
`Host: localhost:8000`  
`Cache-Control: no-cache`    
`Content-Type: multipart/form-data;`    
`Content-Disposition: form-data; name="csv"; filename="names.csv"`    
`Content-Type: application/vnd.ms-excel`  
* File Example:  
`Roger,Federer,fed@gmail.com,908-111-2222`  
`Rafa,Nadal,rafa@gmail.com,213-456-0987`  
`Novak,Djoker,nole@gmail.com,777-777-7777`  

##### GET /users.csv
* Export all users into a CSV file format
