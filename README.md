# Twomes Backoffice API
API for Twomes data collection platform to enable research.

## Table of contents
- [Deploying](#deploying)
- [Developing](#developing)
- [Usage](#usage)
- [Status](#status)
- [License](#license)
- [Credits](#credits)

## Deploying
For our process to deploy the API to our public server, or update it, see these links:
- Deploy: https://github.com/energietransitie/twomes-backoffice-configuration#api
- Update: https://github.com/energietransitie/twomes-backoffice-configuration#updating

### Prerequisites
The Twomes API is available as a Docker image.
You will need to [install Docker](https://docs.docker.com/engine/install/) to run it.

### Images
See all [available images](https://github.com/energietransitie/twomes-backoffice-api/pkgs/container/twomes-backoffice-api):
- Use the `latest` tag to get the latest stable release built from a tagged GitHub release. 
- Use the `main` tag to get the latest development release, built directly from the `main` branch.

### Docker Compose ([more information](https://docs.docker.com/compose/features-uses/))
```yaml
version: "3.8"
services:
  web:
    container_name: twomes-api-web
    image: ghcr.io/energietransitie/twomes-backoffice-api:latest
    ports:
      - 8080:8080
    volumes:
      - /path/to/data:/data
    environment:
      - TWOMES_DSN=root:password@tcp(db:3306)/twomes
      - TWOMES_BASE_URL=http://localhost:8080
    depends_on:
      - db

  db:
    container_name: twomes-api-db
    image: mariadb:latest
    environment:
      - MYSQL_DATABASE=twomes
      - MYSQL_ROOT_PASSWORD=password
```

## Developing

### Requirements
- [Go (minimum 1.20)](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)

### Running
Make sure Docker is running on your local machine, then start the service from a command line terminal from the root of this repository:
```shell
docker compose up --build
```

This generates log messages from both the `web` and the `db` component.
Just keep this running in your terminal.

The API is now available on http://localhost:8080/.

Create a new admin account to use the admin endpoints:
```shell
docker compose exec -i web twomes-backoffice-api admin create -n <name>
```
> Substitute `<name>` with the name of the admin account you want to create.

Example output of running the command, with `johndoe` as admin name:
```text
Admin "johndoe" created. Authorization token: eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJUd29tZXNBUEl2MiIsInN1YiI6IjQiLCJleHAiOjE3MTU3MzEyMDAsIm5iZiI6MTY4NDE1MjA4OSwiaWF0IjoxNjg0MTUyMDg5LCJraW5kIjoiYWRtaW5Ub2tlbiJ9.N_uhPhLsaTq0DVGVPhdfU6Hd2VD0Zb8QxesTaWeILlNkkjQ9Vuxpwe0sfi3Vj0GJgyin2ZilPE6AS-makGm2cg
```

Copy the authorization token without any spaces.

When finished, type the `Ctrl + C` key in the first terminal. The container state is 
preserved, and to restart, simply run `docker compose up --build` again.

To completely remove all docker containers created above:
```shell
docker compose rm
```

To delete the saved data, remove the `data` directory in the root of this repository.

### Folder structure

This repository tries to implement a DDD approach. While some elements are still too tightly coupled to really call it DDD, the structure still tries te represent DDD as best as possible.

| Folder       | Purpose                                                                   |
| ------------ | ------------------------------------------------------------------------- |
| .github      | GitHub Actions workflows and config files                                 |
| cmd          | CLI commands                                                              |
| docs         | Additional documentation                                                  |
| handlers     | HTTP handlers for API endpoints                                           |
| internal     | Utitilities that are not exposed outside of this package                  |
| repositories | Repositories for domain models. Contains al DB logic                      |
| services     | Services tie repositories and subservices together to perform operations. |
| swaggerdocs  | Swagger UI and OpenAPI spec.                                              |
| twomes       | Domain models and logic.                                                  |

### Model diagram

To re-generate the model diagram:
```shell
docker run -i --rm vranac/erd < docs/model.er > docs/model.pdf
```

## Usage

### Tokens
Most end points require a session token to be provided in an authorization
bearer HTTP header. These end points are marked with a 'lock' symbol. Click
on the 'Authorize' button at the upper right of the page, or click on the 
'lock' symbol at the end point, and paste the session token in the value field.
Subsequent calls done through http://localhost:8080/docs will then use the
session token.

There are currently three types of tokens:
- Admin: Used by administrators to manage resources.
- Account: Used by an account to manage its resources.
- Device: Used by a measurement device to upload measurements.

### Managing admins and cloudfeeds
When the container is running, lookup it's name.

Run the following command to see info about how to manage admins:
```shell
docker exec <container-name> twomes-backoffice-api admin --help
```

Run the following command to see info about how to manage cloudfeeds:
```shell
docker exec <container-name> twomes-backoffice-api cloudfeed --help
```

### Administrators on our servers
Contact an administrator to get admin access to the API:
- Nick van Ravenzwaaij
- Henri ter Hofte

## Status
Project is: _work in progress_

## License
This software is available under the [Apache 2.0 license](./LICENSE), 
Copyright 2021 [Research group Energy Transition, Windesheim University of 
Applied Sciences](https://windesheim.nl/energietransitie) 

## Credits
This software is created by:
- Nick van Ravenzwaaij · [@n-vr](https://github.com/n-vr)

Thanks also go to:
- Arjan peddemors · [@arpe](https://github.com/arpe)

Product owner:
- Henri ter Hofte · [@henriterhofte](https://github.com/henriterhofte)

We use and gratefully aknowlegde the efforts of the makers of the following source code and libraries:

- [chi](https://github.com/go-chi/chi), by Peter Kieltyka, Google Inc, licensed under [MIT license](https://github.com/go-chi/chi/blob/master/LICENSE)
- [gorm](https://gorm.io), by Jinzhu, licensed under [MIT license](https://github.com/go-gorm/gorm/blob/master/License)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql), by go-sql-driver, licensed under [Mozilla Public License 2.0](https://github.com/go-sql-driver/mysql/blob/master/LICENSE)
- [logrus](https://github.com/sirupsen/logrus), by Simon Eskildsen, licensed under [MIT license](https://github.com/sirupsen/logrus/blob/master/LICENSE)
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt), by Dave Grijalva, licensed under [MIT license](https://github.com/golang-jwt/jwt/blob/main/LICENSE)
- [crc16](https://github.com/sigurn/crc16), by sigurn, licensed under [MIT license](https://github.com/sigurn/crc16/blob/master/LICENSE)
- [swagger-ui](https://github.com/swagger-api/swagger-ui), by SmartBear Software Inc., licensed under [Apache 2.0 license](https://github.com/swagger-api/swagger-ui/blob/master/LICENSE)
