# github-search [![CircleCI](https://circleci.com/gh/hasyimibhar/github-search/tree/master.svg?style=svg)](https://circleci.com/gh/hasyimibhar/github-search/tree/master)

A sample application using Github v3 [Search API](https://developer.github.com/v3/search/).

## Requirements

- Go 1.11

## Usage

To start the API server:

```sh
$ go run . /path/to/config.yml
```

## Setting up database

1. Make sure you have MySQL or MariaDB running. The easiest way is to use docker:

```sh
# You might need to wait a few seconds for the database to setup before proceeding with the next steps
$ docker run --rm -d \
    -e MYSQL_ROOT_PASSWORD=password \
    -p 3306:3306 \
    --name mariadb \
    mariadb:10.1
```

2. Create the database:

```sh
$ mysql -uroot -h127.0.0.1 -p -e "create database github_search;"
```

3. Run the migrations:

```sh
$ docker run \
    -v $(pwd)/database/migrations:/migrations \
    --network host \
    migrate/migrate -path /migrations -database "mysql://github_search:password@tcp(127.0.0.1:3306)/github_search" up
```

## Configuration

```yaml
http_port: 8081
log_level: trace
cors_enabled: true

database:
  host: 127.0.0.1
  database: github_search
  user: github_search
  password: password

# Generate this here: https://github.com/settings/applications/new
github:
  client_id: XXXX
  client_secret: XXXX
```

## Tests

To run tests:

```sh
# You don't need to set these, but without them, the tests might fail sometimes due to rate limiting
$ export GITHUB_CLIENT_ID=XXXX
$ export GITHUB_CLIENT_SECRET=XXXX
$ go test ./...
```
