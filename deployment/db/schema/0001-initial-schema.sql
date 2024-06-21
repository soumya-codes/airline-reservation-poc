-- Create tables
CREATE TABLE passenger
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE airline
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE trip
(
    id         SERIAL PRIMARY KEY,
    airline_id INT       NOT NULL,
    schedule   TIMESTAMP NOT NULL,
    completed  BOOLEAN   NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_airline FOREIGN KEY (airline_id) REFERENCES airline (id) ON DELETE CASCADE
);

CREATE TABLE reservation
(
    id           SERIAL PRIMARY KEY,
    seat_id      VARCHAR(10) NOT NULL,
    passenger_id INT,
    trip_id      INT         NOT NULL,
    CONSTRAINT fk_passenger FOREIGN KEY (passenger_id) REFERENCES passenger (id) ON DELETE SET NULL,
    CONSTRAINT fk_trip FOREIGN KEY (trip_id) REFERENCES trip (id) ON DELETE CASCADE,
    CONSTRAINT unique_trip_seat UNIQUE (trip_id, seat_id),
    CONSTRAINT unique_passenger_trip UNIQUE (passenger_id, trip_id)
);

-- Create index
CREATE INDEX idx_trip_schedule_airline_id ON trip (schedule, airline_id);
CREATE INDEX idx_reservation_trip_id_passenger_id_seat_id ON reservation (trip_id, passenger_id, seat_id) WHERE passenger_id IS NULL;

-- Trigger function definition to add seats
CREATE OR REPLACE FUNCTION add_seats() RETURNS TRIGGER AS
$$
DECLARE
    seat_id INT;
BEGIN
    FOR seat_id IN 1..30
        LOOP
            INSERT INTO reservation (seat_id, trip_id)
            VALUES (seat_id::TEXT || 'A', NEW.id),
                   (seat_id::TEXT || 'B', NEW.id),
                   (seat_id::TEXT || 'C', NEW.id),
                   (seat_id::TEXT || 'D', NEW.id),
                   (seat_id::TEXT || 'E', NEW.id),
                   (seat_id::TEXT || 'F', NEW.id);
        END LOOP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to add seats after a trip is inserted
CREATE TRIGGER add_seats_trigger
    AFTER INSERT
    ON trip
    FOR EACH ROW
EXECUTE FUNCTION add_seats();

-- Prepopulate the tables with some data

