# Twomes Backoffice API

Twomes API based on FastAPI / SQLAlchemy / MariaDB 

## Table of contents

- [Prerequisites](#prerequisites)
- [Local Service](#local-service)
- [Deployment](#deployment)
- [Development](#development)
- [Status](#status)
- [License](#license)
- [Credits](#credits)


## Prerequisites

Running, and further developing, the API requires a recent Docker setup.
See https://www.docker.com/products/docker-desktop for installation.


## Local Service

To try out the Twomes API, locally, on your machine, it is possible to run 
the database and API server using Docker Compose. 

Start the service from the command line, from the root directory 
of this project.
```shell
docker-compose up
```

This generates log messages from both the `web` and the `db` component.
Just keep this running in your terminal.

Open another terminal, and fill the database with some initial data (e.g., 
Property and DeviceType instances). See also `src/test/fixture.py`.
```shell
docker-compose run web python -c "from test.fixture import base; base()"
```

The API is now available on http://localhost:8000/ , but it is more convenient
to try it out via http://localhost:8000/docs . Note: the `device create` API
end point takes a `device_type` parameter as input. This is the name of one
of the device types defined in `src/test/fixture.py` - 'Gateway' for example.

Most end points require a session token to be provided in an authorization
bearer HTTP header. These end points are marked with a 'lock' symbol. Click
on the 'Authorize' button at the upper right of the page, or click on the 
'lock' symbol at the end point, and paste the session token in the value field.
Subsequent calls done through http://localhost:8000/docs will then use the
session token.

There are currently three types of session tokens:
- Admin: for creating accounts and devices, used by Twomes admins
- Account: for account activation, and attaching devices to accounts
- Device: for uploading device measurements

To add yourself, locally, as one of the administrator:
```shell
docker-compose run web python -c "import user; user.create_admin()"
```
Type your name, and add the admin tuple in `src/user.py`.

Example output of running `create_admin()`, providing `piet` as admin name:
```text
Admin name: piet
Update user.admins with this tuple:
  (2, 'piet', '$2b$12$3wMWc1PK4OqCWuFZGF5XieCTOFBbP6uBTZHtkc9vCFRlYUZciXOuu')
Authorisation bearer token for admin "piet":
  Mg.u6Rcx2fHl-lydbEiKZILGtd9i1hzCES1uXkcPFT-tw0
```

When finished, type Ctrl-C in the first terminal. The container state is 
preserved, and to restart, simply run `docker-compose up` again.

To completely remove all docker containers created above
```shell
docker-compose rm
```

## Deployment

Deployment of the Twomes API is done using the Docker image created after
every commit into the `master` branch. The image is pushed to the Github
docker registry (part of Github Packages).

The image is available at
```text
docker.pkg.github.com/energietransitie/twomes-backoffice-api/api:latest
```

## Development

### Setup

Create a python virtualenv with a recent python version (>= 3.9), and 
activate this virtual environment.

From within your virtualenv, install the required packages.
```shell
$ pip install -r requirements.txt
```

Run a MariaDB server on your local machine, e.g.:
```shell
docker run -p 127.0.0.1:3306:3306  --name mariadb -e MYSQL_DATABASE=twomes -e MYSQL_ROOT_PASSWORD=twomes -d mariadb:10.5.9
```

Set the MariaDB url in your environment:
```shell
export TWOMES_DB_URL="root:twomes@localhost/twomes"
```

### Database migrations

Database migrations are managed with Alembic (see https://alembic.sqlalchemy.org).

To get an up-to-date database schema, at any time, run
```shell
PYTHONPATH=src alembic upgrade head
```

If you change the database models, create migrations
```shell
PYTHONPATH=src alembic revision --autogenerate -m "<short description of migration>"
```

The new revision is stored in `alembic/versions`: check that the newly 
created revision file is correct, and commit to your development branch.
Then run the `alembic upgrade` command again.


## Status

Project is: _work in progress_


## License

This software is available under the [Apache 2.0 license](./LICENSE), 
Copyright 2021 [Research group Energy Transition, Windesheim University of 
Applied Sciences](https://windesheim.nl/energietransitie) 


## Credits

To be done