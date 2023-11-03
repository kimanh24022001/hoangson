#!/usr/bin/bash

curl -X PUT http://localhost:8080/users -d @user_creatch_batch.json | jq
