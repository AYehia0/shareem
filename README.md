# Shareem

Ignore the name, the purpose of this app is to try a new deployment method!


## Running locally
Make sure you have docker and docker-compose installed in your system, simply run `docker compose build` to build the images, then `docker compose up` to run! by default hot reload works!

### Migration
To create migration, make sure [golang-migrate](https://github.com/golang-migrate/migrate) is installed : `migrate create -ext sql -dir database/migrations -seq <migration_name>`

The migrations run automatically when app is restarted!

### Sqlc
To generate the sql code, make sure [sqlc](https://github.com/sqlc-dev/sqlc) is installed: `sqlc generate`

sqlc uses the config file in : `sqlc.yaml`, you can see that I am using `pgx/v4` not `v5`
