# Payroll application

# External Dependencies
The following are required to build the project

- golang-migrate
- postgres@13
- make

# Running the application
- run `make` to compile the code
- run `make migrate_up`
- to put some data in the database run `make database_seed`
- All the variables can be overridden with make to ensure you can connect to the database
- source `./scripts/dbcredentials.sh` this is a best guess at database credentials

## Endpoints

- __GET__ `api/payroll?userid=<int>&month=<int>`
- __POST__ `api/payroll` requires a payload of:
  `{"forename": <string>, "surname": <string>, "grossSalary": <F64>, "countryCode": <int>}`

## Country Codes
- Italy: `1`
- France: `2` 

# Limitations
- You can get one month at a time for a person, not all months of the year
- Cannot add anymore countries (well you can with SQL, but not through the API)
- Real / Future payroll is not handled whatsoever (would be done predicated on the month)
- I have not used an ORM, mainly due to my unfamiliarity with golang's ecosphere
- There are absolutely no tests, though things have been setup in a way that would make testing fairly straight forward

-----

# Payroll Challenge

Each employee has a yearly gross salary, and on top of that, the company can give a net monthly bonus at his discretion.

Pento is a multinational corporation based in Italy, and it is expanding in France, and the employee can be hired only in one of those two countries, to calculate the net salary each of the countries has a different process:

    Italy:
        The employee has 14 salaries per year; August and December have an extra salary
        The employee pays 25% of taxes
        The employee pays 7% for the national health insurance
    France:
        The employee has 12 salaries per year
        The employee pays 30% of taxes
        The employee does not pay for a national health insurance

Have a payroll summary each month for every employee with the gross/net salary breakdown.

Keep in mind that there are two kinds of payrolls current/past payrolls named "Real Payroll" and future payrolls named "Future Preview".
