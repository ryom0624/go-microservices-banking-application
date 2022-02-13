Go Microservices Banking Application

https://www.udemy.com/course/rest-based-microservices-api-development-in-go-lang/

# Feature
- Open New Account
- Deposit or Withdrawal transaction
- Secure by Role Based Access Control(RBAC)

# Architecture
- Hexagonal Architecture
  - Dependency inversion in Go 
  
# Setup
## Build and Run MySQL by docker
```shell
$ docker-compose -f ./resources/docker/docker-compose.yaml up -d
```


## Set Environments
 - SERVER_ADDRESS [IP Address of the machine]
 - SERVER_PORT [Port of the machine]
 - DB_USER [Database username]
 - DB_PASSWORD [Database password]
 - DB_ADDRESS [IP address of the database]
 - DB_PORT [Port of the database]
 - DB_NAME [Name of the database]
 