-- Add users to the passenger table
INSERT INTO passenger (name)
VALUES ('Arjun Mehta'),
       ('Aarav Sharma'),
       ('Vihaan Patel'),
       ('Reyansh Gupta'),
       ('Ayaan Kumar'),
       ('Krishna Reddy'),
       ('Advait Iyer'),
       ('Ishaan Naidu'),
       ('Dhruv Bhat'),
       ('Arnav Ghosh'),
       ('Viraj Rao'),
       ('Shaan Desai'),
       ('Kabir Joshi'),
       ('Aditya Nair'),
       ('Aryan Kaur'),
       ('Samar Khanna'),
       ('Atharv Verma'),
       ('Nirav Bajaj'),
       ('Omkar Sinha'),
       ('Dev Mishra'),
       ('Ananya Reddy'),
       ('Saanvi Sharma'),
       ('Aadhya Patel'),
       ('Aarohi Gupta'),
       ('Myra Kumar'),
       ('Diya Mehta'),
       ('Ira Bhat'),
       ('Kiara Naidu'),
       ('Inaya Iyer'),
       ('Navya Desai'),
       ('Sara Joshi'),
       ('Anika Rao'),
       ('Ahana Ghosh'),
       ('Prisha Bajaj'),
       ('Riya Nair'),
       ('Tara Khanna'),
       ('Meera Verma'),
       ('Swara Sinha'),
       ('Aadhira Mishra'),
       ('Shanaya Kaur'),
       ('Emma Johnson'),
       ('Liam Smith'),
       ('Olivia Brown'),
       ('Noah Wilson'),
       ('Ava Davis'),
       ('William Moore'),
       ('Sophia Taylor'),
       ('James Anderson'),
       ('Isabella Thomas'),
       ('Benjamin Jackson'),
       ('Lucas White'),
       ('Mason Harris'),
       ('Mia Martin'),
       ('Ethan Thompson'),
       ('Amelia Garcia'),
       ('Alexander Martinez'),
       ('Harper Robinson'),
       ('Henry Clark'),
       ('Evelyn Lewis'),
       ('Michael Lee'),
       ('Abhinav Menon'),
       ('Aarush Pillai'),
       ('Ayaan Shetty'),
       ('Krish Pandey'),
       ('Rohan Chatterjee'),
       ('Vivan Roy'),
       ('Atharva Sengupta'),
       ('Neil Banerjee'),
       ('Avi Nambiar'),
       ('Samar Dutta'),
       ('Ishan Sahu'),
       ('Reyansh Bhattacharya'),
       ('Kian Chakraborty'),
       ('Aarav Sanyal'),
       ('Advait Bose'),
       ('Dhruv Mazumdar'),
       ('Arnav Basu'),
       ('Om Malhotra'),
       ('Vihaan Nag'),
       ('Dev Bhargava'),
       ('Aarohi Bhatt'),
       ('Anika Nanda'),
       ('Saanvi Goel'),
       ('Navya Mathur'),
       ('Ira Grover'),
       ('Kiara Jindal'),
       ('Diya Kochhar'),
       ('Ananya Bansal'),
       ('Myra Chawla'),
       ('Inaya Luthra'),
       ('Prisha Advani'),
       ('Ahana Kapoor'),
       ('Riya Sodhi'),
       ('Swara Sabharwal'),
       ('Aadhya Khatri'),
       ('Meera Gill'),
       ('Sara Vohra'),
       ('Tara Tandon'),
       ('Shanaya Arora'),
       ('Aadhira Chauhan'),
       ('Ved Agarwal'),
       ('Aayush Saxena'),
       ('Arjun Tiwari'),
       ('Krishna Dubey'),
       ('Raghav Srivastava'),
       ('Aarush Kaul'),
       ('Reyansh Sood'),
       ('Dhruv Puri'),
       ('Kabir Kohli'),
       ('Ayaan Batra'),
       ('Vihaan Kohli'),
       ('Ishaan Mehra'),
       ('Nirav Trivedi'),
       ('Dev Thakur'),
       ('Arnav Bajpai'),
       ('Omkar Sood'),
       ('Shaan Chopra'),
       ('Krish Chauhan'),
       ('Aditya Mathur'),
       ('Advait Sahni'),
       ('Kabir Anand'),
       ('Viraj Kapoor'),
       ('Samar Gupta'),
       ('Aryan Menon'),
       ('Krishna Iyer'),
       ('Ananya Nair'),
       ('Saanvi Reddy'),
       ('Aarohi Sharma'),
       ('Diya Patel'),
       ('Myra Gupta'),
       ('Inaya Kumar'),
       ('Navya Mehta'),
       ('Aadhya Bhat'),
       ('Ahana Naidu'),
       ('Riya Iyer'),
       ('Prisha Desai'),
       ('Tara Joshi'),
       ('Kiara Rao'),
       ('Ira Ghosh'),
       ('Aadhira Bajaj'),
       ('Shanaya Nair'),
       ('Sara Khanna'),
       ('Meera Verma'),
       ('Swara Sinha'),
       ('Aarav Singh'),
       ('Anika Rana'),
       ('Aadhya Malhotra'),
       ('Kabir Seth'),
       ('Dev Kumar'),
       ('Ishaan Joshi'),
       ('Dhruv Kapoor'),
       ('Vihaan Gupta'),
       ('Arnav Patel'),
       ('Shaan Mehta'),
       ('Krish Sharma'),
       ('Aditya Jain'),
       ('Viraj Gupta'),
       ('Samar Verma'),
       ('Reyansh Singh'),
       ('Kabir Nair'),
       ('Advait Reddy'),
       ('Krishna Sharma'),
       ('Aarav Verma'),
       ('Ananya Mehta'),
       ('Saanvi Patel'),
       ('Diya Singh'),
       ('Myra Jain'),
       ('Inaya Khanna'),
       ('Navya Sharma'),
       ('Ira Verma'),
       ('Kiara Singh'),
       ('Aadhya Patel'),
       ('Ahana Mehta'),
       ('Riya Gupta'),
       ('Prisha Sharma'),
       ('Tara Singh'),
       ('Meera Jain'),
       ('Swara Patel'),
       ('Sara Mehta'),
       ('Shanaya Verma');

-- Insert Airline Entries
INSERT INTO airline (name)
VALUES ('Air India'),
       ('IndiGo'),
       ('SpiceJet'),
       ('Vistara'),
       ('GoAir');

