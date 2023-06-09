# Muerta

Muerta - RESTful API for a term paper on "Web application to control the shelf life of products using computer vision".

## How to run?

First, create an `.env` file and put the following environment variables in it:

```shell
PORT=[port]
API_NAME=[name]
DB_NAME=[db_name]
DB_PORT=[db_port]
DB_HOST=[db_host]
DB_PASSWORD=[db_password]
DB_USER=[db_user]
ALLOWED_ORIGINS=[host1.com,host2.org]
CACHE_HOST=[redis_host]
CACHE_USER=[redis_username]
CACHE_PASSWORD=[redis_password]
CACHE_PORT=[redis_port]
```

Then Start the Docker containers with this command:

```shell
docker compose up -d
```

> Make sure you have open ports for the API and Database

## Features

- [x] Service to recognize shelf life in text from picture
- [x] JWT Authentication
- [x] Logging in JSON format
- [x] Users with roles
- [x] Swagger API documentation
- [x] Redis for caching JWT tokens
