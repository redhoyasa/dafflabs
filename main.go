package main

import (
	"context"
	"fmt"
	"github.com/gojektech/heimdall/v6/hystrix"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	qm "github.com/quickmetrics/qckm-go"
	telegram "github.com/redhoyasa/dafflabs/internal/client/telegram"
	"github.com/redhoyasa/dafflabs/internal/database"
	"github.com/redhoyasa/dafflabs/internal/http/handler"
	"github.com/redhoyasa/dafflabs/internal/migration"
	"github.com/redhoyasa/dafflabs/internal/repository"
	"github.com/redhoyasa/dafflabs/internal/repository/tokopedia"
	"github.com/redhoyasa/dafflabs/internal/service/pricealert"
	"github.com/redhoyasa/dafflabs/internal/service/wishlist"
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

	err = migration.ExecuteMigration(viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	db, err := database.NewConn(viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	repo := repository.NewWishRepo(db)

	e := createHTTPServer(repo)
	c := createScheduler(repo)
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

func createHTTPServer(wishRepo wishlist.WishlistRepoIFace) (e *echo.Echo) {

	tokopediaClient, _ := tokopedia.NewClient(hystrix.NewClient())
	svc := wishlist.NewWishlistSvc(wishRepo, tokopediaClient)
	h := handler.NewHandler(svc)

	e = echo.New()
	api := e.Group("/api")
	api.GET("/api/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &PingResponse{Ping: "pong"})
	})

	wishlist := api.Group("/wishlist")
	wishlist.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == viper.GetString("API_KEY"), nil
	}))

	wishlist.POST("/wish", func(c echo.Context) error {
		return h.AddWish(c)
	})
	wishlist.DELETE("/wish/:id", func(c echo.Context) error {
		return h.DeleteWish(c)
	})
	wishlist.GET("/wish/customer/:customer_ref_id", func(c echo.Context) error {
		return h.FetchCustomerWishlist(c)
	})
	return
}

func createScheduler(wishRepo wishlist.WishlistRepoIFace) (c *cron.Cron) {
	hystrix.NewClient()
	tokopediaClient, err := tokopedia.NewClient(hystrix.NewClient())
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	telegramClient, err := telegram.NewClient(viper.GetString("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	svc, err := pricealert.NewClient(*telegramClient, tokopediaClient, wishRepo)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	c = cron.New()
	c.AddFunc("* * * * *", func() {
		svc.GeneratePriceChecker()
	})
	return
}
