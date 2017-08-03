package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kensodev/micro-github-auth-proxy"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenPort     = kingpin.Flag("listen-port", "Which port should the proxy listen on").Required().Int()
	configLocation = kingpin.Flag("config-location", "Proxy Config Location").Required().String()
)

func main() {
	kingpin.Parse()

	configReader := authproxy.NewConfigurationReader(*configLocation)

	data, err := configReader.ReadConfigurationFile()
	handleError(err)

	config, err := authproxy.NewConfiguration(data)
	handleError(err)

	authContext := authproxy.NewGithubAuthContext(config)
	http.Handle("/callback", authContext)

	authproxy.NewHttpListeners(config)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil); err != nil {
		fmt.Errorf("Could not listen on port %s: %s", *listenPort, err.Error())
		os.Exit(1)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error reading your configuration file: %s", err.Error())
		os.Exit(1)
	}
}
