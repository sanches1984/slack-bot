package main

import (
	"github.com/gorilla/mux"
	"github.com/sanches1984/slackbot/api"
	"log"
	"net/http"
)

const token = "your_token"

func main() {
	bot := api.NewSlackBot(token, "answer_poll", "How was your day?")

	router := mux.NewRouter()
	router.HandleFunc("/interaction", bot.Callback)
	router.HandleFunc("/send/{user_id}", bot.Send)

	log.Printf("[INFO] Server listening")
	if err := http.ListenAndServe("localhost:8082", router); err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}
}
