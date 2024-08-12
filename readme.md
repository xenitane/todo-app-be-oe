# Todo App Backend

> Backend for todo app with auth written in go

## Features

- As it is written in go, a compiled language. The footprint is small, i.e. low resource usage
- Along with that it is highly performant.
- Built using echo framework which makes it easy to extend and handle errors
- Postgres is used as database.
- Authentication based on JWT.

## Setup

### Requirements for local development

- Have build-essentials like `make` on your system.
- Have Golang-1.22.6 installed on your system, refer [here](https://go.dev)
- Have the `.env` file, format [here](#env-file-format)
- Have Docker, Docker compose installed.

#### Steps

- Clone this repository via `git clone https://github.com/xenitane/todo-app-be-oe` and enter it.
- Execute `make docker-run` to start a database container.
- In a seperate terminal instance execute `make watch` to start the development server.

#### env file format

```sh
PORT=8080

DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=todo-app
DB_USERNAME=user
DB_PASSWORD=pass
DB_SCHEMA=public

JWT_SIGNING_KEY=secret
```

### Containerization and deployment

The project has a `Dockerfile` which can be used to build a portable image for the application.

The image can be build using the `docker build -t <tag-name> .` command.

When using that image for deployment, make sure that you have exposed the necessary environment variables to the container as mentioned in the [env file format](#env-file-format).

## HTTP Endpoints

This is the list of all the http endpoints this application has

| PATH                               | METHOD |              REQUIRED HEADERS               |             REQUEST BODY              | ADMIN ONLY |
| :--------------------------------- | :----: | :-----------------------------------------: | :-----------------------------------: | :--------: |
| /                                  |  GET   |                    none                     |                 none                  |     -      |
| /health                            |  GET   |                    none                     |                 none                  |     -      |
| /api/auth/signup                   |  POST  |                    none                     |      [Signup Body](#signup-body)      |     -      |
| /api/auth/signin                   |  POST  |                    none                     |      [Signin Body](#signin-body)      |     -      |
| /api/user                          |  GET   | "Authorization": "Bearer &lt;jwt token&gt;" |                 none                  |    YES     |
| /api/user/{username}               |  GET   | "Authorization": "Bearer &lt;jwt token&gt;" |                 none                  |     -      |
| /api/user/{username}               | PATCH  | "Authorization": "Bearer &lt;jwt token&gt;" | [User Update Body](#user-update-body) |     -      |
| /api/user/{username}/todo          |  GET   | "Authorization": "Bearer &lt;jwt token&gt;" |                 none                  |     -      |
| /api/user/{username}/todo          |  POST  | "Authorization": "Bearer &lt;jwt token&gt;" |    [Add Todo Body](#add-todo-body)    |     -      |
| /api/user/{username}/todo/{todoid} |  GET   | "Authorization": "Bearer &lt;jwt token&gt;" |                 none                  |     -      |
| /api/user/{username}/todo/{todoid} | DELETE | "Authorization": "Bearer &lt;jwt token&gt;" |                 none                  |     -      |
| /api/user/{username}/todo/{todoid} | PATCH  | "Authorization": "Bearer &lt;jwt token&gt;" | [Todo Update Body](#todo-update-body) |     -      |

#### signup body

```json
{
  "username": "username",
  "firstName": "Jhon",
  "lastName": "Meyr",
  "password": "password"
}
```

#### signin body

```json
{
  "username": "username",
  "password": "password"
}
```

#### user update body

```json
{
  "firstName": "string",
  "lastName": "string",
  "password": "password",
  "isAdmin": true // this field can only be accessed by admins
}
```

#### add todo body

```json
{
  "title": "Do my homework",
  "description": "Complete 2 essays,1 assignment and prepare for a quiz.", // can be empty
  "dueDate": "2024-08-12T03:53:59.000Z" //date in ISO String format
}
```

#### update todo body

```json
{
  "title": "new title",
  "description": "new description",
  "dueDate": "2025-01-01T00:00:00.000Z",
  "status": 2 // status values: [0, pending], [1, work in progress], [2, completed]
}
```
