package main

import "fmt"

type API struct {
	Client       *Client
	WebHookEvent WebhookEvent
}

func NewAPI(githubToken string, webhook WebhookEvent) (*API, error) {
	client, err := NewClient(githubToken, nil)
	if err != nil {
		return nil, err
	}

	return &API{
		Client:       client,
		WebHookEvent: webhook,
	}, nil
}

/*
POST /repos/:owner/:repo/issues/:number/labels
[
    "Label1",
    "Label2"
]
*/
func (self *API) addLabels(labels []string) error {
	reqUri := "/repos/" + self.WebHookEvent.RepoFullName + "/issues/" + fmt.Sprintf("%.0f", self.WebHookEvent.Number) + "/labels"
	req, err := self.Client.NewRequest("POST", reqUri, labels)
	if err != nil {
		return err
	}
	// fmt.Println(req)

	r, err := self.Client.Do(req, nil)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("Code: %d Limit: %s %s", r.StatusCode, r.Header.Get("X-Ratelimit-Limit"), r.Header.Get("X-RateLimit-Remaining"))
	}
	return nil
}

/*
POST /repos/:owner/:repo/issues/:number/assignees

{
  "assignees": [
    "hubot",
    "other_assignee"
  ]
}
*/
func (self *API) addAssignees(assignees []string) error {

	payload := struct {
		Assignees []string `json:"assignees"`
	}{Assignees: assignees}

	reqUri := "/repos/" + self.WebHookEvent.RepoFullName + "/issues/" + fmt.Sprintf("%.0f", self.WebHookEvent.Number) + "/assignees"
	req, err := self.Client.NewRequest("POST", reqUri, payload)
	if err != nil {
		return err
	}

	// var data interface{}

	r, err := self.Client.Do(req, nil)
	// r, err := self.Client.Do(req, &data)
	if err != nil {
		return err
	}
	if r.StatusCode != 201 {
		return fmt.Errorf("Code: %d Limit: %s %s", r.StatusCode, r.Header.Get("X-Ratelimit-Limit"), r.Header.Get("X-RateLimit-Remaining"))
	}

	// fmt.Println(data)

	return nil
}
