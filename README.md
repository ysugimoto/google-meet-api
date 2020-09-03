# google-meet-api

Yet another Google Meet URL generation API.

## Installation

```
go get github.com/ysugimoto/google-meet-api
```

## Usage

This packages's interface is similar to google api library.

```go
import (
    "log"
    "context"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    meet "github.com/ysugimoto/google-meet-api/v1"
)

func main() {
    // make oauth config from expected JSON file or GOOGLE_APPLICATION_CREDENTIALS or some way
    config, err := google.ConfigFromJSON("credential.json")
    if err != nil {
        log.Fatalln(err)
    }
    ctx := context.Background()

    // Make oauth token from suitable way and initialize client
    token := ...
    client := config.Client(ctx, token)

    // Api use
    m, err := meet.New(client)
    if err != nil {
        log.Fatalln(err)
    }
    resp, err := m.Conference.Create().Context(ctx).Do()
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Google Meet URL Created: %s\n", resp.URL)
}
```

See [example](https://github.com/ysugimoto/google-meet-api/blob/master/example) in detail.

## How this works

Now Google Meet management API is not in API [Product Index](https://developers.google.com/products), But we can create Meet URL by accessing Google Calendar API. This package wraps thier API calling and provide simple interface like google api packages.

## Testing

```shell
go test ./...
```

## License

MIT

## Author

Yoshiaki Sugimoto
