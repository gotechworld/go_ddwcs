package main

import (
	"github.com/joho/godotenv"
	_ "gitlab.altex.ro/ams/go_ddwcs/awb/controller"
	_ "gitlab.altex.ro/ams/go_ddwcs/delivery/controller"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	"os"
	"path"
	"runtime"
)

/**
 * Important thing - SET the main dir which will be used throughout the app
 */
func init() {
	_, b, _, _ := runtime.Caller(0)
	err := os.Setenv("APP_MAIN_DIR", path.Dir(b))
	if err != nil {
		println("Cannot set APP_MAIN_DIR")
		panic("Cannot start app !!!")
	}
}

/**
 * Entry point
 */
func main() {
	godotenv.Load(os.Getenv("APP_MAIN_DIR") + "/.env")
	lib.Run()
}
