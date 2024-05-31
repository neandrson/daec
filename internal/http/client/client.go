package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Vojan-Najov/daec/internal/result"
	"github.com/Vojan-Najov/daec/internal/task"
)

type Client struct {
	http.Client
}

func (client *Client) GetTask() *task.Task {
	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8081/internal/task",
		nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	answer := struct {
		Task task.Task `json:"task"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &answer.Task
}

func (client *Client) SendResult(result result.Result) {
	fmt.Println("res", result)
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("buf", buf.String())

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8081/internal/task",
		&buf,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println()
		return
	}
	defer resp.Body.Close()
}
