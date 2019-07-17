package main

import (

	"errors"
	"flag"
	"fmt"
	"github.com/annchain/OG/common/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var configPath string

func init()  {
	flag.StringVar(&configPath,"config","config.toml","config config.toml")
}

func main() {
	mergeLocalConfig(configPath)
	m := NewModule(viper.GetString("leveldb.path"))
	c:= NewController(m)
	c.InitRouter()
	ogurl := viper.GetString("og.url")
	if ogurl =="" {
		panicIfError(errors.New("miss og url"),"")
	}
	_, sk:= crypto.Signer.RandomKeyPair()
	s := NewRankSpider(m,ogurl,sk)
	s.Start()
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
		s.stop()
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
	return
}

func panicIfError(err error , message string)  {
	if err!=nil {
		fmt.Println("will panic ",err,message)
		panic(err)
	}

}

