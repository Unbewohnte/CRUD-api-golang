# CRUD api
## A simple CRUD api written in Go`s standart library

API has implementation of "GET", "POST", "PUT", "DELETE" http methods, allowing to Create, Read, Update and Delete json objects in database. 

The structure of a basic json object represents Go struct "RandomData" with exported fields "title" (string), "text" (string) and unexported fields  "DateCreated" (time.Time), "LastUpdated" (time.Time) and ID (int64)

Example of a single object stored in a json database :  {
  "ID": 1618064651615612586,
  "DateCreated": "2021-04-10T14:24:11.615612068Z",
  "LastUpdated": "2021-04-10T14:24:11.615612068Z",
  "title": "Title",
  "text": "text"
 }

This project was my first take on such thing. The goal was to create a basic working example and it looks like that I`ve achieved that goal. 
