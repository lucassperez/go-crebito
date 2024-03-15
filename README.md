# Rinha

<img alt="(Image) Cockfight with Go logo on their heads." src="https://github.com/lucassperez/go-crebito/assets/60318892/507339c9-e4b1-49df-be05-4ecb112a2fc1" />

This is my take at the second edition of <a target="_blank" href="https://github.com/zanfranceschi/rinha-de-backend-2024-q1">Rinha de Backend</a>,
which would roughly translate to Backend Fight (or Cockfight).

Basically, it is a challenge to create an app with specific endpoints that can
handle tests the creators of the challenge have prepared. The theme of this
Backend Fight is to handle multiple data reads and writes to the same resource
and how to avoid data races and data inconsistency.

The stack has to consist at least of:
- Load balancer
- Two instances of the server
- Database

I tackled it using Go, a language I'm just learning, and thus decided to try and
learn more on how to create and use middlewares and the new functionalities of
the standard `serveMux` introduced in go version `1.22.0`. Also, I'm just trying
to think of a directory structure to organize everything.

## Tech Stack

The idea was to use only Go's standard library, but I used a driver for the
database as well.

- **Nginx** as load balancer
- **Go** for the api with the driver <a target="_blank" href="https://github.com/jackc/pgx">`pgx/v4`</a>
- **Postgres** as the database

## How To
I have made two different docker images, the `dev` and the `prod` images, using
Docker's really cool <a target="_blank" href="https://docs.docker.com/build/building/multi-stage/">multi stage builds</a>.
<br />
The difference is that the `dev` also uses <a target="_blank" href="https://github.com/cosmtrek/air">air</a>
for live reloading of the app, making the dev experience better, while the `prod`
image just compiles and run the binary.

To help development and testing, make was also used as a command runner
(not as a build tool). Run `make help` to see the options.

You don't have to have go installed, as you can run everything inside docker containers.

For dev, you can just run:
```sh
docker compose up
# or
make server
```

Now you'll have the server accepting connections at `localhost:4000`.
<br />
By default the port will be 4000, but you can change it in `docker-compose.yml`
by changing the `SERVER_ADDRESS` env var. Don't forget to change it in the
`ports` section, too!

To seed the database (You may want to do this before testing the endpoints.):
```sh
make seed
```

To test the endpoints, there is a convenience script located in `scripts/curl-endpoints.sh`.
<br/>
You'll need curl installed to use it.
<br/>
It also features a help option:
```sh
./scripts/curl-endpoints.sh help
```

## Gatling Test

To run the Gatling tests, you'll need to <a target="_blank" href="https://gatling.io/open-source/">install Gatling</a>.
<br />
To install it, all you need to do is download, unzip it and set an environment variable
named `GATLING_HOME` to the directory after unzipping it. This environment variable
is used by this repo's test script (load-test/start.sh).

Then you should stop the running dev containers, if any, and run the prod containers.
To do so:
```sh
make down
make prod.up
```

To start the tests, run:
```sh
make prod.gatling
```

You can see some stats and logs with:
```sh
make prod.stats
# Services names in prod are api1, api2, db and nginx
docker compose -f docker-compose-prod.yml logs <service_name> -f
```

It will place the results in `load-test/user-files/results/`.

Here is an example ran in 2024/02/23:

<img alt="(Image) Gatling tests statistics, 2024/02/23." src="https://github.com/lucassperez/go-crebito/assets/60318892/145855fd-79cf-42be-879b-e2e41d24a016" />

## Results

On the 14th of March of 2024, the results were published. Every contender started
with a certain amount of points and some criteria could deduct points. They were
as follows:

- If at least 98% of the requests were responded with less than 250ms, you'd get
no deductions. If your api had less than 98%, you'd lose points based on how many
% of requests were responded on 250ms or more.
- If the balance were not consistent, you'd lose some points based on how many
inconsistencies the system had.

Inconsistencies would discount more points than speed of requests, meaning that
consistency was more important.

This API had a grand total of 0 penalties! Meaning it passed with max score. 🎉

This talks a lot about Golang's efficiency. Also, most of the work was done in
Postgresql handling the locks.

Results can be seen <a taget="_blank" href="https://github.com/zanfranceschi/rinha-de-backend-2024-q1/blob/main/RESULTADOS.md">here</a>.

Access this page and search for `lucassperez-go-crebito`. Yay! (:

## Endpoints

### GET /clientes/{id}/extrato

#### Params
Url query param `id`: The id of the `cliente`.

#### Response

`HTTP 200 OK`
```json
{
  "saldo": {
    "total": -9098,
    "data_extrato": "2024-01-17T02:34:41.217753Z",
    "limite": 100000
  },
  "ultimas_transacoes": [
    {
      "valor": 10,
      "tipo": "c",
      "descricao": "descricao",
      "realizada_em": "2024-01-17T02:34:38.543030Z"
    },
    {
      "valor": 90000,
      "tipo": "d",
      "descricao": "descricao",
      "realizada_em": "2024-01-17T02:34:38.543030Z"
    }
  ]
}
```

- `saldo`
    - `total` client's balance
    - `data_extrato` timestamp of when the request was made
    - `limite` client's limit
- `ultimas_transacoes` client's last transactions
    - `valor` value of transaction
    - `tipo` type of transaction (`c` is credit, `d` is debit)
    - `descricao` string between 1 and 10 characters
    - `realizada_em` timestamp of when transaction was made

<hr />

### POST /clientes/{id}/transacoes

#### Params
Url query param `id`: The id of the `cliente`.

Body:
```json
{
    "valor": 1000,
    "tipo" : "c",
    "descricao" : "descricao"
}
```

- `valor` must be a **non negative integer**
- `tipo` one of `"c"` (credit) or `"d"` (debit)
- `descricao` a string **between 1 and 10** characters

#### Response

`HTTP 200 OK`

```json
{
    "limite" : 100000,
    "saldo" : -9098
}
```

- `limite` client's limit
- `saldo` client's new balance

Balance can never be less than -limit.
