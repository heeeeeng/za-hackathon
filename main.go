package main

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.toml", "config config.toml")
}

func main() {
	mergeLocalConfig(configPath)

	viper.Debug()

	m := NewModule()
	c := NewController(m)
	log.SetLevel(log.DebugLevel)
	host := viper.GetString("db.host")
	username := viper.GetString("db.user_name")
	pass := viper.GetString("db.pass")
	dbName := viper.GetString("db.name")
	webUri := viper.GetString("web.uri")
	m.InitDataBase(host, dbName, username, pass)
	ogurl := viper.GetString("og.url")
	contractAddress := viper.GetString("og.contract_address")
	if ogurl == "" || contractAddress == "" {
		panicIfError(errors.New("miss og url or contract address"), "")
	}
	s := NewRankSpider(m, ogurl, contractAddress, webUri)
	c.rankSpider = s
	s.Start()
	c.InitRouter()
	fmt.Println("---------Server Start!---------")
	fmt.Println("Port: ", 10001)
	go func() {
		log.Fatal(http.ListenAndServe(":10001", nil))
	}()

	// prevent sudden stop. Do your clean up here
	var gracefulStop = make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	func() {
		sig := <-gracefulStop
		log.Warnf("caught sig: %+v", sig)
		log.Warn("Exiting... Please do no kill me")
		s.Stop()
		m.Close()
		os.Exit(0)
	}()
}

func mergeLocalConfig(configPath string) {
	absPath, err := filepath.Abs(configPath)
	panicIfError(err, fmt.Sprintf("Error on parsing config file path: %s", absPath))

	file, err := os.Open(absPath)
	panicIfError(err, fmt.Sprintf("Error on opening config file: %s", absPath))
	defer file.Close()

	viper.SetConfigType("toml")
	err = viper.MergeConfig(file)
	panicIfError(err, fmt.Sprintf("Error on reading config file: %s", absPath))

	viper.SetEnvPrefix("hk")
	viper.AutomaticEnv()
	return
}

func panicIfError(err error, message string) {
	if err != nil {
		fmt.Println("will panic ", err, message)
		panic(err)
	}

}
