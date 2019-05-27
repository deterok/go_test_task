## What this
This is the solution of a test task to obtain positions Golang Developer.
Problem: We have users. They have money. They need to manage their funds (for example, to make a transfer). This implementation of a classic solution to managing user accounts. The microservice allows to manage accounts. It allows you to perform various operations with money. Now only two operations are supported: Transfer and deposit.


## How it works
I use three main entities:
- Account: user wallet with only one type of currency
- Operation: set of various actions (transactions) over various Accounts
- Transaction: atomic action on a specific pair of accounts.

When an operation is created, then all participating accounts change the Amount field to the specified number depending on the type of operation and transaction values. Then the operation along with all transactions is saved. Now this operation will be part of the operation history of each account.

## Dependencies
- go-1.*
- docker-18.*
- make-4.2.*

## How to run
- Simple run:
    ```shell
    $ make up
    ```
    or
    ```shell
    $ make up-build
    ```

    After this command, the environment will start and you can request the server at `http://localhost:8800`

- Run tests:

    WARNING: The Command is not optimized. Extra containers can be running.
    ```shell
    $ make test
    ```
    Out:
    ```
    go test  -timeout 60s -v -p 1 ./...
    Starting go_test_task_redis_1    ... done
    Starting go_test_task_postgres_1 ... done
    ?       github.com/deterok/go_test_task/payments/cmd    [no test files]
    ?       github.com/deterok/go_test_task/payments/cmd/service    [no test files]
    ?       github.com/deterok/go_test_task/payments/pkg/endpoint   [no test files]
    ?       github.com/deterok/go_test_task/payments/pkg/http       [no test files]
    === RUN   Test_basicPaymentsService_CreateAccount
    === RUN   Test_basicPaymentsService_CreateAccount/simple_create
    --- PASS: Test_basicPaymentsService_CreateAccount (0.10s)
        --- PASS: Test_basicPaymentsService_CreateAccount/simple_create (0.10s)
    === RUN   Test_basicPaymentsService_GetAccount
    === RUN   Test_basicPaymentsService_GetAccount/simple_getting
    === RUN   Test_basicPaymentsService_GetAccount/another_simple_getting
    === RUN   Test_basicPaymentsService_GetAccount/account_doesn't_exist
    --- PASS: Test_basicPaymentsService_GetAccount (0.30s)
        --- PASS: Test_basicPaymentsService_GetAccount/simple_getting (0.10s)
        --- PASS: Test_basicPaymentsService_GetAccount/another_simple_getting (0.10s)
        --- PASS: Test_basicPaymentsService_GetAccount/account_doesn't_exist (0.10s)
    === RUN   Test_basicPaymentsService_GetAccounts
    === RUN   Test_basicPaymentsService_GetAccounts/simple_getting
    --- PASS: Test_basicPaymentsService_GetAccounts (0.10s)
        --- PASS: Test_basicPaymentsService_GetAccounts/simple_getting (0.10s)
    === RUN   Test_basicPaymentsService_MakeTransfer
    === RUN   Test_basicPaymentsService_MakeTransfer/simple_transfer
    === RUN   Test_basicPaymentsService_MakeTransfer/full_transfer
    --- PASS: Test_basicPaymentsService_MakeTransfer (0.21s)
        --- PASS: Test_basicPaymentsService_MakeTransfer/simple_transfer (0.10s)
        --- PASS: Test_basicPaymentsService_MakeTransfer/full_transfer (0.11s)
    PASS
    ok      github.com/deterok/go_test_task/payments/pkg/service    0.716s
    ```

- To cleanup, run the following command:

    WARNING: It's full cleanup! It can delete important containers (like postgres)!
    ```shell
    make down
    ```

## Examples requests

### Create account
Request:
```shell
curl --request POST \
  --url http://localhost:8800/accounts \
  --header 'Accept: */*' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'Host: localhost:8800' \
  --header 'accept-encoding: gzip, deflate' \
  --header 'cache-control: no-cache' \
  --header 'content-length: 56' \
  --data '{\n	"currency": "USD",\n	"name": "test"\n}'
```

Response:
```
{
    "account": {
        "ID": 1,
        "Name": "test",
        "Currency": "USD",
        "Amount": "0",
        "CreatedAt": "2019-05-28T10:55:17.158888875Z",
        "UpdatedAt": "2019-05-28T10:55:17.158888875Z",
        "DeletedAt": null
    }
}
```

### Create deposite operation
Request:
```shell
curl --request POST \
  --url http://localhost:8800/operations/deposite \
  --header 'Accept: */*' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'Host: localhost:8800' \
  --header 'accept-encoding: gzip, deflate' \
  --header 'cache-control: no-cache' \
  --header 'content-length: 46' \
  --data '{\n    "to": 1,\n    "currency": "USD",\n    "amount": "12"\n}'
```

Response:
```
{
    {
    "operation": {
        "ID": 5,
        "CreatedAt": "2019-05-28T10:56:10.140422707Z",
        "UpdatedAt": "2019-05-28T10:56:10.140422707Z",
        "DeletedAt": null,
        "Participants": [
            -1,
            1
        ],
        "Type": 0,
        "Transactions": [
            {
                "ID": 5,
                "CreatedAt": "2019-05-28T10:56:10.14143324Z",
                "UpdatedAt": "2019-05-28T10:56:10.14143324Z",
                "DeletedAt": null,
                "OperationID": 5,
                "From": -1,
                "To": 1,
                "Currency": "USD",
                "Amount": "12"
            }
        ]
    }
}
}
```

### Fetching operations list of an account
Request:
```shell
curl --request GET \
  --url http://localhost:8800/accounts/1/operations \
  --header 'Accept: */*' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'Host: localhost:8800' \
  --header 'accept-encoding: gzip, deflate' \
  --header 'cache-control: no-cache'
```
Response:

```
{
    "operations": [
        {
            "ID": 5,
            "CreatedAt": "2019-05-28T10:55:22.019653Z",
            "UpdatedAt": "2019-05-28T10:55:22.019653Z",
            "DeletedAt": null,
            "Participants": [
                -1,
                1
            ],
            "Type": 0,
            "Transactions": [
                {
                    "ID": 1,
                    "CreatedAt": "2019-05-28T10:55:22.028264Z",
                    "UpdatedAt": "2019-05-28T10:55:22.028264Z",
                    "DeletedAt": null,
                    "OperationID": 1,
                    "From": -1,
                    "To": 1,
                    "Currency": "USD",
                    "Amount": "12"
                }
            ]
        },
    ]
}
```

## API
You can view the API [there](docs/API.md).


## Directions for improvement
* Separate models and entities
* Add verification of incoming data
* Add better tests (get away from postgresql in tests)
* Add various checks, for example: checking for the existence of currencies
* Add a normal world account
* Add several world accounts for different tasks (taxes, fee, etc.)
