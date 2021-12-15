# Golang playground
Simple CRUD app for <b> job offers</b> written in Go to explore the language and popular libraries. 

<h3>Libraries </h3>

- ORM https://gorm.io/ 
- Logger https://github.com/uber-go/zap
- Web https://github.com/gin-gonic/gin
- CLI https://github.com/urfave/cli/v2

Storing data in PostgreSQL v11.

Project layout https://github.com/golang-standards/project-layout

<h3>First Run</h3>
Make sure your PostgreSQL instance is up and running. 
To create table in your preferred db  just copy content of https://github.com/vladimirconev/go-playground/blob/main/setup_table.sql and execute. <br/>

<h3> Starting server </h3>

Start server on default port `3456` by invoking `go run cmd/main.go server` .

To override any of already pre-set variables:
- server port 8888 `go run cmd/main.go server --server-port "8888"`
- your PostgreSQL password `go run cmd/main.go server --postgres-password "your_pass"`
- and you get the idea already ... of course you can chain them and override multiple variables 

Under `/docs` folder there is a Postman collection ready to be imported and start playing around.

Running tests `go test ./... -short`.

Happy Coding!!!
