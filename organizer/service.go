package organizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type OrganizerService struct {
	client http.Client
	config OrganizerServiceConfig
	url    string
}

func InitializeOrgService() OrganizerService {
	config := LoadConfig()

	var url string

	url = config.Protocol + "://" + config.Domain + ":" + config.Port

	service := OrganizerService{}
	service.url = url
	service.config = config

	service.client = http.Client{
		Timeout: time.Minute,
	}

	return service
}

func (os OrganizerService) Log(l ...interface{}) {
	if false {
		log.Println("[SERVICE]", l, "\n")
	}
}

func (os OrganizerService) get(path string, query map[string]string, respData interface{}) {
	url := os.url + "/" + path + "?taskClient=1"

	for key, value := range query {
		url += "&" + key + "=" + value
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+os.config.Token)

	resp, err := os.client.Do(req)

	if err != nil {
		log.Fatal("Error sending request")
	}

	defer resp.Body.Close()

	os.Log("response Status:", resp.Status)
	os.Log("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//os.Log("response Body:", string(body))

	json.Unmarshal(body, &respData)
}

func (os OrganizerService) post(path string, data interface{}, respData interface{}) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal("Cannot encode data to json")
	}

	req, _ := http.NewRequest("POST", os.url+"/"+path, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", "Bearer "+os.config.Token)
	req.Header.Add("Content-type", "application/json; charset=UTF-8")

	resp, err := os.client.Do(req)

	if err != nil {
		log.Fatal("Error sending request")
	}

	defer resp.Body.Close()

	os.Log("response Status:", resp.Status)
	os.Log("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//os.Log("response Body:", string(body))

	json.Unmarshal(body, &respData)
}

func (os OrganizerService) CreateTask(label string) Task {
	task := Task{}
	task.Label = label

	var respTask Task

	os.post("task", task, respTask)

	os.Log("Created new task at:", respTask.CreateAt)

	return respTask
}

func (os OrganizerService) LoadTasks(limit int, offset int) []Task {
	var respTasks []Task

	query := make(map[string]string)
	query["limit"] = fmt.Sprint(limit)
	query["offset"] = fmt.Sprint(offset)

	os.get("task", query, &respTasks)

	return respTasks
}

func (os OrganizerService) Undone(t *Task) Task {
	var respTask Task

	os.post("task/"+t.Uuid+"/undone", nil, respTask)

	return respTask
}

func (os OrganizerService) Done(t *Task) Task {
	var respTask Task

	os.post("task/"+t.Uuid+"/done", nil, respTask)

	return respTask
}
