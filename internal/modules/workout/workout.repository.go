package workout

import (
	"context"
	"fmt"
	"time"

	db "github.com/joshibbotson/gym-tracker-backend/internal/db"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutRepository interface {
	InsertWorkout(ctx context.Context, workout t.Workout) (*t.Workout, error)
	FetchWorkoutByDate(ctx context.Context, userId primitive.ObjectID, date time.Time) ([]t.Workout, error)
	FetchWorkoutsByUserId(ctx context.Context, userID primitive.ObjectID) ([]t.YearlyData, error)
	FetchActivityCountByUserId(ctx context.Context, userID primitive.ObjectID) (int64, error)
	UpdateWorkout(ctx context.Context, workout t.Workout) (*t.Workout, error)
	RemoveWorkout(ctx context.Context, workoutID primitive.ObjectID) (bool, error)
}

type workoutRepository struct {
	workoutCollection *mongo.Collection
}

func NewWorkoutRepository() WorkoutRepository {
	return &workoutRepository{
		workoutCollection: db.Client.Database(db.DB_NAME).Collection("workout"),
	}
}

func (r *workoutRepository) InsertWorkout(ctx context.Context, workout t.Workout) (*t.Workout, error) {
	_, err := r.workoutCollection.InsertOne(ctx, workout)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *workoutRepository) FetchWorkoutByDate(ctx context.Context, userId primitive.ObjectID, date time.Time) ([]t.Workout, error) {
	// Calculate the start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Define the query filter
	filter := bson.M{
		"userId": userId,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	// Query the database
	var workouts []t.Workout
	// returns a cursor to iterate through results
	// if error with filter returns error
	cursor, err := r.workoutCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// ensures the cursor which is streaming from mongo is closed
	// when the function returns at any point.
	defer cursor.Close(ctx)

	// Decode results
	// goes through each documents and decodes it to struct
	if err = cursor.All(ctx, &workouts); err != nil {
		return nil, err
	}
	fmt.Println("workouts:", workouts)
	return workouts, nil
}

func (r *workoutRepository) FetchWorkoutsByUserId(ctx context.Context, userID primitive.ObjectID) ([]t.YearlyData, error) {
	startOfYear := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(time.Now().Year(), time.December, 31, 23, 59, 59, 999999999, time.UTC)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "userId", Value: userID},
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

	cursor, err := r.workoutCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []t.YearlyData
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decoding error: %v", err)
	}

	return results, nil
}

func (r *workoutRepository) FetchActivityCountByUserId(ctx context.Context, userId primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"userId": userId,
	}
	count, err := r.workoutCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *workoutRepository) UpdateWorkout(ctx context.Context, workout t.Workout) (*t.Workout, error) {
	_, err := r.workoutCollection.UpdateByID(ctx, workout.ID, bson.M{"$set": workout})
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *workoutRepository) RemoveWorkout(ctx context.Context, workoutID primitive.ObjectID) (bool, error) {
	res, err := r.workoutCollection.DeleteOne(ctx, bson.M{"_id": workoutID})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
