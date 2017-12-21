package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gost/sensorthings-connector/configuration"
	"github.com/gost/sensorthings-connector/connector"
	log "github.com/sirupsen/logrus"
	"github.com/tebben/discordrus"
)

var (
	config       = configuration.Config{}
	reportTicker *time.Ticker
)

func main() {
	initShutdownListener()
	initLogRus()
	initConfig()
	addLoggingDiscordHook(config.Logging.Discord)
	startStatusReporter(config.Logging.Status)
	connector.Start(config.Connector)
}

// initLogRus configures logrus for use in sensorthings-connector
func initLogRus() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

// initConfig reads the default or supplied config file.
// sensorthings-connector service stops when unable to read the config or config contains errors
func initConfig() {
	cfgFlag := flag.String("config", "config.json", fmt.Sprintf("path to the %s config file", connector.NAME))
	flag.Parse()
	var err error
	config, err = configuration.GetConfig(*cfgFlag)
	if err != nil {
		log.Fatal("config read error: ", err)
	}
}

// initShutdownListener listen for shutdown to cleanup
func initShutdownListener() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-stop
		if reportTicker != nil {
			reportTicker.Stop()
		}

		connector.Stop()
		log.Infof("%s stopped", connector.NAME)
		os.Exit(1)
	}()
}

func addLoggingDiscordHook(cfg configuration.DiscordConfig) {
	if !cfg.Enabled || len(cfg.Hook) == 0 {
		return
	}

	log.AddHook(discordrus.NewHook(
		cfg.Hook,
		log.InfoLevel,
		&discordrus.Opts{
			Username:           cfg.Name,
			Author:             "",
			DisableTimestamp:   false,
			TimestampFormat:    "Jan 2 15:04:05.00000",
			EnableCustomColors: true,
			CustomLevelColors: &discordrus.LevelColors{
				Debug: 10170623,
				Info:  3581519,
				Warn:  14327864,
				Error: 13631488,
				Panic: 13631488,
				Fatal: 13631488,
			},
			DisableInlineFields: false,
		},
	))
}

func startStatusReporter(cfg configuration.StatusConfig) {
	if !cfg.Enabled {
		return
	}

	interval := cfg.IntervalSeconds
	if interval == 0 {
		interval = 3600 // default to 1 hour if not set
	}

	reportTicker = time.NewTicker(time.Second * time.Duration(interval))
	go func() {
		for range reportTicker.C {
			for _, m := range connector.Modules {
				data := (*m).GetConnectorModuleData()
				status := data.Status
				log.WithFields(log.Fields{
					"Running":          status.Running,
					"Latest GET time":  status.LastGet,
					"Latest POST time": status.LastPost,
					"POST success":     status.ObservationsPostedOk,
					"POST failed":      status.ObservationsPostedFailed,
					"Errors":           status.ErrorCount,
				}).Infof("Status report for module %s", data.ModuleFileName)
			}
		}
	}()
}
