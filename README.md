# crud-api

## Installation

`git clone https://github.com/v1tbrah/crud-api`

## Example usage

* Go to project working directory. For example: `cd ~/go/src/crud-api`
* Run app. For example: `go run cmd/main.go`

## Options
The following options are set by default:
```
api server run address: `:3333`
storage type: `file`
   (currently only `file` is supported)
path to file storage: `users.json`
```
* flag options:
```
   -a string
      api server run address
   -st string
      type of data storage 
   -pf string
      path to file data storage
```
For example: `go run cmd/main.go -a=:8081 -st=file -pf=users.json`
* env options:
```
  RUN_ADDRESS string
      api server run address
  STORAGE_TYPE string
      type of data storage
  -PATH_FILE_STORAGE string
      path to file data storage
      
```