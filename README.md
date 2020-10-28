## Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Environment](#environment)
* [Contact](#contact)

## About The Project

Games Provider API is a 4G Solution. Responsible of managing information emails.

Following main tasks accomplished by the microservice:

* Managing Models
* Mananing Model
* Retrieve list of model
* Update tag with audit information
* Retrieve Mogels with the desired fields
* Create Many Models at once


### Built With

* [Redis](https://redis.io/)
* [Golang](https://golang.org/)
* [MySQL](https://dev.mysql.com/doc/relnotes/mysql/8.0/en/)


## Getting Started

### Prerequisites

Software you need to install to be ready to Go!

* [Golang with Chocolatey](https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10-es)

* [MySQL with Docker](https://phoenixnap.com/kb/mysql-docker-container)

### Installation

1. Locate your workspace

```sh
cd ~$HOME/github.com/briams/
```

2. Clone the repo

```sh
git clone http://github.com/briams/4g-events-api.git
```

3. Install Go packages

```sh
go install .
```

4. Setup environment file

5. Start program

  ```sh
  go run ./cmd/api
  ```

  - Building the binary
  ```sh
  go build ./cmd/api && ./api
  ```

## Environment
| Variable               | Ejemplo                                  | Informacion adicional    |
|------------------------|------------------------------------------|--------------------------|
| GENERAL_API_KEY        | %9*e%c86B57dQNH42X%*9qQ%X$q369          | LLave de autorización    |
| APP_STAGE              | DEV                                      | Stage del microservicio  |
| APP_PORT               | 3040                                     | Puerto del microservicio |
| SERVER_ALLOW_ORIGINS | https://gdpteam.com,http://localhost:3040 | Dominios permitidos por CORS (dejar vacío para aceptar cualquiera) |
| ROUTES_SERVER_TIMEOUT_MINS | 4                                | Timeout del server |
| DB_CONNECTION       |  mysql                                   | Nombre del motor de BD   |
| DB_DATABASE         | 4GEmailingDB                               | Nombre de la BD          |
| DB_HOST             | 172.16.10.134                            | BD Host               |
| DB_PORT             | 3306                                     | BD Port               |
| DB_USERNAME         | root                                     | BD User                  |
| DB_PASSWORD         | 6gKLm6GsYUqheJX1CcGDqwgy.q07gTjg98iIoRiV | BD Password             |
| DB_MAXOPEN          | 5                                        | Open connections         |
| DB_MAXIDLE          | 5                                        | Max Iddle Connections    |
| DB_MAXAGEMINS       | 60                                       | Max age min Connections  |
| REDIS_HOST          | 172.16.10.134                            | Host de Redis            |
| REDIS_PORT          | 6378                                     | Puerto de Redis          |
| REDIS_PASSWORD      |                                          | Redis password (por defecto vacío)   |
| REDIS_DB            |                                           | Nombre de redis BD (por defecto vacío) |

