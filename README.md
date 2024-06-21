# Airline Reservation System

## Project Overview
This project aims to analyze the behavior of a Relational DB in general and Postgres in particular under high resource contention scenarios, with a focus where the following configurations are varied

* Isolation levels
* Locking mechanisms
* Connection pool sizes

The system simulates an airline reservation process where multiple passengers book seats on flights in parallel.

## Prerequisites
* Docker 
* Docker Compose 
* Golang 1.22 or later

## Infrastructure
* Machine: MacBook M2 Pro
* Docker Image: The images are leanest and the most secure images available for Postgres. The Postgres image is based on the official Postgres image occupies just 95 MB in size. Please check [chaingaurd-images](https://www.chainguard.dev/chainguard-images) for more details. They are awesome!

## Setup/TearDown
* Clone the repository
* Repository contains a `Makefile` with the following targets of interest:
* ```make setup``` - Sets up the environment:
  * Brings up a Postgres docker instance with pre-populated airline data.
* ```make teardown``` - Tears down the environment.

## Postgres Instance Coordinates
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: airline_reservation_db

## DB Schema and Prepopulated Test Data
* The initial schema and data are populated through the `create_and_populate_tables.sh` script.
* The schema consists of the following tables:
  * `airline` - Contains the airline details.
  * `passengers` - Contains the passenger details.
  * `trip` - Contains the trip details.
  * `reservations` - Contains the reservation details.
* The tables are pre-populated with data for testing purposes.

## Testing
The project includes table-driven tests that cover various combinations of connection pool sizes, locking strategies, and isolation levels. 
</br>The tests are designed to observe and validate the behavior of the system under these configurations.


### Test Suite
**make test** - Runs the test cases with the below combinations of test variables.
* **Connection pool sizes:** [1, 5, 10, 20]
* **Isolation levels:** 
    * Read Committed
    * Repeatable Read
    * Serializable.
  </br>**Note:** `Read Uncommitted` is not tested as the behaviour is similar to Read Committed in Postgres.
* **Locking mechanisms:** 
    * None
    * Shared
    * Exclusive

**Note:** One can customize the combination in `booking_test.go`

**Example Test Case**: Book seats with a connection pool of size `50` and `exclusive-lock` strategy under `READ COMMITTED` isolation level.

### Testing custom test scenario:
You can run a custom test scenario by setting the options in `main.go`:
```go
cfg := config.NewConfig(
    config.WithMaxConn(50),
    config.WithLockStrategy(seat.GetSeatWithExclusiveLock),
    config.WithTxIsolation(pgtx.READ_COMMITTED),
)
```

## Setting Resource Limits and Configurations
You can also try to test the system and DB behavior by varying the CPU and memory limits for the DB and the application. This can be done by setting the CPU and memory limits of the PostgreSQL container in `docker-compose.yml` 
</br>You can adjust the configuration parameters in `postgresql.conf` to test Postgres behaviour for different set of configurations.


## Getting Started
It is very easy to setup the project and get with the project.</br>
Clone the repository and running the tests using the below commands should get you started:
```bash
git clone git@github.com:soumya-codes/airline-reservation-poc.git
make test
 ```
**Note**: The project uses `go version 1.22`, you may have to tidy the Go modules before running the tests.