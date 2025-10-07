# Driver API v2

The main purpose of the Driver API v2 is to output information about the driver.  
This API has two types of requests – one for a specific driver and one for all drivers.

---

## Request for a specific driver

```
GET /drivers/{driverId}?version=2&api_key=<...>
Host: api.fm-track.com  
Content-Type: application/json;charset=UTF-8
```

**Mandatory parameters**

| Parameter | Type   | Description          |
|-----------|--------|----------------------|
| driverId  | String | Driver identifier    |
| version   | String | Version of the API   |
| api_key   | String | User identification key |

**Response example**

```json
{
  "id": "ABC123",
  "first_name": "Driver",
  "last_name": null,
  "address": "Lithuania, Vilnius",
  "phone": "+3700000000",
  "identifiers": [
    {
      "identifier": "343234323432342",
      "type": "DLT"
    },
    {
      "identifier": "123456789",
      "type": "TACHOGRAPH"
    },
    {
      "identifier": "3AC64785D2FF",
      "type": "WIRELESS"
    },
    {
      "identifier": "123456789",
      "type": "IBUTTON"
    }
  ]
}
```

> **Note:** Parameters without values return a `null` value.

**Response parameters**

| Parameter    | Type   | Array | Description                     | Units |
|---------------|--------|--------|----------------------------------|--------|
| `id`          | String |        | Driver identifier                | Text   |
| `first_name`  | String |        | The driver’s first name          | Text   |
| `last_name`   | String |        | The driver’s last name           | Text   |
| `address`     | String |        | The driver's address             | Text   |
| `phone`       | String |        | The driver's telephone number    | Text   |
| `identifiers` | Array  |        | Container for identification codes |        |
| &nbsp; `identifier` | String |        | Identification code               | Text   |
| &nbsp; `type`       | String |        | Identification code type          | Text   |

Possible values for `type`: `DLT`, `TACHOGRAPH`, `WIRELESS`, `IBUTTON`

---

## Request for all drivers

```
GET /drivers?version=2&api_key=<...>&limit=<...>&continuation_token=<...>&identifier_type=<...>&identifier=<...>
Host: api.fm-track.com  
Content-Type: application/json;charset=UTF-8
```

**Parameters**

| Parameter           | Type   | Description |
|---------------------|--------|-------------|
| version             | String | Version of the API (required) |
| api_key             | String | User identification key (required) |
| limit               | Number | How many drivers should be included in the response <br> Default: 100 <br> Max: 1000 |
| continuation_token  | Number | Token for paginated continuation |
| identifier_type     | String | Identification type, used to filter (possible values: `DLT`, `TACHOGRAPH`, `WIRELESS`, `IBUTTON`) |
| identifier          | String | Identification code, used to filter specific driver |

**Response example**

```json
{
  "count": 100,
  "continuation_token": 123,
  "items": [
    {
      "id": "ABC123",
      "first_name": "Driver",
      "last_name": null,
      "address": "Lithuania, Vilnius",
      "phone": "+3700000000",
      "identifiers": [
        {
          "identifier": "343234323432342",
          "type": "DLT"
        },
        {
          "identifier": "123456789",
          "type": "TACHOGRAPH"
        },
        {
          "identifier": "3AC64785D2FF",
          "type": "WIRELESS"
        },
        {
          "identifier": "123456789",
          "type": "IBUTTON"
        }
      ]
    }
  ]
}
```

> **Note:** Parameters without values return a `null` value.

**Response parameters**

| Parameter           | Type   | Array | Description |
|----------------------|--------|--------|-------------|
| `count`              | Number |        | How many records are included in the response |
| `continuation_token` | Number |        | From which record the data is shown if limit was reached |
| `items`              | Array  |        | Container for all drivers |
| &nbsp; `id`          | String |        | Driver identifier |
| &nbsp; `first_name`  | String |        | The driver’s first name |
| &nbsp; `last_name`   | String |        | The driver’s last name |
| &nbsp; `address`     | String |        | The driver’s address |
| &nbsp; `phone`       | String |        | The driver’s telephone number |
| &nbsp; `identifiers` | Array  |        | Identification codes container |
| &nbsp;&nbsp; `identifier` | String |   | Identification code |
| &nbsp;&nbsp; `type`        | String |   | Identification code type |

Possible values for `type`: `DLT`, `TACHOGRAPH`, `WIRELESS`, `IBUTTON`
