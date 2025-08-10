#!/bin/bash

API_KEY="mx0vglWtzbBOGF34or"
SECRET="77658a3144bd469fa8050b9c91b9cd4e"
recvWindow=5000
timestamp=$(($(date +%s%3N)))
query="recvWindow=$recvWindow&timestamp=$timestamp"
signature=$(echo -n "$query" | openssl dgst -sha256 -hmac "$SECRET" | sed 's/^.* //')

curl -s -X GET "https://api.mexc.com/api/v3/account?$query&signature=$signature" \
  -H "X-MEXC-APIKEY: $API_KEY"
