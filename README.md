# ibmctoken

ibmctoken gets IBM Cloud Oauth access token for provided API key.

Usage

```go
apiKey := os.Getenv("MY_API_KEY")
token := ibmctoken.NewToken(apiKey)
_ = token.RequestToken()
fmt.Println(token.AccessToken)
```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
