# PDFMonkey Golang SDK

This is a simple SDK for the [PDFMonkey](https://pdfmonkey.io) API written in golang.

# Usage

```go
	client, err := pdfmonkeysdkgo.NewClient()
	if err != nil {
		panic(err)
	}

	var c *pdfmonkeysdkgo.GetCurrentUserResponse

	c, err = client.GetCurrentUser()
	if err != nil {
        panic(err)
	}
	fmt.Printf("Fetched user with name: %s\n\n", c.CurrentUser.DesiredName)
```

## Client customization

The client can be customized by passing API key and API endpoint:

```go
	client, err := pdfmonkeysdkgo.NewClient(WithAPIKey("my-api-key"), WithEndpoint("https://custom.api.endpoint"))
	if err != nil {
		panic(err)
	}

```

The API key and the endpoint can be set using environment variables:
- PDFMONKEY_API_KEY
- PDFMONKEY_API_ENDPOINT

The precedence for the client configuration is as follow:
1. Environment variables
2. Options functions (WithAPIKey, WithEndpoint)
3. Default values

# TODO

- [ ] write test suite
- [ ] tidy up the code
