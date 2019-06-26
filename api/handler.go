package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack/slackevents"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (c *bot) Send(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := c.SendPollMessage([]string{userId}, c.Question)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func (c *bot) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf))
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr = strings.TrimPrefix(jsonStr, "payload=")

	var msg slackevents.MessageAction
	err = json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		log.Printf("[ERROR] Failed to parse request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if msg.CallbackID != c.CallbackID || len(msg.Actions) == 0 {
		return
	}

	if msg.Actions[0].Name != actionName {
		return
	}

	log.Println(msg.User.ID, msg.User.Name, msg.Actions[0].Value)
	// TODO processing answer

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Thank you!"))
}
