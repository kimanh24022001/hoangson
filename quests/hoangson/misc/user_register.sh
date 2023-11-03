#!/usr/bin/bash

curl -X PUT http://localhost:8080/users/register -d '{"email": "abc@smatyx.com"}' | jq