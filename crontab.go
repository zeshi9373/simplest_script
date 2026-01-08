package main

import (
	"simplest_script/crontab"
)

func InitCrontab() {
	go crontab.Init()
}
