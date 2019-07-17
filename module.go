package main

import "github.com/annchain/OG/ogdb"

type Module struct {
	db ogdb.Database
}

func (m*Module)Close() {
	m.db.Close()
}

func NewModule(dbPath string ) *Module {
	m:= &Module{

	}
	db,err :=  ogdb.NewLevelDB(dbPath,16,16)
	if err!=nil {
		panicIfError(err,"")
	}
	m.db = db
	return m
}