package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/werniq/gym/models"
	"math/rand"
	"net/http"
	"strings"
)

const (
	back   = 16
	chest  = 17
	biceps = 15
	legs   = 15
	glutes = 13
)

var (
	muscles = [7]string{
		"back", "chest", "biceps", "triceps", "abs", "legs", "glutes",
	}
)

// ReturnOneExercise gets one exercise from table given in r.Body, and
// writes exercise in json format in response. Will be used for re-generating
// exercise. Need to create a button, in lower part of exercise card, and if clicked
// re-generating exercise.
func (web *webapp) ReturnOneExercise(c *gin.Context) {
	var req struct {
		Table string `json:"table"`
		Id    int    `json:"current_exercise_id"`
	}
	json.NewDecoder(c.Request.Body).Decode(&req)
	var num1 int
	var exercise models.Exercise
	var err error
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	rand.Seed(legs)

	for {
		num1 = rand.Intn(legs + 1)
		if num1 != req.Id {
			break
		}
	}

	exercise, err = web.db.GetExerciseById(num1, req.Table)
	if err != nil {
		payload.Error = true
		payload.Message = "Error getting exercise. Please, try again."
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   payload.Error,
			"message": payload.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exercise": exercise,
	})
}

// ReturnExercises returns amount of exercises, inputted by user.
func (web *webapp) ReturnExercises(c *gin.Context) {
	var info struct {
		TableName     string `json:"tableName"`
		ExerciseCount int    `json:"exercise_count"`
	}

	var exercises []models.Exercise
	var exercisesID []int
	json.NewDecoder(c.Request.Body).Decode(&info.ExerciseCount)
	if info.ExerciseCount <= 0 {
		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = true
		payload.Message = "Count of exercises in workout should be more than 0"
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   payload.Error,
			"message": payload.Message,
		})
		return
	}
	for i := 0; i < info.ExerciseCount; i++ {
		var num1 int
		for {
			num1 = rand.Intn(legs + 1)
			if web.elementInArray(i, exercisesID, num1) == false {
				break
			}
		}
		if ex, err := web.db.GetExerciseById(num1, info.TableName); err == nil {
			exercises = append(exercises, ex)
			exercisesID = append(exercisesID, num1)
		} else {
			var payload struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}
			payload.Error = true
			payload.Message = "Actually, idgaf what the error is"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   payload.Error,
				"message": payload.Message,
			})
			return
		}
	}
	// at this point should not be any errors, so we can send exercises to front-end
	c.JSON(http.StatusOK, gin.H{
		"exercises": exercises,
	})
	fmt.Println(exercises)
}

// GenerateWorkout handler will return whole workout
// for given muscle and amount of exercises, also given by user
func (web *webapp) GenerateWorkout(c *gin.Context) {
	var payload struct {
		ExercisesCount []int  `json:"exercises_count"`
		Muscles        string `json:"muscles"`
	}

	var response struct {
		Muscles   []string          `json:"muscles"`
		Exercises []models.Exercise `json:"exercises"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		web.errorLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Error decoding request body: %v", err),
		})
		return
	}:

	muscles := strings.Split(payload.Muscles, "")
	response.Muscles = muscles

	for i := 0; i <= len(payload.ExercisesCount)-1; i++ {
		if payload.ExercisesCount[i] <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "exercises count can not be less or equal 0",
			})
		}
	}

	// generating exercises for workout (get random once from table)
	for i := 0; i < len(payload.ExercisesCount); i++ {
		ex, err := web.db.GetOneRandomExercise(response.Muscles[i])
		if err != nil {
			var info struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}
			info.Error = true
			info.Message = "something went wrong. please, try again"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   info.Error,
				"message": info.Message,
			})
			return
		}
		response.Exercises = append(response.Exercises, ex)
	}

	//out, _ := json.MarshalIndent(payload, "", "\t")
	c.JSON(http.StatusOK, gin.H{
		"exercises": response.Exercises,
		"muscles":   response.Muscles,
	})
}

// ReturnAllBackExercises is used in /api/back route, for receiving all exercises for back
func (web *webapp) ReturnAllBackExercises(c *gin.Context) {
	stmt := `SELECT * FROM back`
	var back []models.Exercise

	row, err := web.db.DB.Query(stmt)

	if err != nil {
		fmt.Printf("Error selecting EXERCISES from BACK table: %v", err)
		return
	}

	err = row.Scan(&back)

	c.JSON(http.StatusOK, gin.H{
		"exercises": back,
	})

}

// ReturnAllLegsExercises is used in /api/legs route, for receiving all exercises for legs
func (web *webapp) ReturnAllLegsExercises(c *gin.Context) {
	var exercises []models.Exercise
	var legs []models.Exercise

	stmt := `SELECT * FROM legs`

	row, err := web.db.DB.Query(stmt)
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	for row.Next() {
		var title string
		var technique string
		var videoURI string
		var ex models.Exercise
		if err := row.Scan(&title, &technique, &videoURI); err != nil {
			web.errorLog.Println(err)
		}
		ex = models.Exercise{
			Title:     title,
			Technique: technique,
			VideoURI:  videoURI,
		}
		legs = append(legs, ex)
	}

	out, err := json.MarshalIndent(legs, "", "\t")
	if err != nil {
		web.errorLog.Println(err)
		return
	}

	json.Unmarshal(out, &exercises)

	c.JSON(http.StatusOK, gin.H{
		"exercises": exercises,
	})
}

// ReturnAllChestExercises is used in /api/chest route, for receiving all exercises for chest
func (web *webapp) ReturnAllChestExercises(c *gin.Context) {
	stmt := `SELECT * FROM chest`
	var chest []models.Exercise

	row, err := web.db.DB.Query(stmt)

	if err != nil {
		fmt.Printf("Error selecting EXERCISES from CHEST table: %v", err)
		return
	}

	err = row.Scan(&chest)

	c.JSON(http.StatusOK, gin.H{
		"exercises": chest,
	})
}

// ReturnAllBicepsExercises is used in /api/biceps route, for receiving all exercises for biceps
func (web *webapp) ReturnAllBicepsExercises(c *gin.Context) {
	stmt := `SELECT * FROM biceps`
	var biceps []models.Exercise

	row, err := web.db.DB.Query(stmt)

	if err != nil {
		fmt.Printf("Error selecting EXERCISES from BICEPS table: %v", err)
		return
	}

	err = row.Scan(&biceps)

	c.JSON(http.StatusOK, gin.H{
		"exercises": biceps,
	})
}

// ReturnAllTricepsExercises is used in /api/triceps route, for receiving all exercises for triceps
func (web *webapp) ReturnAllTricepsExercises(c *gin.Context) {
	stmt := `SELECT * FROM triceps`
	var triceps []models.Exercise

	row, err := web.db.DB.Query(stmt)

	if err != nil {
		fmt.Printf("Error selecting EXERCISES from TRICEPS table: %v", err)
		return
	}

	err = row.Scan(&triceps)

	c.JSON(http.StatusOK, gin.H{
		"exercises": triceps,
	})
}

// ReturnAllShouldersExercises is used in /api/shoulders route, for receiving all exercises for shoulders
func (web *webapp) ReturnAllShouldersExercises(c *gin.Context) {
	stmt := `SELECT * FROM shoulders`
	var shoulders []models.Exercise

	row, err := web.db.DB.Query(stmt)

	if err != nil {
		fmt.Printf("Error selecting EXERCISES from SHOULDERS table: %v", err)
		return
	}

	err = row.Scan(&shoulders)

	c.JSON(http.StatusOK, gin.H{
		"exercises": shoulders,
	})
}
