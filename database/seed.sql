BEGIN;

INSERT INTO users(forename, surname, country_code) VALUES
  ('Alice', 'Smith', 1),
  ('Bob', 'Thompson', 1),
  ('Json', 'DeRulo', 2),
  ('Brian', 'Kernighan', 2);

COMMIT;

BEGIN;
INSERT INTO countries(id, country_name) VALUES
  (1, 'France'),
  (2, 'Italy');

INSERT INTO payroll(gross_salary, date_payable, user_id) VALUES
  (60000, '2022-06-22 05:00:00-07', 1), -- Alice
  (90000, '2022-06-22 05:00:00-07', 2), -- Bob
  (40000, '2022-06-22 05:00:00-07', 3), -- Json
  (20000, '2022-06-22 05:00:00-07', 4); -- Brian

INSERT INTO taxes(percent,country_code) VALUES
  (30, 2), -- FRANCE
  (25, 1); -- ITALY

INSERT INTO deductions(country_code, percentage, reason) VALUES
  (1, 7, 'national insurance');

/* Bonuses for countries */
INSERT INTO country_bonuses(country_code, percentage, month_payable) VALUES
  (1, 100, 8),
  (1, 100, 12);
COMMIT;
