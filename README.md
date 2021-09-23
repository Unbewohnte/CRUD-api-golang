# CRUD api
## A simple CRUD api written in Go

API has implementation of "GET", "POST", "PATCH", "DELETE" http methods, allowing to Read, Create, Update and Delete objects in sqlite3 database via json input.

---

## Status

Implemented:
- **GET**
- **POST**
- **PATCH**
- **DELETE**

## Examples
- `curl localhost:8000/randomdata` - to get EVERYTHING (obviously a bad idea if you have lots of data)
- `curl localhost:8000/randomdata -H "content-type:application/json" -d '{"title":"This is a title","text":"This is a text"}' -X POST` - to create a new RandomData (IDs are created automatically from 1-∞)
- `curl localhost:8000/randomdata/1` - to get the first RandomData you`ve created
- `curl localhost:8000/randomdata/1 -H "content-type:application/json" -d '{"title":"This is an updated title","text":"This is an updated text"}' -X PATCH` - to update the first RandomData
- `curl localhost:8000/randomdata/1  -X DELETE` - to delete the first RandomData

---

It`s not a recommended or even the correct way of doing a CRUD api of such sort, I'm just practicing  

