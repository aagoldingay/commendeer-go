[![Build Status](https://travis-ci.org/aagoldingay/commendeer-go.svg?branch=master)](https://travis-ci.org/aagoldingay/commendeer-go)

# commendeer-go

"Commend here" - A Final Year Project prototype for a feedback system 

Developed using Go 1.11.x

Follow [this guide](https://grpc.io/docs/quickstart/go.html) to install protobufs and gRPC

## Guide to run

You must first have the Postgres database set up. Insert password when prompted, if necessary

```
psql -U [username] -f [path-to]\dbconfig\db_create.sql
psql -U [username] -f [path-to]\dbconfig\db_populate.sql
```

Import packages:

```
go get -u github.com/lib/pq
```
