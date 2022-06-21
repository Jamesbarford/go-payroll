/* Assumes the database is representative of one company
 * otherwise there would need to be a company table etc...
 */

BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  forename VARCHAR(200) NOT NULL,
  surname VARCHAR(200) NOT NULL,
  country_code INT  NOT NULL,
  UNIQUE(forename, surname)
);

/* Country id to pretty name */
CREATE TABLE IF NOT EXISTS countries (
  id INT UNIQUE NOT NULL,
  -- Are there countries with the same name? I do not know
  country_name VARCHAR(200),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS payroll (
  id SERIAL PRIMARY KEY,
  user_id INT UNIQUE,
  gross_salary DOUBLE PRECISION NOT NULL,
  -- Nuke salary from database if user is deleted ... questionable if we'd want this?
  CONSTRAINT fk_user
   FOREIGN KEY(user_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);

/* Idea being every country has some form of tax */
CREATE TABLE IF NOT EXISTS taxes (
  id SERIAL PRIMARY KEY,
  percent DOUBLE PRECISION NOT NULL,
  country_code INT UNIQUE NOT NULL
);

/* Allows for additional deductions, `n` per country */
CREATE TABLE IF NOT EXISTS deductions (
  id SERIAL PRIMARY KEY,
  country_code INT NOT NULL,
  percentage DOUBLE PRECISION,
  reason VARCHAR(200)
);

/* There could also be no reason why we might not have a user <> bonus personal structure
 * i.e for investment banking */
CREATE TABLE IF NOT EXISTS country_bonuses (
  id SERIAL PRIMARY KEY,
  percentage DOUBLE PRECISION,
  month_payable INT NOT NULL,
  country_code INT NOT NULL
);

CREATE INDEX deductions_country_code_idx ON deductions (country_code);
CREATE INDEX user_country_code_idx ON users(country_code);
CREATE UNIQUE INDEX tax_country_code_idx ON taxes(country_code);
CREATE INDEX bonuses_country_code_idx ON country_bonuses(country_code);

COMMIT;
