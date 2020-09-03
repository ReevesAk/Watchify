package main

import (
	"Watchify/helper"

	bolt "go.etcd.io/bbolt"
)

type Database struct {
	database bolt.DB
}

func main() {
	fyneDb := Database{}

	helper.Run()
	fyneDb.database.Close()
}
