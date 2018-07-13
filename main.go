package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("environment variable missing: " + name)
	}
	return v
}

func main() {

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "keyword, k",
			Usage: "keywords to watch and retweet if matched",
		},
	}
	app.Name = "twitterbot"
	app.Usage = "retweets tweets based on keywords. Specify using -keyword flags"

	app.Action = func(c *cli.Context) error {

		anaconda.SetConsumerKey(consumerKey)
		anaconda.SetConsumerSecret(consumerSecret)
		api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
		log := logger{
			logrus.New(),
		}
		api.SetLogger(log)

		if len(c.StringSlice("keyword")) == 0 {
			log.Fatalln("At least one keyword must be specified using the -k flag")
		}

		stream := api.PublicStreamFilter(url.Values{
			"track": c.StringSlice("keyword"),
		})
		defer stream.Stop()
		defer logrus.Info("shutting down")

		for v := range stream.C {
			t, ok := v.(anaconda.Tweet)
			if !ok {
				log.Printf("recieved unexpected value of type %T", v)
				continue
			}

			if t.RetweetedStatus != nil {
				logrus.Debugln("already retweeted. ignoring")
				continue
			}
			logrus.Infoln("Retweeting:\t", t.Id)
			api.Retweet(t.Id, true)
		}
		return nil
	}

	//Start app
	_ = app.Run(os.Args)
}

type logger struct {
	*logrus.Logger
}

func (l logger) Critical(args ...interface{})            { l.Error(args...) }
func (l logger) Criticalf(f string, args ...interface{}) { l.Errorf(f, args...) }

func (l logger) Notice(args ...interface{})            { l.Info(args...) }
func (l logger) Noticef(f string, args ...interface{}) { l.Infof(f, args...) }
