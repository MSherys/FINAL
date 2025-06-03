package server

import (
	"fmt"
	nextdate "go1f/pkg/api"
	"net/http"
)

func Run() error {
	//nextdate.Init()
	newr := nextdate.Init()
	port := 7540
	fmt.Println("Запуск сервера на порту %d\n", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), newr)
	//http.Handle("/", http.FileServer(http.Dir("web")))
	//fmt.Println("WWWWWW")
	//return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
