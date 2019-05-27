API overview
============

## Contents

- [API overview](#api-overview)
  - [Contents](#contents)
  - [Endpoints](#endpoints)
    - [Accounts](#accounts)
      - [Create an account:](#create-an-account)
      - [Fetching accounts:](#fetching-accounts)
      - [Fetching an account's operations:](#fetching-an-accounts-operations)
    - [Operations](#operations)
      - [Make deposit](#make-deposit)
      - [Make deposit](#make-deposit-1)
  - [Entities](#entities)
    - [Account](#account)
    - [Operation](#operation)
      - [Operation type](#operation-type)
      - [Transaction](#transaction)




## Endpoints

### Accounts

#### Create an account:

    POST /accounts

Body request:

| Attribute  | Description                 |
| ---------- | --------------------------- |
| `id`       | The ID of the account       |
| `name`     | The username of the account |
| `currency` | The currency of the account |

Creates and returns a new [Account](#account),

#### Fetching accounts:

    GET /accounts

Returns list of exists [Accounts](#account).

#### Fetching an account's operations:

    GET /accounts/{id}/operations

Returns a list of account [operations](#operation).

### Operations

#### Make deposit

    POST /operations/deposit

Body request:

| Attribute  | Description                   |
| ---------- | ----------------------------- |
| `to`       | The ID of the account         |
| `currency` | The currency of the operation |
| `amount`   | Amount of the operation       |

Creates and returns new deposit [operation](#operation) for an account.

#### Make deposit

    POST /operations/transfer

Body request:

| Attribute  | Description                   |
| ---------- | ----------------------------- |
| from       | Account - donor               |
| to         | Account - recipient           |
| `currency` | The currency of the operation |
| `amount`   | Amount of the operation       |

Creates and returns new deposit [operation](#operation) for an account.

## Entities

### Account

| Attribute  | Description                 |
| ---------- | --------------------------- |
| `Id`       | The ID of the account       |
| `Name`     | The username of the account |
| `Currency` | The currency of the account |
| `Amount`   | Amount of the account       |

### Operation
Simple entity for description operations between accounts.

| Attribute      | Description                  |
| -------------- | ---------------------------- |
| `Id`           | The ID of the operation      |
| `Participants` | Accounts ids of participants |
| `Type`         | The type of the operation    |
| `Transactions` | List of related transactions |

#### Operation type
|Value| Description|
| 0 | Deposit type. Used to send money to the account from the outside world|
| 1 | Transfer type. Used to transfer money between accounts|

#### Transaction
Low-level entity for describing operations between 2 accounts or an account and the world.


| Attribute     | Description                 |
| ------------- | --------------------------- |
| `OperationID` | Operation id                |
| `From`        | Account - donor             |
| `To`          | Account - recipient         |
| `Currency`    | Currency of the transaction |
| `Amount`      | Amount of the transaction   |
