# README #

This README would normally document whatever steps are necessary to get your application up and running..

### How to run http rest api
- Ensure you are in the appsku-golang directory.
- Create a .env file in the appsku-golang directory based on the example in the `.env.example` file.
- Generate the `wire_gen.go` file for dependency injection using the command `wire ./...`
- To run the REST HTTP API, use the command go run `./app/cmd/ main`

### How to run consumer
- Ensure you are in the appsku-golang directory.
- Create a .env file in the appsku-golang directory based on the example in the `.env.example` file.
- Generate the `wire_gen.go` file for dependency injection using the command `wire ./...`
- To run the REST HTTP API, use the command go run `./app/cmd/ --name=consumer_name consumer`

### How to use google wire
- Install google wire (you can check in https://github.com/google/wire)
- Add new controller, usecase, service and repository in app folder
- Add new set to wire.go file
- Run :
```
wire ./...
```#   c l e a n - g o l a n g - b o i l e r p l a t e  
 