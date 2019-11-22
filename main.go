package todofun

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/y-yagi/goext/arr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Title string `json:"title"`
		URL   string `json:"url"`
		ID    string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		// TODO: error handling
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if d.Title == "" {
		log.Println("Title is empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	supportedIDs := strings.Split(os.Getenv("SUPPORTED_IDS"), ",")
	if !arr.Contains(supportedIDs, d.ID) {
		log.Println("ID is not supported:", d.ID)
		w.WriteHeader(http.StatusNotFound)
	}

	if err := insertTask(os.Getenv("TASK_LIST_ID"), d.Title, d.URL); err != nil {
		// TODO: error handling
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func buildTaskService() (*tasks.Service, error) {
	ctx := context.Background()

	config, err := google.ConfigFromJSON([]byte(os.Getenv("OAUTH_CREDENTIALS")), tasks.TasksScope)
	if err != nil {
		return nil, err
	}

	client, err := getClient(config)
	ts, err := tasks.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func getClient(config *oauth2.Config) (*http.Client, error) {
	tok := &oauth2.Token{}
	if err := json.NewDecoder(strings.NewReader(os.Getenv("OAUTH_TOKEN"))).Decode(tok); err != nil {
		return nil, err
	}

	return config.Client(context.Background(), tok), nil
}

func insertTask(tasklistid string, title string, notes string) error {
	var task tasks.Task
	task.Title = title
	task.Notes = notes

	ts, err := buildTaskService()
	if err != nil {
		return err
	}

	_, err = ts.Tasks.Insert(tasklistid, &task).Do()
	return err
}
