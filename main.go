package main

import (
	"github.com/EtienneBerube/cat-scribers/cmd"
	"github.com/EtienneBerube/cat-scribers/pkg/config"
)

func main() {
	config.Init()
	cmd.RunServer()
}