package main 

import (
	"os"
	"log"
	"runtime"
	"github.com/h2object/h2object/commands"
)

func main() {
	defer func(){
		if r := recover(); r != nil {
            log.Println(os.Args[0], " ", r)
        }
	}()
	
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.App().Run(os.Args)
}