# goauth2.0

### Authorization server following the oauth 2.0 flow

### How to get started

To get started, follow the instructions on how to run the REST API

##### Generate the public and private RSA keys

```bash
make keys
```

##### Load environment variables

```bash
export $(grep -v ^# .env.example)
```

##### Load private rsa key

```bash
export PRIVATE_KEY=$(cat ./scripts/private_key.pem)
```

##### Run DynamoDB migrations

```bash
  make migrations
```

##### Run it!

```bash
make run
```

### How to run http handler test

To get started, follow the instructions on how to run http handler test

##### Generate the public and private RSA keys

```bash
make keys
```

##### Load private rsa key

```bash
export PRIVATE_KEY=$(cat ./scripts/private_key.pem)
```

##### Run command

###### Note: adding values to flags

```bash
go test ./... -v -client_id={client_id} -client_secret={client_secret} -redirect_url={redirect_url}
```

<strong>Note:</strong> You must register your project to get the environment variable needed to set up Google oauth. [Google Cloud Platform](https://cloud.google.com)

##### Oauth2

```
 cookie storage                                                                       server storege
|-------------|   |----------|               Request HTTP             |----------|   |-------------|
|code verifier|   |          | random state && code challenge && code |          |   |    code     |
|     &&      |   |          |=======================================>|          |   |  challenge  |
|    state    |   |          |                                        |Authoriza-|   |             |
|-------------|   |          |          redirected back to            |cion serv-|   |-------------|
                  |          |            state && code               |er        |
                  |  Client  |<=======================================|          |
|-------------|   |          |                                        |  Google  |   |-------------|
|client state |   |          |               Exchange                 |  Login   |   |code verifier|
|   matches   |   |          |        code verifier && code           |          |   |   matches   |
|  redirect   |   |          |=======================================>|          |   |    code     |
|    state    |   |          |                                        |          |   |  challenge  |
|-------------|   |----------|                                        |----------|   |-------------|
```

###### References

- [OAuth 2.0 Simplified](https://www.oauth.com/)
- [ The OAuth 2.0 Authorization Framework
  ](https://datatracker.ietf.org/doc/html/rfc6749)
