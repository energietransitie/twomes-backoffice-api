# Twomes Backoffice API

Twomes API based on FastAPI / SQLAlchemy / MariaDB 


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

Some end points require a session token to be provided in an authorization
bearer HTTP header. These end points are marked with a 'lock' symbol. Click
on the 'Authorize' button, and paste the session token in the value field.
Subsequent calls done through http://localhost:8000/docs will then use the
session token.

When finished, type Ctrl-C in the first terminal. The container state is 
preserved, and to restart, simply run `docker-compose up` again.

To completely remove all docker containers created above
```shell
docker-compose rm
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


