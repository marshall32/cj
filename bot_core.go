package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cardigann/go-cloudflare-scraper"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// App stores program state
type App struct {
	config         Config
	mongo          *mgo.Session
	accounts       *mgo.Collection
	chat           *mgo.Collection
	discordClient  *discordgo.Session
	httpClient     *http.Client
	ready          chan bool
	cache          *cache.Cache
	commandManager *CommandManager
}

// Start starts the app with the specified config and blocks until fatal error
func Start(config Config) {
	scrpr, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		config:     config,
		httpClient: &http.Client{Transport: scrpr},
		cache:      cache.New(5*time.Minute, 30*time.Second),
	}

	logger.Debug("started with debug logging enabled",
		zap.Any("config", app.config))

	app.ConnectDB()
	app.StartCommandManager()
	app.ConnectDiscord()

	app.newPostAlert("3", func() {
		app.discordClient.ChannelMessageSend(app.config.PrimaryChannel, "New Kalcor Post: http://forum.sa-mp.com/search.php?do=finduser&u=3")
	})

	done := make(chan bool)
	<-done
}
