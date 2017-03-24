package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"encoding/json"
)

type WebhookEvent struct {
	Action       string        // action
	Sender       string        // sender.login
	Number       float64       // pull_request.number
	Assignees    []interface{} // pull_request.assignees
	RepoFullName string        // repository.full_name
}

func parse(body interface{}) (WebhookEvent, error) {
	var webhookEvent WebhookEvent

	// action
	if _, ok := body.(map[string]interface{})["action"]; !ok {
		return webhookEvent, fmt.Errorf("The action attribute is not exist")
	}
	webhookEvent.Action = body.(map[string]interface{})["action"].(string)

	// sender
	if _, ok := body.(map[string]interface{})["sender"].(map[string]interface{})["login"]; !ok {
		return webhookEvent, fmt.Errorf("The sender.login attribute is not exist")
	}
	webhookEvent.Sender = body.(map[string]interface{})["sender"].(map[string]interface{})["login"].(string)

	// number
	if _, ok := body.(map[string]interface{})["pull_request"].(map[string]interface{})["number"]; !ok {
		return webhookEvent, fmt.Errorf("The pull_request.number attribute is not exist")
	}
	webhookEvent.Number = body.(map[string]interface{})["pull_request"].(map[string]interface{})["number"].(float64) // float64

	// assignees
	if _, ok := body.(map[string]interface{})["pull_request"].(map[string]interface{})["assignees"]; !ok {
		return webhookEvent, fmt.Errorf("The pull_request.assignees attribute is not exist")
	}
	webhookEvent.Assignees = body.(map[string]interface{})["pull_request"].(map[string]interface{})["assignees"].([]interface{})

	// repo full name
	if _, ok := body.(map[string]interface{})["repository"].(map[string]interface{})["full_name"]; !ok {
		return webhookEvent, fmt.Errorf("The repository.full_name attribute is not exist")
	}
	webhookEvent.RepoFullName = body.(map[string]interface{})["repository"].(map[string]interface{})["full_name"].(string)

	return webhookEvent, nil
}

func isPullRequestEvent(xGitHubEvent string) bool {
	if xGitHubEvent == "pull_request" {
		return true
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {

	if !isPullRequestEvent(r.Header.Get("X-GitHub-Event")) {
		w.WriteHeader(200)
		return
	}

	var body interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	webhookEvent, err := parse(body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	if _, ok := config.RepoConfigs[webhookEvent.RepoFullName]; !ok {
		w.WriteHeader(200)
		return
	}

	if !(webhookEvent.Action == "opened" || webhookEvent.Action == "reopened") {
		w.WriteHeader(200)
		return
	}

	token := config.RepoConfigs[webhookEvent.RepoFullName].Token
	api, err := NewAPI(token, webhookEvent)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	if len(webhookEvent.Assignees) == 0 {
		members := config.RepoConfigs[webhookEvent.RepoFullName].Members
		assignFunc := func(s []string, name string) []string {
			var ret []string
			for _, v := range s {
				if v != name {
					ret = append(ret, v)
				}
			}
			return ret
		}

		assignees := assignFunc(members, webhookEvent.Sender)
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(assignees))

		if err := api.addAssignees([]string{assignees[i]}); err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
	}

	labels := config.RepoConfigs[webhookEvent.RepoFullName].Labels
	if err := api.addLabels(labels); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("healthy"))
}

func runServer(listen string) {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthCheckHandler)
	http.ListenAndServe(listen, nil)
}
