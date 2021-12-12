# go-playground
simple CRUD app written in Go to explore and learn 

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

Start server on default port 3456 by invoking:
`go run cmd/main.go server`.

<h3> Sample payloads </h3> 

- POST localhost:3456/offers <br/>
```json
{
    "company": "test_company",
    "email": "test.e@e-on.com",
    "expiration_date": "2021-05-25 17:27:43.48878",
    "link": "http://test.com",
    "details": "we are looking for Jr Python developer...",
    "salary": 4500.00,
    "phone": "+38978323177"
}
```
- PUT localhost:3456/offers/:offerID <br/>
```json
{
    "email": "test2.e@e-on.com",
    "link": "http://test2.com",
    "salary": 6500.00,
    "phone": "+38978323177"
}
```


Happy Coding!!!
