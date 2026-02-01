package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	`github.com/agnostic-play/ditoo/cmd/http`
	"github.com/agnostic-play/ditoo/internal/config"
)

const banner = `
BonClay 
version %s | OS %s %s %s CPU %v
`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf(banner+"\n", "1.0", runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.NumCPU())
	conf := config.LoadConfigFile("./resources/config/config.yaml")

	appExitSignal, serverExitSignal := http.Runner(conf)

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
