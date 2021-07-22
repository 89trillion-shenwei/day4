package main

import (
	"day4/app/http"
	"day4/internal/util"
)

func main() {
	//初始化redis池和mongo
	util.Init()
	http.Start()

}
