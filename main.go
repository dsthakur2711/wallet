package main

import (
	"github.com/dsthakur2711/wallet/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	logs "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func dbConn() (db *gorm.DB) {
	/*dbDriver := "mysql"
	dbUser := "root"
	dbPass := "password"
	dbName := "walletDB"*/
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/walletDB")
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/walletDB?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {

	// SetFormatter sets the standard logger formatter.
	logs.SetFormatter(&logs.TextFormatter{})
	logs.Println("starting wallet service")

	//db := dbConn()
	//defer db.Close()

	server.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
}