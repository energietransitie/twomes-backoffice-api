# Twomes Backoffice API

Twomes API based on FastAPI / SQLAlchemy / MariaDB 


## Prerequisites

Running, and further developing, the API requires a recent Docker setup.
See https://www.docker.com/products/docker-desktop for installation.


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

