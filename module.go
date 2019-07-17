package main

type Module struct {
	Db *LevelDB
}

func NewModule(dbPath string ) *Module {
	m:= &Module{

	}
	db,err :=  NewLevelDB(dbPath,16,16)
	if err!=nil {
		panicIfError(err,"")
	}
	m.Db = db
	return m
}