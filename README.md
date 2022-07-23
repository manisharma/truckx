# Theatre Booking Engine

Theatre Booking Engine mimics a minimal movie booking `server`.
The `Engine` use `JWT` to secure sensitive routes.
The `Engine` has its dependency on `postgres` database. To facilitate exploring the `database` there's another sidekick `pgadmin`

> `Schema` internal/db/migrations/000001_schema.up.sql

## Endpoints

```json
{
    "requests": [
        {
            "name": "movies _avalibility",
            "url": "http://localhost:8080/api/v1/movies/2/availability?state=Maharashtra",
            "method": "GET",
            "params": [
                {
                    "name": "state",
                    "value": "Maharashtra",
                    "isPath": false
                }
            ],
            "body": {
                "type": "json",
                "raw": "{\n    \"email\": \"admin@truckx.com\",\n    \"password\": \"passw0rd\"\n}",
                "form": []
            },
            "auth": {
                "type": "bearer",
                "bearer": "<token>"
            }
        },
        {
            "name": "login",
            "url": "http://localhost:8080/api/login",
            "method": "POST",
            "body": {
                "type": "json",
                "raw": "{\n    \"email\": \"admin@truckx.com\",\n    \"password\": \"passw0rd\"\n}",
                "form": []
            },
            "auth": {
                "type": "bearer",
                "bearer": "<token>"
            }
        },
        {
            "name": "get_bookings_by_id",
            "url": "http://localhost:8080/api/v1/bookings/1",
            "method": "GET",
            "auth": {
                "type": "bearer",
                "bearer": "<token>"
            }
        },
        {
            "name": "get_bookings_by_user",
            "url": "http://localhost:8080/api/v1/bookings",
            "method": "POST",
            "body": {
                "type": "json",
                "raw": "{\n    \"theatreId\": 1,\n    \"state\": \"Maharashtra\",\n    \"movieId\": 1,\n    \"movieShowId\": 3,\n    \"seats\": [{\n        \"seat_id\":\"1\",\n        \"guest_name\": \"guest 1\",\n        \"guest_email\": \"guest1@yahoo.com\",\n        \"guest_phone\": \"8765423456\"\n    },{\n        \"seat_id\":\"2\",\n        \"guest_name\": \"guest 2\",\n        \"guest_email\": \"guest2@yahoo.com\",\n        \"guest_phone\": \"4565234567\"\n    },{\n        \"seat_id\":\"3\",\n        \"guest_name\": \"guest3\",\n        \"guest_email\": \"guest3@yahoo.com\",\n        \"guest_phone\": \"7563452456\"\n    },{\n        \"seat_id\":\"4\",\n        \"guest_name\": \"guest4\",\n        \"guest_email\": \"guest4@yahoo.com\",\n        \"guest_phone\": \"7653456789\"\n    }],\n    \"date\": \"2022-07-23\"\n}",
                "form": []
            },
            "auth": {
                "type": "bearer",
                "bearer": "<token>"
            }
        },
        {
            "name": "get_all_bookings_for_a_user",
            "url": "http://localhost:8080/api/v1/bookings",
            "method": "GET",
            "auth": {
                "type": "bearer",
                "bearer": "<token>"
            }
        },
        {
            "name": "movies",
            "url": "http://localhost:8080/api/v1/movies",
            "method": "GET",
            "body": {
                "type": "json",
                "raw": "{\n    \"email\": \"admin@truckx.com\",\n    \"password\": \"passw0rd\"\n}",
                "form": []
            },
            "auth": {
                "type": "bearer",
                "bearer": "<token>"
            }
        }
    ]
}
```


## Prerequisites
The `Engine` at bare minimim needs `docker` installed/running and a `.env` under `/cmd` with following details
```
DB_HOST=truckx_postgres
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=passw0rd
DB_DATABASE=postgres
DB_EXPLORER_EMAIL=admin@truckx.com
DB_EXPLORER_PASSWORD=passw0rd
JWT_KEY=super_secret_key
``` 

## Execution

> `make up` will bring up all services.

After all services have successfully come up visit
[localhost](http://localhost:8080) and star playing with the endpoints.
To explore the postgres data visit [postgres](http://localhost:5050) and login with `DB_EXPLORER_EMAIL` and `DB_EXPLORER_PASSWORD` and configure the database in there with the said configuration from `.env` file.

> `make migrate` to run migration

> `make down` will teardown all services
