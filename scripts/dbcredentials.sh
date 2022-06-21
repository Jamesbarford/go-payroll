#!/usr/bin/env bash

setDevelopmentDatabaseCredentials() {
  export DB_USER=""
  export DB_PASSWORD=""
  # grab user
  export DB_NAME=$(whoami)
  export DB_HOST="localhost"
  export DB_PORT=5432
}

setDevelopmentDatabaseCredentials
