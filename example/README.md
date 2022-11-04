Example demo of how to use all of these modules.

## Start the app
Using a proper `NEW_RELIC_KEY`, run the following:
```shell
env \
  SERVICE_NAME=foobar \
  ENV=local \
  TS_LOG_LEVEL=info \
  TS_LOG_FORMAT=json \
  NEW_RELIC_FORWARD_LOGS=true \
  NEW_RELIC_KEY=foobar \
    go run main.go
```

## Hit the endpoint
```shell
curl localhost:8081/ok
```