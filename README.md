#Metrics-Collector

Metrics-Collector is one of the components of CF `app-autoscaler`. It is used to collect application metrics from CF loggretator. The current version only support memory metrics, it will be extended to include other metrics like throughput and response time. 

## Getting started

###System requirements:

* Go 1.5 or above
* Cloud Foundry relese 235 or later

### build and test

1. Create a directory where you would like to store the source for Go projects and their binaries 
2. Add this directory to your `$GOPAHT`
3. create the directory `src/github.com/cloudfoundry-incubator` 
3. change directory to `src/github.com/cloudfoundry-incubator`
4. clone the `app-autoscaler` project: `clone git@github.com:cloudfoundry-incubator/app-autoscaler.git`
5. change directory to 'src/github.com/cloud-foundry-incubator/app-autoscaler/metrics-collector'
6. build the project: `go build -o out/mc`
7. test the project: `ginkgo -r`

### run the metrics-collector

Firstly a configuration file needs to be created. Examples can be found under `example-config` directory. Here is an example:

```
cf:
  api: "https://api.bosh-lite.com"
  grant_type: "password"
  user: "admin"
  pass: "admin"
server:
  port: 8080
logging:
  level: "info"
  file: "logs/mc.log"
  log_to_stdout: true
```


The config parameters are explained as below

* `cf` : cloudfoundry config
 * `api`: API endpoint of cloudfoundry
 * `grant_type`: the grant type when you login into coud foundry, can be "password" or "client_credentials"
 * `user`: the user name when using password grant to login cloudfoundry
 * `pass`: the password when using password grant to login cloudfoundry
 * `client_id`: the client id when using client_credentials grant to login cloudfoundry
 * `secret`: the client secret when using client_credentials grant to login cloudfoundry
* server: API sever config
 * `port`: 8080 - the port API sever will use
* `logging`
 * level: the level of logging, can be 'debug', 'info', 'error' and 'fatal'
 * file:  the log file name, "" means not log to file
 * log_to_stdout: whehter show logs to stdout


To run the metrics-collector, use `./out/mc -c config_file_name'

## API

Metrics Collector exposes the following APIs for other CF App-Autoscaler components to retrieve metrics.

| PATH                      | METHOD  | Description                              |
|---------------------------|---------|------------------------------------------|
| /v1/apps/{appid}/metrics/memory | GET | Get the latest memroy metric for an application |







