# goco
 
# example

```

package main

import (
	"log"

	"github.com/lonelymous/goco"
)

type ServerConfig struct {
	Host string
	Port int
}

func main() {
	serverConfig := &ServerConfig{}

	// Setup config
	err := goco.InitializeConfig(serverConfig)
	if err != nil {
		log.Fatalln("error while setup config", err)
	}

	log.Println(serverConfig)

}

```