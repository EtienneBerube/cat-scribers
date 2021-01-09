package main

import (
	"github.com/EtienneBerube/only-cats/cmd"
	"github.com/EtienneBerube/only-cats/pkg/config"
)

func main() {
	conf := config.Init()
	cmd.RunServer(conf)
}