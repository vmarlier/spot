# SPOT

Single Point of Testing

## TODO

[ ] Decide naming convention for API

[ ] Feature 1
POST request with body describing the request to send to which target and the expected answer.
e.g.:
```txt
POST /v1/request/send
{
    "target": "localhost:8080",
    "path": "/test",
    "body": "{}",
    "expectedAnswer": "200",
}
```

[ ] Feature 2
POST request with body describing the number of requests to send to which target and the expected answers.
e.g.:
```txt
POST /v1/request/sendMass
{
    "target": "localhost:8080",
    "path": "/test",
    "body": "{}",
    "volume": "200",
    "elapsed": "10ms",
    "expectedBehaviour": {
        "whileHealthyStatusCode": 200,
        "whileUnhealthyStatusCode": "4xx, 5xx",
    }
}
```

[ ] Create a GET help handler to receive all available endpoints
[ ] Create a health handler
