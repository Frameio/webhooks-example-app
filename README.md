# webhooks-example-app

This is an example application which shows how build a Webhook receiver for Frame.io.

For more information about how Webhooks work, check out our [documentation](https://docs.frame.io/docs/webhooks).

## Usage

```
$ go run main.go
```

## Tests

```
$ go test
```

By default the application runs on port 3000, however this can be overrided by setting the `PORT` environment variable.

The service will need to run at a publically accessible address. For quick testing, you can use [ngrok](https://ngrok.com/) or similar service.
