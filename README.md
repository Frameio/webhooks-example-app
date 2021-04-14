<img width="1644" alt="artboard_small" src="https://user-images.githubusercontent.com/19295862/66240171-ba8dd280-e6b0-11e9-9ccf-573a4fc5961f.png">

# Frame.io Webhook Receiver in Go
Frame.io is a cloud-based collaboration hub that allows video professionals to share files, comment on clips real-time, and compare different versions and edits of a clip. 

This is an example application which shows how to build a webhook receiver for Frame.io. You can use webhooks to receive notifications about events that occur in the Frame.io app. These notifications can be sent to external systems for processing, API callback, and workflow automation. 

For more information about how Webhooks work, check out our [documentation](https://docs.frame.io/docs/webhooks).

## Pre-requisites 

* Developer account with Frame.io - [https://developer.frame.io](https://developer.frame.io/)
* Web server with a publicly accessible address - we recommend trying [ngrok](https://ngrok.com/)

## Usage

```
$ go run main.go
```

## Tests

```
$ go test
```

If the test is working, you should receive something that looks similar to this: 

```
2019/10/11 13:39:03 POST http://example.com/ping HTTP/1.1
Content-Type: application/json
X-Frameio-Request-Timestamp: 1570826343
X-Frameio-Signature: v0=9aab38e9622907c1683429ea491ce6f291c4fb7404c4b55089eeca824465a454

{"id":"123","name":"ping","resource":null}
```
At the end of the test, you'll get a note saying `PASS` in the terminal. 

By default the application runs on port 3000, however this can be overridden by setting the `PORT` environment variable.

The service will need to run at a publicly accessible address. For quick testing, you can use [ngrok](https://ngrok.com/) or similar service. If you need help setting up or troubleshooting ngrok, see [How to Setup and Troubleshoot ngrok (Mac)](https://docs.frame.io/docs/how-to-setup-and-troubleshoot-ngrok-mac).
