# FreeStuff API Go Client

This is a client for the [FreeStuff API](https://docs.freestuffbot.xyz/), written in Go.

Refer to `internal/pkg/client/v1/cli/main.go` for example usage of it.

## Testing

You can execute unit tests by executing the following:

```
make test
```

### Integration Testing

You can verify the ability of this client to interact with the live FreeStuff API by executing the following:

```
FREESTUFF_API_KEY="<your key here>" go run internal/pkg/client/v1/cli/main.go
```