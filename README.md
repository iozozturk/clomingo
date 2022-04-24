# CLOMINGO

Minimum Viable Cloud Native GO Microservice to start your own

Best suitable for your REST, Mobile or Single Page Application(SPA) backends

Tried to showcase best practices about project structure, dependency injection, config management, logging etc. Let me
know if you think something need to be changed via issues.

Enjoy! Or buy me a coffee and we enjoy together :)

<a href="https://www.buymeacoffee.com/ismet" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>

### What is this repository for? ###

* Bootstrap your Golang Microservice on Google Cloud
* Does not have many dependencies, you can integrate your own choices along the way
* Uses native `net/http`
* Uses `Zap` from Uber for logging
* Uses `Google Cloud Datastore` for database connection
* Uses `Google Sign In` for authentication flow
* Uses `App Engine` for deployment
* Contains `Dockerfile` for deployment to other environments such as `Cloud Run, Kubernetes`
* Contains `Makefile` for building and deploying
* Includes sample requests in `requests.http` file that you can easily run with Intellij IDEs
* No frontend implemented
* REST backend for all kind of your SPA or mobile applications

### Implemented routines/endpoints ###

* social sign in - Google Sign In
* logout
* healthcheck
* sample middlewares for authentication and logging

### How do I get set up? ###

* Download the project
* Setup your Google Cloud environment
* Make changes at config_env.json if necessary
* run `go run ./cmd` command at root directory
* hit to localhost:8080

### Sample Requests ###

####  Sign In With Google

```
$ curl -H "Content-Type: application/json" -H "X-Installation-Id: some-unique-id" -X POST -d '{
                                                                                  "token":"googe sign in token captured",
                                                                                  "pushToken": "push notification token if you have one"
                                                                                }' http://localhost:8080/auth/google
```

##### Response:

```
{
  "Id": 5643280054222748,
  "SessionToken": "961253fa-7082-11ec-b758-acde48001122",
  "SocialToken": "googe sign in token captured",
  "SessionType": 1,
  "UserId": 5704568653556992,
  "User": {
    "Id": 5704568653556992,
    "Name": "user name",
    "Email": "user-emai@gmail.com",
    "Photo": "https://lh3.googleusercontent.com/a-/blabla=s96-c",
    "Timestamp": "2022-01-07T20:00:08.291634Z",
    "Timeupdate": "2022-01-07T20:00:08.291634Z"
  },
  "UserAgent": "Apache-HttpClient/4.5.13 (Java/11.0.13)",
  "UserIp": "127.0.0.1:55580",
  "DeviceId": "idea-device-id",
  "PushToken": "push-notification-token",
  "PushEnabled": true,
  "Timestamp": "2022-01-08T13:57:50.748715+01:00",
  "Timeupdate": "2022-01-08T13:57:50.748715+01:00"
}
```

#### Sign Out

Cookies should be preserved between requests. After sign-out cookies will be cleared automatically.

```
curl -X POST --location "http://localhost:8080/auth/signout" \
    -H "Content-Type: application/json" \
    -H "X-Installation-Id: idea-device-id"
```
