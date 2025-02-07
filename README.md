### How to use

Environment variables necessary to run this project are stored in `.env`, according to the template `.env.dist`.

After setting these environment variables, the project can be run using `make`.

`make run` starts the web server, after which the server listens for any HTTP requests.

I used [Postman](https://www.postman.com/downloads/) for my HTTP requests. Below you can find an example request.

```
POST /payment-request
Host: localhost:8080 
Authorization: Basic
    - username: redacted (according to .env)
    - password: redacted (according to .env)
Content-Type: application/json

{ "debtor_iban": "FR1112739000504482744411A64", "debtor_name": "company1", "creditor_iban": "DE65500105179799248552", "creditor_name": "beneficiary", "ammount": 42.99, "idempotency_unique_key": "JXJ984XXXZ" }
```

You can also find a Python file to do some bulk requests. I used this to test the async capabilities of the application.

### Project structure

numeral  
Root folder. Contains `.env` for environment variables, `Makefile` for easy launching of the project and the README.

./data  
Sample data

./cmd/server  
Home of the main package and the web server launch function

./internal/payment  
Home of the payment package containing the services. Mainly the HTTP handler as well as the bank connector. 

### Customer journery

- (Re-)initialize DB
- REST API server exposing `/payment-request` on `http://localhost:8080`
- When request received, do the following:
    - Check authorization with user details stored in env vars
    - Check whether HTTP request is POST method
    - Decode request into payment struct
    - Validate request according to given JSONSchema
    - Store payment in DB using goroutine
    - Try to store payment as XML according to given XML in location of BANK_FOLDER env var using goroutine
    - Start goroutine that waits for 'bank' to process the payment and updates the DB entry of the idempotency ID
    - Wait for payment to be stored in DB, then send back 200 HTTP status with payment encoded

### Design

Go in production seems to be a hotly debated topic. Found many differing ideas about project organization. Will link these two that I found most promising. I ended up going with a server project according to the official Go guide.

- [Official Go project organization](https://go.dev/doc/modules/layout)
- [Project organization in Go at SoundCloud](http://peter.bourgon.org/go-in-production/)

Personally also found it quite confusing how exactly to handle environment variables, (relative) paths (different in tests than at runtime?!) and error/panic returns. Lots to learn still! :)

Decided to use `handler.go` as my main orchestrator throughout the customer journey, with most if not all functionality in the `internal` folder being called from the main HTTP handler `HandleCreatePayment`.

Found out pretty late in the three hours about the `idempotency_unique_key` and its meaning as to handling multiple payment requests with a shared idempotency ID and that all tying into the DB design. Definitely a bit much for the three hours I had.

Introduced asynchronicity in this later stage, because DB handling and waiting for the 'bank' response should definitely be done asynchronously.

### Changes from 3 hour version

- SQLite DB support
- Listen for bank response in a goroutine, update DB
- Working tests for HTTP validation (not complete test coverage)
- Use panic instead of log.Fatal to be more Go idiomatic
- Make everything run asynchronously

### Limitations

- The project is missing the following features and/or functionality:
    - XML validation
- Besides that, limitations from my perspective:
    - No complete test coverage, just showing off some testing
    - Can be expanded on in all kinds of directions (e.g. auth with middleware, having a separate payments service, DB connection pool, better error handling, etc.)
