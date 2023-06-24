package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func main() {
	var wg sync.WaitGroup
	commentsCh := make(chan []Comment)

	// 控制并发请求数量
	concurrency := 10
	semaphore := make(chan struct{}, concurrency)

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // 获取一个信号量，限制并发数量

		go func(postID int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d/comments", postID)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error requesting URL:", err)
				return
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			var comments []Comment
			err = json.Unmarshal(body, &comments)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				return
			}

			commentsCh <- comments
		}(i)
	}

	go func() {
		wg.Wait()
		close(commentsCh)
	}()

	file, err := os.Create("go_email.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for comments := range commentsCh {
		for _, comment := range comments {
			fmt.Println("Post ID:", comment.PostID)
			fmt.Println("Comment ID:", comment.ID)
			fmt.Println("Name:", comment.Name)
			fmt.Println("Email:", comment.Email)
			fmt.Println("Body:", comment.Body)
			fmt.Println("--------------------")

			_, err := file.WriteString(comment.Email + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}
}
