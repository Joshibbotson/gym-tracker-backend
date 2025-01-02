package workout

import (
	"context"
	"fmt"
	"strconv"
	"time"

	db "github.com/joshibbotson/gym-tracker-backend/internal/db"
	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutService interface {
	CreateWorkout(userID primitive.ObjectID, workout t.CreateWorkoutRequest) (*t.Workout, error)
	GetWorkoutsByUserId(userId primitive.ObjectID) ([]t.YearlyData, error)
	// getWorkoutsByDate(userId string, date time.Time) (t.Workout, error)
}

type workoutService struct{}

func NewWorkoutService() WorkoutService {
	return &workoutService{}
}

func (r *workoutService) CreateWorkout(userID primitive.ObjectID, workout t.CreateWorkoutRequest) (*t.Workout, error) {
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
		UserId:    userID,
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

// in the future add a year to get
// it also doesn't currently get multiple workouts on the same day.
// we need to include the ID
func (r *workoutService) GetWorkoutsByUserId(userId primitive.ObjectID) ([]t.YearlyData, error) {
	collection := db.Client.Database(db.DB_NAME).Collection("workout")

	startOfYear := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(time.Now().Year(), time.December, 31, 23, 59, 59, 999999999, time.UTC)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "userId", Value: userId},
			{Key: "date", Value: bson.D{
				{Key: "$gte", Value: startOfYear},
				{Key: "$lte", Value: endOfYear},
			}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "ID", Value: "$_id"},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$year", Value: "$date"}}},
			{Key: "month", Value: bson.D{{Key: "$month", Value: "$date"}}},
			{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$date"}}},
			{Key: "date", Value: "$date"},
			{Key: "config", Value: `$workout`},
			{Key: "ID", Value: 1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "year", Value: "$year"}, {Key: "month", Value: "$month"}}},
			{Key: "workouts", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "ID", Value: "$ID"},
				{Key: "date", Value: "$date"},
				{Key: "config", Value: "$config"},
			}}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.year"},
			{Key: "months", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "month", Value: "$_id.month"},
				{Key: "workouts", Value: "$workouts"},
			}}}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "months.month", Value: 1},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "year", Value: `$_id`},
			{Key: "months", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$months"},
				{Key: "as", Value: "month"},
				{Key: "in", Value: bson.D{
					{Key: "ID", Value: "$ID"},
					{Key: "month", Value: bson.D{{Key: "$toString", Value: "$$month.month"}}},
					{Key: "workouts", Value: `$$month.workouts`},
				}},
			}}}},
		}}},
	}

	// Execute the aggregation pipeline
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation error: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Decode the results into a slice of YearlyData
	var results []t.YearlyData
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("decoding error: %v", err)
	}

	return fillMissingDates(results), nil
}

func fillMissingDates(workoutData []t.YearlyData) []t.YearlyData {
	location := time.UTC

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

// Should be able to get any workouts related to a date clicked and userId
// can we retrieve the userID from the session and add it to context?
// func (r *workoutService) getWorkoutsByDate(userId string, date time.Time) (t.Workout, error) {

// }
