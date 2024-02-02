# goauth2.0

### Authorization server following the oauth 2.0 flow

### How to get started

To get started, follow the instructions on how to run the REST API

##### Load environment variables

```bash
export $(grep -v ^# .env.example)
```

##### Run it!

```bash
make run
```

<strong>Note:</strong> You must register your project to get the environment variable needed to set up Google oauth.  [Google Cloud Platform](https://cloud.google.com)