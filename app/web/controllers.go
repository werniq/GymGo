package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/werniq/gym/models"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Exercise struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Technique string `json:"technique"`
	VideoURI  string `json:"VideoURI"`
}

type ExercisesResponse struct {
	Exercises []models.Exercise `json:"exercises"`
}

func (web *webapp) HomePage(c *gin.Context) {
	if err := web.renderTemplate(c.Writer, c.Request, "index", &templateData{}); err != nil {
		fmt.Printf("Error rendering INDEX page: %v", err)
	}
}

func (web *webapp) GenerateWorkout(c *gin.Context) {

	data := make(map[string]interface{})
	var payload struct {
		Muscles       string   `json:"muscles"`
		ExerciseCount []string `json:"exercisesCount"`
	}

	var request struct {
		Muscles       string `json:"muscles"`
		ExerciseCount []int  `json:"exercisesCount"`
	}

	var res1 ExercisesResponse

	var info struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		web.errorLog.Println(err)
		info.Error = true
		info.Message = err.Error()
		fmt.Println(c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   info.Error,
			"message": info.Message,
		})
		return
	}

	if len(payload.ExerciseCount) == 0 {
		info.Error = true
		info.Message = "you should input amount of exercises in form"
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   info.Error,
			"message": info.Message,
		})
		return
	} else {
		for i := 0; i < len(payload.ExerciseCount); i++ {
			if payload.ExerciseCount[i] == "" {
				info.Error = true
				info.Message = fmt.Sprintf("Cannot format empty string on index: %d", i)
				return
			}
			num, err := strconv.Atoi(payload.ExerciseCount[i])
			if err != nil {
				info.Error = true
				info.Message = fmt.Sprintf("Error converting string to number: %v", err)
				web.errorLog.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   info.Error,
					"message": info.Message,
				})
				return
			}
			request.ExerciseCount = append(request.ExerciseCount, num)
		}
	}
	request.Muscles = payload.Muscles

	body, err := json.Marshal(request)
	if err != nil {
		info.Error = true
		info.Message = fmt.Sprintf("Error marshalling data: %v", err)
		web.errorLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   info.Error,
			"message": info.Message,
		})
		return
	}
	reader := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://127.0.0.1:4001/api/generate-workout", reader)
	if err != nil {
		info.Error = true
		info.Message = fmt.Sprintf("Error creating new POST request: %v", err)
		web.errorLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   info.Error,
			"message": info.Message,
		})
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		info.Error = true
		info.Message = fmt.Sprintf("Error executing the request: %v", err)
		web.errorLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   info.Error,
			"message": info.Message,
		})
		return
	}

	err = json.NewDecoder(res.Body).Decode(&res1)
	if err != nil {
		info.Error = true
		info.Message = fmt.Sprintf("Error decoding body: %v", err)
		web.errorLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   info.Error,
			"message": info.Message,
		})
		return
	}

	data["exercises"] = res1
	data["title"] = "WORKOUT! WORKOUT!"

	body, err = json.Marshal(res1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Error marshalling data: %v", err),
		})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	// id, title, technique, videoURI
	rand.Seed(time.Now().Unix())
	//num := rand.Intn(100000000)

	_, _, err = web.db.SaveExercises(res1.Exercises)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	//session := sessions.Default(c)

	//session.Set("len", len(res1.Exercises))
	//session.Set("startID", num)
	//err = session.Save()
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error":   true,
	//		"message": fmt.Sprintf("Error saving session data: %v", err),
	//	})
	//	return
	//}

	c.Redirect(301, "/receipt")
	c.Abort()
}

func (web *webapp) Receipt(c *gin.Context) {
	s := sessions.Default(c)
	id := s.Get("startID")
	l := s.Get("len")

	num1 := strconv.Itoa(id.(int))
	exercises, err := web.db.RetrieveExercises(num1, l.(int))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Error decoding RECEIPT body: %v", err),
		})
		return
	}
	data := make(map[string]interface{})
	data["title"] = "Your workout, please!"
	data["exercises"] = exercises

	c.Redirect(http.StatusFound, "/receipt/muscle")
}

func (web *webapp) ReceiptMuscle(c *gin.Context) {
	var res ExercisesResponse
	err := json.NewDecoder(c.Request.Body).Decode(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Error decoding Request body: %v", err),
		})
		return
	}
	data := make(map[string]interface{})

	data["title"] = "Your "
	data["exercises"] = res.Exercises

	if err := web.renderTemplate(c.Writer, c.Request, "receipt", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println(err)
	}
}

func (web *webapp) GenerateWorkoutPage(c *gin.Context) {
	data := make(map[string]interface{})
	data["title"] = "Generate workout"

	if err := web.renderTemplate(c.Writer, c.Request, "generate-workout", &templateData{
		Data: data,
	}, "workout"); err != nil {
		fmt.Printf("Error rendering GENERATE-WORKOUT page: %v", err)
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
		return
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
		return
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

func (web *webapp) Glutes(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://localhost:4001/api/legs/glutes", nil)

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
		return
	}

	data := make(map[string]interface{})
	data["exercises"] = exercises.Exercises
	data["title"] = "All Legs Exercises"
	fmt.Println("alright")
	if err := web.renderTemplate(c.Writer, c.Request, "muscles", &templateData{
		Data: data,
	}); err != nil {
		web.errorLog.Println("Error rendering LEGS page: %v", err)
	}
}

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
