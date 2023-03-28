# Twomes Backoffice API

Twomes API based on FastAPI / SQLAlchemy / MariaDB 

## Table of contents

- [Prerequisites](#prerequisites)
- [Deploying](#deploying)
- [Developing](#developing)
- [Status](#status)
- [License](#license)
- [Credits](#credits)


## Prerequisites

Running, and further developing, the API requires a recent Docker setup.
See https://www.docker.com/products/docker-desktop for installation.


## Deploying 

The Twomes API can be be deployed on a local test server or a server in the cloud.

### Deploying on your local machine

To try out the Twomes API, locally, on your machine, it is possible to run 
the database and API server using Docker Compose. 

Make sure Docker is running on your local machine, then start the service from a command line terminal, from the root directory 
of this project.
```shell
docker-compose up
```

This generates log messages from both the `web` and the `db` component.
Just keep this running in your terminal.

Open another terminal, and fill the database with initial Property and 
DeviceType instances. See also `src/data/loader.py`. This loads the data
available in `src/data/sensors.csv`.
```shell
docker-compose run web python -c "from data.loader import csv_create_update; csv_create_update()"
```

The API is now available on http://localhost:8000/ , but it is more convenient
to try it out via http://localhost:8000/docs . Note: the `device create` API
end point takes a `device_type` parameter as input. This is the name of one
of the device types defined in `src/data/sensors.csv` - 'OpenTherm-Monitor' 
for example.

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

To add yourself, locally, as one of the administrator, first pull to make sure that your main branch is up do date and then:
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

When finished, type the `Ctrl + C` key in the first terminal. The container state is 
preserved, and to restart, simply run `docker-compose up` again.

To completely remove all docker containers created above
```shell
docker-compose rm
```

### Deploying to `api.tst.energietransitiewindesheim.nl`

Deployment of the Twomes API to `api.tst.energietransitiewindesheim.nl` is done automatically using the Docker image created after
every commit into the `main` branch. The image is pushed to the Github
docker registry (part of Github Packages).

The image is available at
```text
docker.pkg.github.com/energietransitie/twomes-backoffice-api/api:latest
```
To deploy, recreate the container for `api.tst.energietransitiewindesheim.nl` while using the 'Pull latest image' option.

### Creating new admin accounts to `api.tst.energietransitiewindesheim.nl`
If you need the ability to create user accounts for testing purposes, first follow the procedure to create an admin account as described under [Deploying to a local test server](#deploying-to a-local-test-server) and test it locally via http://localhost:8000/docs. Then commit and push changes in the main branch to origin and ask the admin for the `api.tst.energietransitiewindesheim.nl` server (currently [@henriterhofte](https://github.com/henriterhofte)) to activate your newly created admin account. He will then recreate the container for `api.tst.energietransitiewindesheim.nl` while using the 'Pull latest image' option, which activates the new account.

### Registering new properties to `api.tst.energietransitiewindesheim.nl`
If you need a new property at `api.tst.energietransitiewindesheim.nl`, first update it in `src/data/sensors.csv` and test it locally as described under [Deploying on your local machine](#deploying-on-yourlocal-machine) and test it locally via http://localhost:8000/docs. Commit changes in a separate branch, push and create a Pull Request; ask admin for the `api.tst.energietransitiewindesheim.nl` server (currently [@henriterhofte](https://github.com/henriterhofte)) to review, merge and activate the new definitions of devices and properties in sensors.csv. 

He will then log in via SSH to the [Twomes backoffice server](https://github.com/energietransitie/twomes-backoffice-configuration) at energietransitiewindesheim.nl  and execute the following command:
```shell
docker pull ghcr.io/energietransitie/twomes_api:latest && \
cd /root/api/tst && \
docker-compose up -d && \
docker exec twomes-api-tst python3 -c "from data.loader import csv_create_update; csv_create_update()"
```

## Developing

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

### Running

Run the API server on your local machine:
```shell
PYTHONPATH=src uvicorn api:app --reload
```

Open http://127.0.0.1:8000/docs in your browser, to see the API documentation.


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


### Model diagram

To re-generate the model diagram
```shell
docker run -i --rm vranac/erd < docs/model.er > docs/model.pdf 
```


## Status

Project is: _work in progress_


## License

This software is available under the [Apache 2.0 license](./LICENSE), 
Copyright 2021 [Research group Energy Transition, Windesheim University of 
Applied Sciences](https://windesheim.nl/energietransitie) 


## Credits

This software is created by:
* Nick van Ravenzwaaij  ·  [@n-vr](https://github.com/n-vr)

Thanks also go to:
* Arjan peddemors  ·  [@arpe](https://github.com/arpe)

Product owner:
* Henri ter Hofte  ·  [@henriterhofte](https://github.com/henriterhofte)

We use and gratefully aknowlegde the efforts of the makers of the following source code and libraries:

* [chi](https://github.com/go-chi/chi), by Peter Kieltyka, Google Inc, licensed under [MIT license](https://github.com/go-chi/chi/blob/master/LICENSE)
* [gorm](https://gorm.io), by Jinzhu, licensed under [MIT license](https://github.com/go-gorm/gorm/blob/master/License)
* [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql), by go-sql-driver, licensed under [Mozilla Public License 2.0](https://github.com/go-sql-driver/mysql/blob/master/LICENSE)
* [logrus](https://github.com/sirupsen/logrus), by Simon Eskildsen, licensed under [MIT license](https://github.com/sirupsen/logrus/blob/master/LICENSE)
* [golang-jwt/jwt](https://github.com/golang-jwt/jwt), by Dave Grijalva, licensed under [MIT license](https://github.com/golang-jwt/jwt/blob/main/LICENSE)
