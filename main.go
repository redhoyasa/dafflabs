package main

import (
	"context"
	"fmt"
	"github.com/gojektech/heimdall/v6/hystrix"
	"github.com/labstack/echo/v4"
	qm "github.com/quickmetrics/qckm-go"
	telegram "github.com/redhoyasa/dafflabs/internal/client/telegram"
	"github.com/redhoyasa/dafflabs/internal/repository/tokopedia"
	"github.com/redhoyasa/dafflabs/internal/service/pricealert"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultConfigName = "application"
	defaultConfigType = "yaml"
	defaultPortKey    = "PORT"
	defaultPortValue  = "3000"
)

type PingResponse struct {
	Ping string `json:"ping"`
}

func main() {
	viper.AutomaticEnv()
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType(defaultConfigType)
	viper.SetDefault(defaultPortKey, defaultPortValue)

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn().Msg("failed to read config file, reading from Environment Variables instead" + err.Error())
	} else {
		log.Warn().Msg("reading config from file")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	qm.Init(qm.Options{
		ApiKey: viper.GetString("QUICKMETRICS_KEY"),
	})

	e := createHTTPServer()
	c := createScheduler()
	listenerPort := fmt.Sprintf(":%s", viper.GetString("PORT"))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.Start(listenerPort); err != nil {
			log.Info().Msg("HTTP server has been stopped...")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Start()
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-stopChan

	log.Info().Msg("Stopping HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg("Stopping scheduler...")
	ctx = c.Stop()
	<-ctx.Done()

	wg.Wait()
	log.Info().Msg("FIN")
}

func createHTTPServer() (e *echo.Echo) {
	e = echo.New()
	e.GET("/api/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &PingResponse{Ping: "pong"})
	})
	return
}

func createScheduler() (c *cron.Cron) {
	hystrix.NewClient()
	tokopediaClient, err := tokopedia.NewClient(hystrix.NewClient())
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	telegramClient, err := telegram.NewClient(viper.GetString("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	svc, err := pricealert.NewClient(*telegramClient, tokopediaClient)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	c = cron.New()
	c.AddFunc("30 9 * * *", func() {
		svc.CheckPrice(context.Background(), "https://www.tokopedia.com/matchamu/matchamu-matcha-latte-20pcs")
		svc.CheckPrice(context.Background(), "https://www.tokopedia.com/unicharm/tokocabang-mamypoko-popok-perekat-royal-soft-nb-52-2-packs")
	})
	return
}
