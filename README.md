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
* **Connection pool sizes:** [1, 5, 50, 180]
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
    config.WithTxIsolation(pgtx.ReadCommitted),
    config.WithLockStrategy(seat.GetSeatWithSharedLock),
    config.WithMaxConn(5),
    config.WithMaxRetries(3),
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

## Sample Output
The test cases generate output in the below format:
```
=== RUN   TestBookSeats/IsolationLevel=READ COMMITED_LockStrategy=GetSeatWithSharedLock_PoolSize=50_Retries=3
INFO[0006] Total time taken to book the seats for trip-id: 15 is 6.723361792s 


INFO[0006] Booking details:                             
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Advait Iyer: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Aarav Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Aditya Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Arjun Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Meera Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1A for passenger Advait Iyer: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1A for passenger Aarav Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1A for passenger Arjun Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1A for passenger Aditya Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1A for passenger Meera Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1A for passenger Aarav Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1A for passenger Arjun Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1A for passenger Advait Iyer: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1A for passenger Meera Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Aryan Kaur: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Samar Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1A for passenger Atharv Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 1A is booked for passenger: Aditya Nair 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1B for passenger Aryan Kaur: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1B for passenger Samar Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1B for passenger Nirav Bajaj: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1B for passenger Omkar Sinha: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1B for passenger Aryan Kaur: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1B for passenger Samar Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1B for passenger Nirav Bajaj: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1B for passenger Omkar Sinha: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 1B is booked for passenger: Atharv Verma 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Nirav Bajaj: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Ananya Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Saanvi Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Dev Mishra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Aadhya Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Ananya Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Saanvi Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Dev Mishra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Omkar Sinha: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Aadhya Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Ananya Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Saanvi Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Dev Mishra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Aarohi Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Myra Kumar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Aadhya Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Diya Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Aarohi Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Kabir Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Ishaan Naidu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Diya Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Aarohi Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Myra Kumar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Kabir Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Ishaan Naidu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Dhruv Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Diya Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Myra Kumar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Kabir Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Ishaan Naidu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Arnav Ghosh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Viraj Rao: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Shaan Desai: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Dhruv Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Arnav Ghosh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1C for passenger Ira Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Shaan Desai: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Dhruv Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1C for passenger Arnav Ghosh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1C for passenger Viraj Rao: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 1C is booked for passenger: Ira Bhat   
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Kiara Naidu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Shaan Desai: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Inaya Iyer: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Navya Desai: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Viraj Rao: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Kiara Naidu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Inaya Iyer: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Navya Desai: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Anika Rao: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Sara Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Inaya Iyer: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Navya Desai: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Kiara Naidu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Sara Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Ahana Ghosh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Anika Rao: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Prisha Bajaj: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Sara Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Ahana Ghosh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Anika Rao: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Riya Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Tara Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Prisha Bajaj: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1D for passenger Ahana Kapoor: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Ahana Ghosh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Tara Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1D for passenger Prisha Bajaj: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Ahana Kapoor: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1D for passenger Riya Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 1D is booked for passenger: Reyansh Gupta 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Vihaan Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Ahana Kapoor: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Tara Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Riya Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1E for passenger Vihaan Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Krish Pandey: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Ayaan Kumar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Aadhira Mishra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Swara Sinha: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1E for passenger Krish Pandey: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1E for passenger Ayaan Kumar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1E for passenger Aadhira Mishra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Vihaan Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Ayaan Kumar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1E for passenger Swara Sinha: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Aadhira Mishra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Krish Pandey: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Emma Johnson: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Shanaya Kaur: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Liam Smith: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1E for passenger Rohan Chatterjee: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1E for passenger Emma Johnson: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1E for passenger Swara Sinha: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 1E is booked for passenger: Shanaya Kaur 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1F for passenger Liam Smith: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1F for passenger Rohan Chatterjee: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1F for passenger Emma Johnson: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1F for passenger Vihaan Nag: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1F for passenger Dev Bhargava: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1F for passenger Rohan Chatterjee: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1F for passenger Aarohi Bhatt: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1F for passenger Vihaan Nag: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1F for passenger Dev Bhargava: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1F for passenger Anika Nanda: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1F for passenger Aarohi Bhatt: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1F for passenger Vihaan Nag: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1F for passenger Dev Bhargava: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1F for passenger Aarohi Bhatt: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 1F for passenger Anika Nanda: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 1F for passenger Saanvi Goel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 1F for passenger Liam Smith: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 1F is booked for passenger: Riya Sodhi 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 2A for passenger Anika Nanda: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Tara Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 2A for passenger Saanvi Goel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Aadhya Khatri: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Swara Sabharwal: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Meera Gill: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 2A for passenger Tara Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 2A for passenger Saanvi Goel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 2A for passenger Swara Sabharwal: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 2A for passenger Aadhya Khatri: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 2A for passenger Tara Joshi: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 2A for passenger Meera Gill: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 2A for passenger Swara Sabharwal: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Sara Vohra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 2A for passenger Aadhya Khatri: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Tara Tandon: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Shanaya Arora: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 2A for passenger Sara Vohra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 2A for passenger Aadhira Chauhan: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 2A is booked for passenger: Meera Gill 
INFO[0006] Seat: 2B is booked for passenger: Tara Tandon 
INFO[0006] Seat: 2C is booked for passenger: Ved Agarwal 
INFO[0006] Seat: 2D is booked for passenger: Aayush Saxena 
INFO[0006] Seat: 2E is booked for passenger: Arjun Tiwari 
INFO[0006] Seat: 2F is booked for passenger: Navya Mathur 
INFO[0006] Seat: 3A is booked for passenger: Ira Grover 
INFO[0006] Seat: 3B is booked for passenger: Kiara Rao  
INFO[0006] Seat: 3C is booked for passenger: Kiara Jindal 
INFO[0006] Seat: 3D is booked for passenger: Diya Kochhar 
INFO[0006] Seat: 3E is booked for passenger: Ira Ghosh  
INFO[0006] Seat: 3F is booked for passenger: Aadhira Bajaj 
INFO[0006] Seat: 4A is booked for passenger: Shanaya Nair 
INFO[0006] Seat: 4B is booked for passenger: Om Malhotra 
INFO[0006] Seat: 4C is booked for passenger: Vivan Roy  
INFO[0006] Seat: 4D is booked for passenger: Atharva Sengupta 
INFO[0006] Seat: 4E is booked for passenger: Krishna Dubey 
INFO[0006] Seat: 4F is booked for passenger: Raghav Srivastava 
INFO[0006] Seat: 5A is booked for passenger: Aarush Kaul 
INFO[0006] Seat: 5B is booked for passenger: Neil Banerjee 
INFO[0006] Seat: 5C is booked for passenger: Reyansh Sood 
INFO[0006] Seat: 5D is booked for passenger: Dhruv Puri 
INFO[0006] Seat: 5E is booked for passenger: Kabir Kohli 
INFO[0006] Seat: 5F is booked for passenger: Ayaan Batra 
INFO[0006] Seat: 6A is booked for passenger: Vihaan Kohli 
INFO[0006] Seat: 6B is booked for passenger: Ishaan Mehra 
INFO[0006] Seat: 6C is booked for passenger: Nirav Trivedi 
INFO[0006] Seat: 6D is booked for passenger: Dev Thakur 
INFO[0006] Seat: 6E is booked for passenger: Arnav Bajpai 
INFO[0006] Seat: 6F is booked for passenger: Omkar Sood 
INFO[0006] Seat: 7A is booked for passenger: Shaan Chopra 
INFO[0006] Seat: 7B is booked for passenger: Krish Chauhan 
INFO[0006] Seat: 7C is booked for passenger: Aditya Mathur 
INFO[0006] Seat: 7D is booked for passenger: Advait Sahni 
INFO[0006] Seat: 7E is booked for passenger: Kabir Anand 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Shanaya Arora: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Sara Vohra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Aadhira Chauhan: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Viraj Kapoor: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Shanaya Arora: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Diya Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Aadhira Chauhan: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Viraj Kapoor: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Samar Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Ananya Bansal: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Krishna Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Diya Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Samar Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Viraj Kapoor: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Ananya Bansal: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Diya Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Krishna Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Myra Chawla: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Samar Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 7F for passenger Ananya Bansal: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Inaya Luthra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Myra Chawla: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Aryan Menon: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 7F for passenger Prisha Advani: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 7F for passenger Inaya Luthra: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 7F is booked for passenger: Krishna Reddy 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 8A for passenger Myra Chawla: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 8A for passenger Advait Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 8A for passenger Aryan Menon: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 8A for passenger Inaya Luthra: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 8A for passenger Meera Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 8A for passenger Advait Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 8A is booked for passenger: Prisha Advani 
INFO[0006] Seat: 8B is booked for passenger: Aryan Menon 
INFO[0006] Seat: 8C is booked for passenger: Myra Gupta 
INFO[0006] Seat: 8D is booked for passenger: Sara Khanna 
INFO[0006] Seat: 8E is booked for passenger: Olivia Brown 
INFO[0006] Seat: 8F is booked for passenger: Ethan Thompson 
INFO[0006] Seat: 9A is booked for passenger: Amelia Garcia 
INFO[0006] Seat: 9B is booked for passenger: Avi Nambiar 
INFO[0006] Seat: 9C is booked for passenger: Krishna Iyer 
INFO[0006] Seat: 9D is booked for passenger: Krishna Sharma 
INFO[0006] Seat: 9E is booked for passenger: Alexander Martinez 
INFO[0006] Seat: 9F is booked for passenger: Harper Robinson 
INFO[0006] Seat: 10A is booked for passenger: Henry Clark 
INFO[0006] Seat: 10B is booked for passenger: Evelyn Lewis 
INFO[0006] Seat: 10C is booked for passenger: Ananya Nair 
INFO[0006] Seat: 10D is booked for passenger: Michael Lee 
INFO[0006] Seat: 10E is booked for passenger: Swara Sinha 
INFO[0006] Seat: 10F is booked for passenger: Aadhya Patel 
INFO[0006] Seat: 11A is booked for passenger: Ananya Mehta 
INFO[0006] Seat: 11B is booked for passenger: Saanvi Patel 
INFO[0006] Seat: 11C is booked for passenger: Diya Singh 
INFO[0006] Seat: 11D is booked for passenger: Myra Jain 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 11E for passenger Inaya Khanna: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 11E for passenger Navya Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 11E for passenger Meera Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 11E for passenger Ira Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 11E for passenger Advait Reddy: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 11E for passenger Navya Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 11E for passenger Meera Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 11E for passenger Ira Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 11E is booked for passenger: Inaya Khanna 
INFO[0006] Seat: 11F is booked for passenger: Kiara Singh 
INFO[0006] Seat: 12A is booked for passenger: Abhinav Menon 
INFO[0006] Seat: 12B is booked for passenger: Aarush Pillai 
INFO[0006] Seat: 12C is booked for passenger: Ayaan Shetty 
INFO[0006] Seat: 12D is booked for passenger: Mia Martin 
INFO[0006] Seat: 12E is booked for passenger: Aarav Sanyal 
INFO[0006] Seat: 12F is booked for passenger: Ishan Sahu 
INFO[0006] Seat: 13A is booked for passenger: Reyansh Bhattacharya 
INFO[0006] Seat: 13B is booked for passenger: Kian Chakraborty 
INFO[0006] Seat: 13C is booked for passenger: Aarav Singh 
INFO[0006] Seat: 13D is booked for passenger: Anika Rana 
INFO[0006] Seat: 13E is booked for passenger: Aadhya Malhotra 
INFO[0006] Seat: 13F is booked for passenger: Kabir Seth 
INFO[0006] Seat: 14A is booked for passenger: Dev Kumar 
INFO[0006] Seat: 14B is booked for passenger: Ishaan Joshi 
INFO[0006] Seat: 14C is booked for passenger: Dhruv Kapoor 
INFO[0006] Seat: 14D is booked for passenger: Ahana Naidu 
INFO[0006] Seat: 14E is booked for passenger: Inaya Kumar 
INFO[0006] Seat: 14F is booked for passenger: Riya Iyer 
INFO[0006] Seat: 15A is booked for passenger: Prisha Desai 
INFO[0006] Seat: 15B is booked for passenger: Aarav Verma 
INFO[0006] Seat: 15C is booked for passenger: Meera Jain 
INFO[0006] Seat: 15D is booked for passenger: Swara Patel 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Navya Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Sara Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Aadhya Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Navya Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Ira Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Aadhya Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Samar Dutta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Navya Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Advait Bose: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Sara Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Samar Dutta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Navya Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Advait Bose: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Aadhya Bhat: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Sara Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Dhruv Mazumdar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Samar Dutta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Arnav Basu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Advait Bose: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Dhruv Mazumdar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Vihaan Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Arnav Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Arnav Basu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Dhruv Mazumdar: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Vihaan Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Arnav Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Arnav Basu: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Shaan Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Vihaan Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Krish Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Aditya Jain: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Shaan Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Arnav Patel: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Viraj Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Krish Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Shaan Mehta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Samar Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Aditya Jain: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Viraj Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Reyansh Singh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Krish Sharma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Aditya Jain: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Samar Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Viraj Gupta: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Kabir Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Reyansh Singh: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Shanaya Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Noah Wilson: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 15E for passenger Samar Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Kabir Nair: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Shanaya Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 15E for passenger Noah Wilson: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 1/3 failed: error booking seat 15E for passenger Ava Davis: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 15E is booked for passenger: Reyansh Singh 
INFO[0006] Seat: 15F is booked for passenger: Kabir Nair 
INFO[0006] Seat: 16A is booked for passenger: William Moore 
INFO[0006] Seat: 16B is booked for passenger: Sophia Taylor 
INFO[0006] Seat: 16C is booked for passenger: James Anderson 
INFO[0006] Seat: 16D is booked for passenger: Isabella Thomas 
INFO[0006] Seat: 16E is booked for passenger: Benjamin Jackson 
INFO[0006] Seat: 16F is booked for passenger: Lucas White 
INFO[0006] Seat: 17A is booked for passenger: Mason Harris 
INFO[0006] Seat: 17B is booked for passenger: Ahana Mehta 
INFO[0006] Seat: 17C is booked for passenger: Prisha Sharma 
INFO[0006] Seat: 17D is booked for passenger: Tara Singh 
INFO[0006] Seat: 17E is booked for passenger: Aarohi Sharma 
INFO[0006] Seat: 17F is booked for passenger: Riya Gupta 
INFO[0006] Seat: 18A is booked for passenger: Saanvi Reddy 
ERRO[0006] ERROR: couldn't book seat: retry 3/3 failed: error booking seat 18B for passenger Shanaya Verma: ERROR: deadlock detected (SQLSTATE 40P01) 
ERRO[0006] ERROR: couldn't book seat: retry 2/3 failed: error booking seat 18B for passenger Ava Davis: ERROR: deadlock detected (SQLSTATE 40P01) 
INFO[0006] Seat: 18B is booked for passenger: Noah Wilson 
INFO[0006] Seat: 18C is booked for passenger: Ava Davis 


INFO[0006] Final seat reservation details:              
INFO[0006] Seat: 1A, is assigned to Passenger: Aditya Nair 
INFO[0006] Seat: 1B, is assigned to Passenger: Atharv Verma 
INFO[0006] Seat: 1C, is assigned to Passenger: Ira Bhat 
INFO[0006] Seat: 1D, is assigned to Passenger: Reyansh Gupta 
INFO[0006] Seat: 1E, is assigned to Passenger: Shanaya Kaur 
INFO[0006] Seat: 1F, is assigned to Passenger: Riya Sodhi 
INFO[0006] Seat: 2A, is assigned to Passenger: Meera Gill 
INFO[0006] Seat: 2B, is assigned to Passenger: Tara Tandon 
INFO[0006] Seat: 2C, is assigned to Passenger: Ved Agarwal 
INFO[0006] Seat: 2D, is assigned to Passenger: Aayush Saxena 
INFO[0006] Seat: 2E, is assigned to Passenger: Arjun Tiwari 
INFO[0006] Seat: 2F, is assigned to Passenger: Navya Mathur 
INFO[0006] Seat: 3A, is assigned to Passenger: Ira Grover 
INFO[0006] Seat: 3B, is assigned to Passenger: Kiara Rao 
INFO[0006] Seat: 3C, is assigned to Passenger: Kiara Jindal 
INFO[0006] Seat: 3D, is assigned to Passenger: Diya Kochhar 
INFO[0006] Seat: 3E, is assigned to Passenger: Ira Ghosh 
INFO[0006] Seat: 3F, is assigned to Passenger: Aadhira Bajaj 
INFO[0006] Seat: 4A, is assigned to Passenger: Shanaya Nair 
INFO[0006] Seat: 4B, is assigned to Passenger: Om Malhotra 
INFO[0006] Seat: 4C, is assigned to Passenger: Vivan Roy 
INFO[0006] Seat: 4D, is assigned to Passenger: Atharva Sengupta 
INFO[0006] Seat: 4E, is assigned to Passenger: Krishna Dubey 
INFO[0006] Seat: 4F, is assigned to Passenger: Raghav Srivastava 
INFO[0006] Seat: 5A, is assigned to Passenger: Aarush Kaul 
INFO[0006] Seat: 5B, is assigned to Passenger: Neil Banerjee 
INFO[0006] Seat: 5C, is assigned to Passenger: Reyansh Sood 
INFO[0006] Seat: 5D, is assigned to Passenger: Dhruv Puri 
INFO[0006] Seat: 5E, is assigned to Passenger: Kabir Kohli 
INFO[0006] Seat: 5F, is assigned to Passenger: Ayaan Batra 
INFO[0006] Seat: 6A, is assigned to Passenger: Vihaan Kohli 
INFO[0006] Seat: 6B, is assigned to Passenger: Ishaan Mehra 
INFO[0006] Seat: 6C, is assigned to Passenger: Nirav Trivedi 
INFO[0006] Seat: 6D, is assigned to Passenger: Dev Thakur 
INFO[0006] Seat: 6E, is assigned to Passenger: Arnav Bajpai 
INFO[0006] Seat: 6F, is assigned to Passenger: Omkar Sood 
INFO[0006] Seat: 7A, is assigned to Passenger: Shaan Chopra 
INFO[0006] Seat: 7B, is assigned to Passenger: Krish Chauhan 
INFO[0006] Seat: 7C, is assigned to Passenger: Aditya Mathur 
INFO[0006] Seat: 7D, is assigned to Passenger: Advait Sahni 
INFO[0006] Seat: 7E, is assigned to Passenger: Kabir Anand 
INFO[0006] Seat: 7F, is assigned to Passenger: Krishna Reddy 
INFO[0006] Seat: 8A, is assigned to Passenger: Prisha Advani 
INFO[0006] Seat: 8B, is assigned to Passenger: Aryan Menon 
INFO[0006] Seat: 8C, is assigned to Passenger: Myra Gupta 
INFO[0006] Seat: 8D, is assigned to Passenger: Sara Khanna 
INFO[0006] Seat: 8E, is assigned to Passenger: Olivia Brown 
INFO[0006] Seat: 8F, is assigned to Passenger: Ethan Thompson 
INFO[0006] Seat: 9A, is assigned to Passenger: Amelia Garcia 
INFO[0006] Seat: 9B, is assigned to Passenger: Avi Nambiar 
INFO[0006] Seat: 9C, is assigned to Passenger: Krishna Iyer 
INFO[0006] Seat: 9D, is assigned to Passenger: Krishna Sharma 
INFO[0006] Seat: 9E, is assigned to Passenger: Alexander Martinez 
INFO[0006] Seat: 9F, is assigned to Passenger: Harper Robinson 
INFO[0006] Seat: 10A, is assigned to Passenger: Henry Clark 
INFO[0006] Seat: 10B, is assigned to Passenger: Evelyn Lewis 
INFO[0006] Seat: 10C, is assigned to Passenger: Ananya Nair 
INFO[0006] Seat: 10D, is assigned to Passenger: Michael Lee 
INFO[0006] Seat: 10E, is assigned to Passenger: Swara Sinha 
INFO[0006] Seat: 10F, is assigned to Passenger: Aadhya Patel 
INFO[0006] Seat: 11A, is assigned to Passenger: Ananya Mehta 
INFO[0006] Seat: 11B, is assigned to Passenger: Saanvi Patel 
INFO[0006] Seat: 11C, is assigned to Passenger: Diya Singh 
INFO[0006] Seat: 11D, is assigned to Passenger: Myra Jain 
INFO[0006] Seat: 11E, is assigned to Passenger: Inaya Khanna 
INFO[0006] Seat: 11F, is assigned to Passenger: Kiara Singh 
INFO[0006] Seat: 12A, is assigned to Passenger: Abhinav Menon 
INFO[0006] Seat: 12B, is assigned to Passenger: Aarush Pillai 
INFO[0006] Seat: 12C, is assigned to Passenger: Ayaan Shetty 
INFO[0006] Seat: 12D, is assigned to Passenger: Mia Martin 
INFO[0006] Seat: 12E, is assigned to Passenger: Aarav Sanyal 
INFO[0006] Seat: 12F, is assigned to Passenger: Ishan Sahu 
INFO[0006] Seat: 13A, is assigned to Passenger: Reyansh Bhattacharya 
INFO[0006] Seat: 13B, is assigned to Passenger: Kian Chakraborty 
INFO[0006] Seat: 13C, is assigned to Passenger: Aarav Singh 
INFO[0006] Seat: 13D, is assigned to Passenger: Anika Rana 
INFO[0006] Seat: 13E, is assigned to Passenger: Aadhya Malhotra 
INFO[0006] Seat: 13F, is assigned to Passenger: Kabir Seth 
INFO[0006] Seat: 14A, is assigned to Passenger: Dev Kumar 
INFO[0006] Seat: 14B, is assigned to Passenger: Ishaan Joshi 
INFO[0006] Seat: 14C, is assigned to Passenger: Dhruv Kapoor 
INFO[0006] Seat: 14D, is assigned to Passenger: Ahana Naidu 
INFO[0006] Seat: 14E, is assigned to Passenger: Inaya Kumar 
INFO[0006] Seat: 14F, is assigned to Passenger: Riya Iyer 
INFO[0006] Seat: 15A, is assigned to Passenger: Prisha Desai 
INFO[0006] Seat: 15B, is assigned to Passenger: Aarav Verma 
INFO[0006] Seat: 15C, is assigned to Passenger: Meera Jain 
INFO[0006] Seat: 15D, is assigned to Passenger: Swara Patel 
INFO[0006] Seat: 15E, is assigned to Passenger: Reyansh Singh 
INFO[0006] Seat: 15F, is assigned to Passenger: Kabir Nair 
INFO[0006] Seat: 16A, is assigned to Passenger: William Moore 
INFO[0006] Seat: 16B, is assigned to Passenger: Sophia Taylor 
INFO[0006] Seat: 16C, is assigned to Passenger: James Anderson 
INFO[0006] Seat: 16D, is assigned to Passenger: Isabella Thomas 
INFO[0006] Seat: 16E, is assigned to Passenger: Benjamin Jackson 
INFO[0006] Seat: 16F, is assigned to Passenger: Lucas White 
INFO[0006] Seat: 17A, is assigned to Passenger: Mason Harris 
INFO[0006] Seat: 17B, is assigned to Passenger: Ahana Mehta 
INFO[0006] Seat: 17C, is assigned to Passenger: Prisha Sharma 
INFO[0006] Seat: 17D, is assigned to Passenger: Tara Singh 
INFO[0006] Seat: 17E, is assigned to Passenger: Aarohi Sharma 
INFO[0006] Seat: 17F, is assigned to Passenger: Riya Gupta 
INFO[0006] Seat: 18A, is assigned to Passenger: Saanvi Reddy 
INFO[0006] Seat: 18B, is assigned to Passenger: Noah Wilson 
INFO[0006] Seat: 18C, is assigned to Passenger: Ava Davis 
x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  .  .  .  .  .  .  .  .  .  .  .  .
x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  .  .  .  .  .  .  .  .  .  .  .  .
x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  .  .  .  .  .  .  .  .  .  .  .  .


x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  .  .  .  .  .  .  .  .  .  .  .  .  .
x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  .  .  .  .  .  .  .  .  .  .  .  .  .
x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  x  .  .  .  .  .  .  .  .  .  .  .  .  .

```