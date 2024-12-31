package workout

import (
	"context"
	"time"

	db "github.com/joshibbotson/gym-tracker-backend/internal/db"
	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutService interface {
	createWorkout(workout t.CreateWorkoutRequest) (*t.Workout, error)
	GetWorkoutsByUserId(userId string) (t.WorkoutData, error)
	getWorkoutsByDate(userId string, date time.Time) (t.Workout, error)
}

type workoutService struct{}

func NewWorkoutService() WorkoutService {
	return &workoutService{}
}

func (r *workoutService) createWorkout(workout t.CreateWorkoutRequest) (*t.Workout, error) {
	collection := db.Client.Database(db.DB_NAME).Collection("workout")
	Config := t.WorkoutConfig{
		Weight:       workout.Weight,
		WorkoutType:  workout.WorkoutType,
		CaloriePhase: workout.CaloriePhase,
		ChestSize:    workout.ChestSize,
		WaistSize:    workout.WaistSize,
		BicepSize:    workout.BicepSize,
		ForearmSize:  workout.ForearmSize,
		ThighSize:    workout.ThighSize,
		CalfSize:     workout.CalfSize,
	}

	newWorkout := t.Workout{
		ID:        primitive.NewObjectID(),
		Date:      workout.Date,
		Workout:   &Config,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), newWorkout)
	if err != nil {
		return nil, err
	}

	return &newWorkout, nil
}

// need to create an aggregation here to capture workouts data.
func (r *workoutService) GetWorkoutsByUserId(userId string) (t.WorkoutData, error) {

}

// Should be able to get any workouts related to a date clicked and userId
// can we retrieve the userID from the decision
func (r *workoutService) getWorkoutsByDate(userId string, date time.Time) (t.Workout, error) {

}
