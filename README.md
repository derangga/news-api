# [News API](https://github.com/derangga/news-api)

A simple CRUD news and topic management.

# Database schema

The project has 4 tables: User, Topics, News Articles, and News Topics

Generated from [dbdiagrams](https://dbdiagram.io/d/news-api-68429013ba2a4ac57b1c0116)

<img src="docs/news-api-table.png" />
<br><br/>

# Development

## Pre-requisite

- Install Python3
- PostgreSQL
- [Golang migrate](https://github.com/golang-migrate/migrate)
- [Golang mock](https://github.com/golang/mock)

## Api Documentation

Before develop any new endpoints, ensure you write the API documentation.

1. Navigate to `/api-docs` dir, update the `news-api.yaml`
2. Run the API Docs by using python `http.server`.

```bash
$ cd api-docs/
$ python3 -m http.server 8081
```

## Development

1. Make `.env` file

```bash
cp env.sample .env
```

2. Start the database service. If you choose using docker, you can use `docker-compose` configuration

```bash
docker compose up
```

3. Push the database schema

```bash
make migrate_up
```

4. Run the project

```bash
make run_api
```

5. To update `/mocks` file for unit test you can run

```bash
make mocks_gen
```
