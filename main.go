package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/astronomer/astro-cli/cmd"
	"github.com/astronomer/astro-cli/config"
	"github.com/astronomer/astro-cli/houston"
	"github.com/astronomer/astro-cli/pkg/httputil"

	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	config.InitConfig(fs)
	httpClient := httputil.NewHTTPClient()
	// configure http transport
	dialTimeout := config.CFG.HoustonDialTimeout.GetInt()
	// #nosec
	httpClient.HTTPClient.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(dialTimeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(dialTimeout) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.CFG.SkipVerifyTLS.GetBool()},
	}
	houstonClient := houston.NewClient(httpClient)
	// setup log level before we start command since we will miss the feature flag checks other wise
	if err := cmd.SetUpLogs(os.Stdout, config.CFG.Verbosity.GetString()); err != nil {
		os.Exit(1)
	}
	if err := cmd.NewRootCmd(houstonClient, os.Stdout).Execute(); err != nil {
		os.Exit(1)
	}
}
