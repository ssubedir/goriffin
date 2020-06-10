<p align="center"><img src="https://i.imgur.com/s74yXT7.png" width="400"></p>
<p align="center">
<a href="https://github.com/ssubedir/goriffin/actions"><img src="https://github.com/ssubedir/goriffin/workflows/Go/badge.svg" alt="Build Test"></a>
<a href="https://goreportcard.com/report/github.com/ssubedir/goriffin"><img src="https://goreportcard.com/badge/github.com/ssubedir/goriffin" alt="Latest Stable Version"></a>
<a href="https://github.com/ssubedir/goriffin/blob/master/LICENSE"><img src="https://poser.pugx.org/laravel/framework/license.svg" alt="License"></a>
</p>

# About Goriffin
Goriffin is written in go, a heartbeat surveillance program. Goriffin is packed with goriffin API and the goriffin background processor. The API is used to configure the services being monitored and the background processor checks the services periodically and reports the status of that same services to a Mongo db

## Getting Started

These instructions will get you a copy of the project up and run the service on your local machine.See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You will need to install go version go1.11 +

### Download

Use git clone to get your local copy 
```
git clone https://github.com/ssubedir/goriffin
```

## Deployment

Use the go build command to build each components 

### Build Project

Build Goriffin API
```
go build -v cmd/goriffin/goriffin.go 
```
Build Goriffin background processor
```
go build -v cmd/background/goriffin_background.go 
```
### Deploy

Before you run Goriffin create .env file with the following Environmental Variables

```
# MongoDB info
DB0_USER=username
DB0_PASSWORD=password
DB0_CLUSTER=address/cluster

  
#Database
DB_DBNAME= database name
SERVICES_HEARTBEAT=collection name
HEARTBEAT_HTTP=collection name

 
# Number of workers for the background processing 
WORKER_COUNT=16

#Goriffin Api port
API_PORT=9000
#Goriffin background processor grpc port
GRPC_PORT=9092
```
## Goriffin Api

### Sattus of Background Processor 
```
GET localhost:API_PORT/status
```
Returns the status of the background goriffin processor by making gRPC request

### Fetch all services
```
GET localhost:API_PORT/services
```
Returns all the services that are being monitored.

```
// Response
{
	"status": true,
  "time": "2020-05-04 23:39:50.8786948 -0400 EDT"
}
```

### Fetch a service
```
GET localhost:API_PORT/service
```
```
// Example request template
{
	 "host":"example.com",
	 "stype":"http"
}
```

Returns a single service.

```
// Response
{
	"host":"example.com",
	"name":"example",
	"stype":"http",
	"accepted_status":{
		"200":true
		},
	"frequency":30,
	"request_method":"GET",
	"request_payload":"",
	"request_headers":[["Key","Value"]]
}
```
### Add a service
```
POST localhost:API_PORT/service
```
```
// Example request template
{
	"host":"example.com",
	"name":"example",
	"stype":"http",
	"accepted_status":{
		"200":true
		},
	"frequency":30,
	"request_method":"GET",
	"request_payload":"",
	"request_headers":[["Key","Value"]]
}
```


Returns the the added service
```
// Response
{
	"host":"example.com",
	"name":"example",
	"stype":"http",
	"accepted_status":{
		"200":true
		},
	"frequency":30,
	"request_method":"GET",
	"request_payload":"",
	"request_headers":[["Key","Value"]]
}

```

### Remove a service
```
DELETE localhost:API_PORT/service
```
```
// Example request template
{
	 "host":"example.com",
	 "stype":"http"
}
```

Returns the status of the action

```
// Response
{
	"host": "example.com"
	"status": true
}

```




## Built With

* [GO](https://golang.org/) - Programming language


## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/ssubedir/goriffin/blob/master/LICENSE) file for details

