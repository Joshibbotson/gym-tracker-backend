package workout

import (
	"context"
	"fmt"
	"strconv"
	"time"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutService interface {
	CreateWorkout(userID primitive.ObjectID, workout t.CreateWorkoutRequest) (*t.Workout, error)
	GetWorkoutsByUserId(userID primitive.ObjectID) ([]t.YearlyData, error)
	GetWorkoutsByDate(userID primitive.ObjectID, date time.Time) ([]t.Workout, error)
	UpdateWorkout(userID primitive.ObjectID, workout t.UpdateWorkoutRequest) ([]t.Workout, error)
	DeleteWorkout(workoutID primitive.ObjectID) (bool, error)
}

type workoutService struct {
	repo WorkoutRepository
}

func NewWorkoutService(repo WorkoutRepository) WorkoutService {
	return &workoutService{repo: repo}
}

func (s *workoutService) CreateWorkout(userID primitive.ObjectID, workout t.CreateWorkoutRequest) (*t.Workout, error) {
	config := t.WorkoutConfig{
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
		UserId:    userID,
		Date:      workout.Date,
		Workout:   &config,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.InsertWorkout(context.TODO(), newWorkout)
}

func (s *workoutService) GetWorkoutsByDate(userId primitive.ObjectID, date time.Time) ([]t.Workout, error) {
	return s.repo.FetchWorkoutByDate(context.TODO(), userId, date)

}

func (s *workoutService) GetWorkoutsByUserId(userID primitive.ObjectID) ([]t.YearlyData, error) {
	workouts, err := s.repo.FetchWorkoutsByUserId(context.TODO(), userID)
	if err != nil {
		return nil, err
	}

	return fillMissingDates(workouts), nil
}

func fillMissingDates(workoutData []t.YearlyData) []t.YearlyData {
	location := time.UTC

	if len(workoutData) == 0 {
		year := time.Now().Year()
		// Initialize a new YearlyData with all months and days for the year
		workoutData = []t.YearlyData{{Year: year, Months: []t.MonthlyData{}}}
	}

	for i, year := range workoutData {
		yearInt := year.Year
		// Create a map for existing months for quick lookup
		existingMonths := make(map[int][]t.DailyWorkout)
		for _, month := range year.Months {
			monthNumber, err := strconv.Atoi(month.Month)
			if err != nil {
				fmt.Printf("Error parsing month: %v\n", err)
				continue
			}
			existingMonths[monthNumber] = month.Workouts
		}

		// Fill all months for the year
		var filledMonths []t.MonthlyData
		for monthNumber := 1; monthNumber <= 12; monthNumber++ {
			// Get the first and last day of the month
			firstDay := time.Date(yearInt, time.Month(monthNumber), 1, 0, 0, 0, 0, location)
			lastDay := firstDay.AddDate(0, 1, -1)

			// Create a map of all existing workout dates for this month
			existingDates := make(map[string]bool)
			for _, workout := range existingMonths[monthNumber] {
				existingDates[workout.Date.Format("2006-01-02")] = true
			}

			// Create a full list of dates for the month
			var filledWorkouts []t.DailyWorkout
			for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
				dateStr := d.Format("2006-01-02")
				if existingDates[dateStr] {
					// Add all existing workouts for this date
					for _, workout := range existingMonths[monthNumber] {
						if workout.Date.Format("2006-01-02") == dateStr {
							filledWorkouts = append(filledWorkouts, workout)
						}
					}
				} else {
					// Add placeholder for missing date
					filledWorkouts = append(filledWorkouts, t.DailyWorkout{
						Date:   d,
						Config: nil,
					})
				}
			}

			// Add the filled month to the result
			filledMonths = append(filledMonths, t.MonthlyData{
				Month:    fmt.Sprintf("%02d", monthNumber), // Ensure consistent month format
				Workouts: filledWorkouts,
			})
		}

		// Update the year with the filled months
		workoutData[i].Months = filledMonths
	}

	return workoutData
}

func (s *workoutService) UpdateWorkout(userID primitive.ObjectID, workout t.UpdateWorkoutRequest) ([]t.Workout, error) {
	config := t.WorkoutConfig{
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

	updatedWorkout := t.Workout{
		ID:      workout.ID,
		Date:    workout.Date,
		UserId:  userID,
		Workout: &config,
	}

	data, err := s.repo.UpdateWorkout(context.TODO(), updatedWorkout)
	if err != nil {
		fmt.Println("error updating!")
		// probably should return the workouts by date here instead of nil? Or maybe not tbf
		return nil, err
	}
	return s.repo.FetchWorkoutByDate(context.TODO(), userID, data.Date)
}

func (s *workoutService) DeleteWorkout(workoutID primitive.ObjectID) (bool, error) {
	return s.repo.RemoveWorkout(context.TODO(), workoutID)
}
