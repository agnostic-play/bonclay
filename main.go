package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/agnostic-play/ditoo/cmd"
	"github.com/agnostic-play/ditoo/internal/config"
)

const banner = `
██████╗ ██╗████████╗ ██████╗  ██████╗ 
██╔══██╗██║╚══██╔══╝██╔═══██╗██╔═══██╗
██║  ██║██║   ██║   ██║   ██║██║   ██║
██║  ██║██║   ██║   ██║   ██║██║   ██║
██████╔╝██║   ██║   ╚██████╔╝╚██████╔╝
╚═════╝ ╚═╝   ╚═╝    ╚═════╝  ╚═════╝ 
DITOO service for Mocking Service
version %s | OS %s %s %s CPU %v
`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf(banner+"\n", "1.0", runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.NumCPU())
	conf := config.LoadConfigFile("./config.yaml")

	appExitSignal, serverExitSignal := cmd.Init(conf)

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)

	for range interruptSignal {
		appExitSignal <- true
		serverExitSignal <- true
		<-appExitSignal    // Finish
		<-serverExitSignal // Finish
		return             // Now we can safely exit the app
	}
}
