CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE readings (
    id SERIAL PRIMARY KEY,
    device_id INT REFERENCES devices(id),
    value NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    device_id INT REFERENCES devices(id),
    threshold NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE alerts (
    id SERIAL PRIMARY KEY,
    device_id INT REFERENCES devices(id),
    reading_id INT REFERENCES readings(id),
    message TEXT NOT NULL,
    serviced BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now()
);