-- Insert Trip Entries
INSERT INTO trip (airline_id, schedule)
VALUES
-- Airline 1 (Air India)
(1, '2024-06-23 09:00:00'),
(1, '2024-06-24 09:00:00'),
(1, '2024-06-25 09:00:00'),
(1, '2024-06-26 09:00:00'),
(1, '2024-06-27 09:00:00'),
(1, '2024-06-28 09:00:00'),
(1, '2024-06-29 09:00:00'),
(1, '2024-06-30 09:00:00'),
(1, '2024-07-01 09:00:00'),
(1, '2024-07-02 09:00:00'),
(1, '2024-07-03 09:00:00'),
(1, '2024-07-04 09:00:00'),
(1, '2024-07-05 09:00:00'),
(1, '2024-07-06 09:00:00'),
(1, '2024-07-07 09:00:00'),
(1, '2024-07-08 09:00:00'),
(1, '2024-07-09 09:00:00'),
(1, '2024-07-10 09:00:00'),
(1, '2024-07-11 09:00:00'),
(1, '2024-07-12 09:00:00'),
(1, '2024-07-13 09:00:00'),
(1, '2024-07-14 09:00:00'),
(1, '2024-07-15 09:00:00'),
(1, '2024-07-16 09:00:00'),
(1, '2024-07-17 09:00:00'),
-- Airline 2 (IndiGo)
(2, '2024-06-24 12:00:00'),
(2, '2024-06-25 12:00:00'),
(2, '2024-06-26 12:00:00'),
(2, '2024-06-27 12:00:00'),
(2, '2024-06-28 12:00:00'),
(2, '2024-06-29 12:00:00'),
(2, '2024-06-30 12:00:00'),
(2, '2024-07-01 12:00:00'),
(2, '2024-07-02 12:00:00'),
(2, '2024-07-03 12:00:00'),
(2, '2024-07-04 12:00:00'),
(2, '2024-07-05 12:00:00'),
(2, '2024-07-06 12:00:00'),
(2, '2024-07-07 12:00:00'),
(2, '2024-07-08 12:00:00'),
(2, '2024-07-09 12:00:00'),
(2, '2024-07-10 12:00:00'),
(2, '2024-07-11 12:00:00'),
(2, '2024-07-12 12:00:00'),
(2, '2024-07-13 12:00:00'),
(2, '2024-07-14 12:00:00'),
(2, '2024-07-15 12:00:00'),
(2, '2024-07-16 12:00:00'),
(2, '2024-07-17 12:00:00'),
-- Airline 3 (SpiceJet)
(3, '2024-06-24 15:00:00'),
(3, '2024-06-25 15:00:00'),
(3, '2024-06-26 15:00:00'),
(3, '2024-06-27 15:00:00'),
(3, '2024-06-28 15:00:00'),
(3, '2024-06-29 15:00:00'),
(3, '2024-06-30 15:00:00'),
(3, '2024-07-01 15:00:00'),
(3, '2024-07-02 15:00:00'),
(3, '2024-07-03 15:00:00'),
(3, '2024-07-04 15:00:00'),
(3, '2024-07-05 15:00:00'),
(3, '2024-07-06 15:00:00'),
(3, '2024-07-07 15:00:00'),
(3, '2024-07-08 15:00:00'),
(3, '2024-07-09 15:00:00'),
(3, '2024-07-10 15:00:00'),
(3, '2024-07-11 15:00:00'),
(3, '2024-07-12 15:00:00'),
(3, '2024-07-13 15:00:00'),
(3, '2024-07-14 15:00:00'),
(3, '2024-07-15 15:00:00'),
(3, '2024-07-16 15:00:00'),
(3, '2024-07-17 15:00:00'),
-- Airline 4 (Vistara)
(4, '2024-06-26 09:00:00'),
(4, '2024-06-27 09:00:00'),
(4, '2024-06-28 09:00:00'),
(4, '2024-06-29 09:00:00'),
(4, '2024-06-30 09:00:00'),
(4, '2024-07-01 09:00:00'),
(4, '2024-07-02 09:00:00'),
(4, '2024-07-03 09:00:00'),
(4, '2024-07-04 09:00:00'),
(4, '2024-07-05 09:00:00'),
(4, '2024-07-06 09:00:00'),
(4, '2024-07-07 09:00:00'),
(4, '2024-07-08 09:00:00'),
(4, '2024-07-09 09:00:00'),
(4, '2024-07-10 09:00:00'),
(4, '2024-07-11 09:00:00'),
(4, '2024-07-12 09:00:00'),
(4, '2024-07-13 09:00:00'),
(4, '2024-07-14 09:00:00'),
(4, '2024-07-15 09:00:00'),
(4, '2024-07-16 09:00:00'),
(4, '2024-07-17 09:00:00'),
-- Airline 5 (GoAir)
(5, '2024-06-26 12:00:00'),
(5, '2024-06-27 12:00:00'),
(5, '2024-06-28 12:00:00'),
(5, '2024-06-29 12:00:00'),
(5, '2024-06-30 12:00:00'),
(5, '2024-07-01 12:00:00'),
(5, '2024-07-02 12:00:00'),
(5, '2024-07-03 12:00:00'),
(5, '2024-07-04 12:00:00'),
(5, '2024-07-05 12:00:00'),
(5, '2024-07-06 12:00:00'),
(5, '2024-07-07 12:00:00'),
(5, '2024-07-08 12:00:00'),
(5, '2024-07-09 12:00:00'),
(5, '2024-07-10 12:00:00'),
(5, '2024-07-11 12:00:00'),
(5, '2024-07-12 12:00:00'),
(5, '2024-07-13 12:00:00'),
(5, '2024-07-14 12:00:00'),
(5, '2024-07-15 12:00:00'),
(5, '2024-07-16 12:00:00'),
(5, '2024-07-17 12:00:00');