# zk-map-server

## Build

```shell
go build -o zk-map-server main.go
```

## Running

```shell
./zk-map-server --config config/zk-map-server-example.yaml
```

# API docs

## Get Proof

get proof by chain id and block height

### Request Path

/proof

### Request Method

GET

### Request Parameters

| Parameters | Type   | Explanation  |
|------------|--------|--------------|
| chain_id   | number | chain id     |
| height     | string | block height |

### Response Parameters

| Parameters | Type   | Explanation |
|------------|--------|-------------|
| code       | int    |             |
| msg        | int    |             |
| data       | object |             |

### Data structure

| Parameters | Type   | Explanation              |
|------------|--------|--------------------------|
| height     | string |                          |
| status     | number | 1: Pending, 3: Completed |
| result     | object |                          |
| error_msg  | string |                          |

### Result structure

| Parameters   | Type         | Explanation |
|--------------|--------------|-------------|
| proof        | object       |             |
| public_input | string array |             |

### Example

#### Request:

```shell        
curl --location 'http://127.0.0.1:8181/proof?chain_id=212&height=60000'
```

#### Response:

Status Pending

```shell

{
    "code": 2000,
    "msg": "Success",
    "data": {
        "height": "60000",
        "status": 1,
        "result": {},
        "error_msg": ""
    }
}
```

Status Completed

```shell
{
    "code": 2000,
    "msg": "Success",
    "data": {
        "height": "60000",
        "status": 3,
        "result": {
            "proof": {
                "pi_a": [
                    "15454543049305956149566727234832623020128796206286670808542742536256596816758",
                    "17529398676163506942832940247284184970670446222082998632049949597725571761940",
                    "1"
                ],
                "pi_b": [
                    [
                        "17141476719120535317401571072921657417091535445907074814178180109440650100061",
                        "15155667921174618731178702110360034957523986495831724415253203584773871159156"
                    ],
                    [
                        "11859746840569313754971312058373580852456923226139956097374129164583068595319",
                        "20365816851161492615463491357423417314245381049064245238157074921122522877015"
                    ],
                    [
                        "1",
                        "0"
                    ]
                ],
                "pi_c": [
                    "15293573022955680275234953682558526729573946116942012059117642796867751491327",
                    "13937753267051210259367129412169820039326429849451621423430959370595357654639",
                    "1"
                ],
                "protocol": "groth16"
            },
            "public_input": [
                "12869250466664455014759573801630797168758698215523897593611610845197028483503"
            ]
        },
        "error_msg": ""
    }
}

```    