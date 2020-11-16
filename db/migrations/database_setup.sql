CREATE EXTENSION pgcrypto;

CREATE TABLE subscriptions (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  amount NUMERIC NOT NULL,
  date_due DATE NOT NULL 
);

CREATE TABLE users(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name TEXT NOT NULL,
email TEXT NOT NULL
);