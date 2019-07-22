package main

import (
	"errors"
	"github.com/globalsign/mgo"
	"time"
)

type Module struct {
	db *mgo.Database
}

func (m *Module) Close() {
	//todo
	//m.db.Session.Close()
}

func NewModule() *Module {
	m := &Module{}
	return m
}

func (m *Module) InitDataBase(host string, dbName string, userName string, pass string) {
	dialInfo := mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  time.Second * 10,
		Database: dbName,
		Username: userName,
		Password: pass,
		Source:   "admin",
	}
	session, err := mgo.DialWithInfo(&dialInfo)
	panicIfError(err, "mgo dial error")
	// set mode
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbName)
	if db == nil {
		panicIfError(errors.New("db is nil"), "db is nil")
	}
	m.db = db
}

func (m *Module) GetCollection(collectionName string) *mgo.Collection {
	return m.db.C(collectionName)
}
