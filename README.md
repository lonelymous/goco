# goco
 
# example

``` go
package main

import (
	"log"

	"github.com/lonelymous/goco"
)

type testNestedConfig struct {
	Name   string `docker:"NAME"`
	Server string `docker:"SERVER"`
	Port   int    `docker:"PORT"`
}

type testConfig struct {
	IsDevelopment bool             `docker:"IS_DEVELOPMENT"`
	PortNumber    int              `docker:"PORT_NUMBER"`
	Nested        testNestedConfig `docker:"NESTED"`
}

func main() {
	// Load config from file 
	godotenv.Load(".env")

	// Setup config
	serverConfig := &ServerConfig{}
	err := InitializeConfig(&serverConfig)
	if err != nil {
		fmt.Println("error while setup config", err)
		return
	}

	fmt.Println(serverConfig)
}
```

# config.ini file

``` ini
IsDevelopment = True
PortNumber = 3000

# Config for the nested server
[Nested]
Name = server_name
Server = http://localhost
Port = 8080
```

# .env file

``` env
DOCKER = true
IS_DEVELOPMENT = true
PORT_NUMBER = 3001

NESTED_NAME = server_name
NESTED_SERVER = http://localhost
NESTED_PORT = 8081
```