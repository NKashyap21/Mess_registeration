CREATE TABLE IF NOT EXISTS mess (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS mess_capacities (
	id SERIAL PRIMARY KEY,
	mess_id INT NOT NULL CONSTRAINT fk_mess_capacity_mess REFERENCES mess(id),
	total_capacity INT NOT NULL,
	from_month INT NOT NULL,
	from_year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS registrants (
	id SERIAL PRIMARY KEY,
	email TEXT NOT NULL,
	mess_id INT NOT NULL CONSTRAINT fk_registrant_mess REFERENCES mess(id),
	month INT NOT NULL,
	year INT NOT NULL
); 

-- Redis is used for rest of the active data
