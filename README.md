# short-to-me
Web service exposing URL shortening functions written in Golang

License: [MIT](https://opensource.org/licenses/MIT)

![golang_sherlock](https://encrypted-tbn0.gstatic.com/images?q=tbn%3AANd9GcR9JlE4ZofSFGEjKgIEW8dPHrlLg9YuO0Ve0a0VlVlkFBzKUS5B)

## Requirements
In order to correctly build and run the service, the following tools are required:
* Git - [Download and Install](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* Go (v1.15) - [Download and Install](https://golang.org/doc/install)
    * Make sure to install the latest version of the go tools  
    `go get -u golang.org/x/tools/...`
* Mongo (v3.0+) - [Download and Install](https://docs.mongodb.com/manual/installation/)

## Installation
The service can be installed using the terminal, in two ways:
* <b>Cloning the repository</b>:  
run the command `git clone github.com/gsiragusa/short-to-me`  
(it is preferable to run the command from inside your `$GOPATH` folder)
* <b>Using Go Get:</b>  
run the command `go get github.com/gsiragusa/short-to-me`

A folder named `short-to-me` should now be in your workspace.  
Once the service is installed, make sure to download the necessary dependencies.  
The project uses Go Modules, therefore change directory to `short-to-me` and run the command:  

`go mod tidy`

### Build and Run
In order to build the project, change directory to `short-to-me` and run the following command:  

`go build -o short-to-me cmd/short-to-me/main.go`

An executable file `short-to-me` will be generated in your project folder.
  
Before running the service, you can set the following environment variables in order to customize the port in which the service will run, and Mongo parameters. Here are the default values:
```
PORT=8081
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=short-to-me
```

You should be ready to run the service now!  
Run the executable file: `./short-to-me`  
Logs should be visible in your console and opening http://localhost:8081/ from your browser should display a `404` error message.  
You're all set!

#### Tests
The project includes some test files. If you wish to run them, from the same folder run the command:  

`go test ./... -tags=integration`

It will run the tests and show their output on the console.

## Documentation
With the service running on your machine, a Swagger providing all the endpoints specifications can be found at [this address](http://localhost:8081/docs/swagger-ui/)

## Examples
Use [Postman](https://www.postman.com/) or `curl` to:
 
#### Generate a short url
`curl -X POST "http://localhost:8081/api?url=www.google.com" -H "accept: application/json"`

Sample response
```
{
    "status": "ok",
    "operation": "create",
    "url": "http://localhost:8081/pRA4OEy"
}
```

#### Read a short url
`curl -X GET "http://localhost:8081/api?url=http%3A%2F%2Flocalhost%3A8081%2FpRA4OEy" -H "accept: application/json"`

Sample response
```
{
    "status": "ok",
    "operation": "read",
    "url":"http://www.google.com"
}
```

#### Delete a short url
`curl -X DELETE "http://localhost:8081/api?url=http%3A%2F%2Flocalhost%3A8081%2FpRA4OEy" -H "accept: application/json"`

Sample response
```
{
    "status": "ok",
    "operation": "delete",
    "url": "http://localhost:8081/pRA4OEy"
}
```

#### Count redirections
`curl -X GET "http://localhost:8081/api/count?url=http%3A%2F%2Flocalhost%3A8081%2FpRA4OEy" -H "accept: application/json"`

Sample response
```
{
    "status": "ok",
    "operation": "count",
    "count": 4
}
```

#### Redirect
Open url `http://localhost:8081/pRA4OEy` in your browser