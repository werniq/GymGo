package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Exercise struct {
	Title     string `json:"title"`
	Technique string `json:"technique"`
	VideoURI  string `json:"VideoURI"`
}

type ExercisesResponse struct {
	Exercises []Exercise `json:"exercises"`
}

func (web *webapp) HomePage(c *gin.Context) {
	if err := web.renderTemplate(c.Writer, c.Request, "index", &templateData{}); err != nil {
		fmt.Printf("Error rendering INDEX page: %v", err)
	}
}

// Legs created request to server, and renders page with exercises data, received from server
func (web *webapp) Legs(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/legs/legs", nil)

	if err != nil {
		web.errorLog.Println("Error creating new request: %v", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		web.errorLog.Println("Error executing the request: %v", err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var exercises ExercisesResponse
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		web.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Legs Exercises"

	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering LEGS page: %v", err)
	}
}

func (web *webapp) Chest(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/front-upper/chest", nil)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		web.errorLog.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var exercises ExercisesResponse
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		web.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Chest Exercises"

	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering CHEST page: %v", err)
	}
}

func (web *webapp) Triceps(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/front-upper/Triceps", nil)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		web.errorLog.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var exercises ExercisesResponse
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		web.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Triceps Exercises"

	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering Triceps page: %v", err)
	}
}

func (web *webapp) Back(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/rear-upper/back", nil)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		web.errorLog.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var exercises ExercisesResponse
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		web.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Chest Exercises"

	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering Back page: %v", err)
	}
}

/*
	/rear-upper
		/triceps
		/shoulders
		/back
	/front-upper
		/biceps
		/chest
		/shoulders
	/legs
		/glutes
		/legs
*/

func (web *webapp) Biceps(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/rear-upper/biceps", nil)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		web.errorLog.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var exercises ExercisesResponse
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		web.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Biceps Exercises"

	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering BICEPS page: %v", err)
	}
}

// Shoulders
func (web *webapp) Shoulders(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/front-upper/chest", nil)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		web.errorLog.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var exercises ExercisesResponse
	err = json.Unmarshal(body, &exercises)
	if err != nil {
		web.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Chest Exercises"

	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering CHEST page: %v", err)
	}
}
