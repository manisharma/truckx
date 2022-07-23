CREATE TABLE IF NOT EXISTS roles (
   id serial PRIMARY KEY,
   name VARCHAR (255) UNIQUE NOT NULL
);

-- Insert 2 roles
INSERT INTO roles(name) VALUES ('SuperUser');
INSERT INTO roles(name) VALUES ('EndUser');

CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name VARCHAR (100) NOT NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL,
    phone_no VARCHAR (15),
    role_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- create a SuperUser
INSERT INTO users(name, email, password, phone_no, role_id) VALUES ('TruckX Admin', 'admin@truckx.com', '$2a$08$6POj4kowAo0B23opQA/GMux/khdkwtcvBJwDDAVS/kGMm600yiuIO', '9876352345', 1);

CREATE TABLE IF NOT EXISTS theatres (
   id serial PRIMARY KEY,
	name VARCHAR (100) NOT NULL,
   address VARCHAR (255) NOT NULL,
	state VARCHAR (50) NOT NULL,
	pin_code INT NOT NULL,
	total_seats INT NOT NULL
);
-- create a Theatre
INSERT INTO theatres(name, address, state, pin_code, total_seats) VALUES ('Mohini Thretre', 'Borivali', 'Maharashtra', '400066', 5);

CREATE TABLE IF NOT EXISTS seats (
   id serial PRIMARY KEY,
   theatre_id INT NOT NULL,
   FOREIGN KEY (theatre_id) REFERENCES theatres(id)
);
-- create 5 seats in the Theatre
INSERT INTO seats(theatre_id) VALUES (1);
INSERT INTO seats(theatre_id) VALUES (1);
INSERT INTO seats(theatre_id) VALUES (1);
INSERT INTO seats(theatre_id) VALUES (1);
INSERT INTO seats(theatre_id) VALUES (1);

CREATE TABLE IF NOT EXISTS shows (
   id serial PRIMARY KEY,
   theatre_id INT NOT NULL,
   from_hr INT NOT NULL,
   from_min INT NOT NULL,
   till_hr INT NOT NULL,
   till_min INT NOT NULL,
   FOREIGN KEY (theatre_id) REFERENCES theatres(id),
   CONSTRAINT Unq_Theatre_Slot UNIQUE(theatre_id, from_hr, from_min, till_hr, till_min)
);
-- Theatre will host 4 shows eveyday
INSERT INTO shows(theatre_id, from_hr, from_min, till_hr, till_min) VALUES (1, 9, 0, 12, 0);
INSERT INTO shows(theatre_id, from_hr, from_min, till_hr, till_min) VALUES (1, 12, 30, 15, 30);
INSERT INTO shows(theatre_id, from_hr, from_min, till_hr, till_min) VALUES (1, 16, 0, 19, 0);
INSERT INTO shows(theatre_id, from_hr, from_min, till_hr, till_min) VALUES (1, 19, 30, 22, 30);


CREATE TABLE movies (
   id serial PRIMARY KEY,
   name VARCHAR (255) NOT NULL,
   release_date TIMESTAMP NOT NULL
);
-- create 2 movies
INSERT INTO movies(name, release_date) VALUES ('Hera Pheri', CURRENT_TIMESTAMP);
INSERT INTO movies(name, release_date) VALUES ('Phir Hera Pheri', CURRENT_DATE + INTERVAL '1 day');


CREATE TABLE IF NOT EXISTS movies_shows (
   id serial PRIMARY KEY,
   movie_id INT NOT NULL,
   show_id INT NOT NULL,
   from_date TIMESTAMP NOT NULL,
   till_date TIMESTAMP NOT NULL,
   FOREIGN KEY (movie_id) REFERENCES movies(id),
   FOREIGN KEY (show_id) REFERENCES shows(Id),
   CONSTRAINT Unq_Movie_Show UNIQUE(movie_id, show_id, from_date, till_date)  
);
INSERT INTO movies_shows(movie_id, show_id, from_date, till_date) VALUES (1, 1, CURRENT_DATE, CURRENT_DATE + INTERVAL '90 day');
INSERT INTO movies_shows(movie_id, show_id, from_date, till_date) VALUES (1, 2, CURRENT_DATE, CURRENT_DATE + INTERVAL '90 day');
INSERT INTO movies_shows(movie_id, show_id, from_date, till_date) VALUES (1, 3, CURRENT_DATE, CURRENT_DATE + INTERVAL '90 day');
INSERT INTO movies_shows(movie_id, show_id, from_date, till_date) VALUES (1, 4, CURRENT_DATE, CURRENT_DATE + INTERVAL '90 day');

CREATE TABLE IF NOT EXISTS bookings (
   id serial PRIMARY KEY,
   movies_shows_id INT NOT NULL,
   creation_date TIMESTAMP NOT NULL,
   booking_date TIMESTAMP NOT NULL,
   FOREIGN KEY (movies_shows_id) REFERENCES movies_shows(id)
);

CREATE TABLE bookings_seats (
   booking_id INT NOT NULL,
   seat_id INT NOT NULL,
   guest_name VARCHAR (100),
   guest_email VARCHAR (255),
   guest_phone VARCHAR (15),
   FOREIGN KEY (booking_id) REFERENCES bookings(id),
   FOREIGN KEY (seat_id) REFERENCES seats(id)
);
CREATE INDEX bookings_seats_booking_id_idx ON bookings_seats(booking_id);

CREATE TABLE user_bookings (
   user_id INT NOT NULL,
   booking_id INT NOT NULL,
   FOREIGN KEY (user_id) REFERENCES users(id),
   FOREIGN KEY (booking_id) REFERENCES bookings(id)
);
CREATE INDEX user_bookings_user_id_idx ON user_bookings(user_id);
CREATE INDEX user_bookings_booking_id_idx ON user_bookings(booking_id);