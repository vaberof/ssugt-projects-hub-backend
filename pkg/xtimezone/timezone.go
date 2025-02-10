package xtimezone

import "time"

var nsk, _ = time.LoadLocation("Asia/Novosibirsk")
var msk, _ = time.LoadLocation("Europe/Moscow")

func NovosibirskLocation() *time.Location {
	return nsk
}

func MoscowLocation() *time.Location {
	return msk
}
