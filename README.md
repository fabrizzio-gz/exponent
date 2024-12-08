# exponent

Send push notifications to Expo apps using Golang

revised version of: https://github.com/oliveroneill/exponent-server-sdk-golang/tree/master

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/9ssi7/exponent.svg)](https://pkg.go.dev/github.com/9ssi7/exponent)

## Installation
```
go get github.com/9ssi7/exponent
```

## Usage
```go
package main

import (
	"context"
	"time"

	"github.com/9ssi7/exponent"
)

func main() {
	c := exponent.NewClient(exponent.WithAccessToken("your-access-token"))

	tkn := exponent.MustParseToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := c.PublishSingle(ctx, &exponent.Message{
		To:       []*exponent.Token{tkn},
		Body:     "This is a test notification",
		Data:     exponent.Data{"withSome": "data"},
		Sound:    "default",
		Title:    "Notification Title",
		Priority: exponent.DefaultPriority,
	})

	if err != nil {
		panic(err)
	}

	if res.IsOk() {
		println("Notification sent successfully")
	} else {
		println("Notification failed")
	}
}

```

## Contributing

We welcome contributions! Please see our [Contribution Guidelines](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License. See [LICENSE](LICENSE) for more details.
