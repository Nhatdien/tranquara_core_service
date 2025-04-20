package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserCompletedSelfGuideActivity struct {
	ActivityId      int64     `json:"activity_id"`
	UserId          uuid.UUID `json:"user_id"`
	ExerciseId      int       `json:"exercise_id"`
	DurationMinutes int       `json:"duration_minutes"`
	Notes           string    `json:"notes"`
	CompletedAt     time.Time `json:"completed_at"`
}

type UserCompletedSelfGuideActivityModel struct {
	DB *sql.DB
}

func (usa *UserCompletedSelfGuideActivityModel) Insert(completeSelfGuidActivity *UserCompletedSelfGuideActivity) error {
	query := `INSERT INTO user_self_guided_activities (user_id, exercise_id, duration_minutes, notes)
			 VALUES ($1, $2, $3, $4)
			 RETURNING *`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{completeSelfGuidActivity.UserId, completeSelfGuidActivity.ExerciseId, completeSelfGuidActivity.DurationMinutes, completeSelfGuidActivity.Notes}
	argsResponse := []any{&completeSelfGuidActivity.ActivityId, &completeSelfGuidActivity.UserId,
		&completeSelfGuidActivity.ExerciseId, &completeSelfGuidActivity.DurationMinutes,
		&completeSelfGuidActivity.Notes, &completeSelfGuidActivity.CompletedAt}

	return usa.DB.QueryRowContext(ctx, query, args...).Scan(argsResponse...)
}

func (a UserCompletedSelfGuideActivityModel) GetList(fromTime, toTime time.Time, userID uuid.UUID, filter Filter) ([]*UserCompletedSelfGuideActivity, Metadata, error) {
	query := fmt.Sprintf(`
					SELECT COUNT(*) OVER(), user_id, exercise_id, duration_minutes, notes, completed_at FROM user_self_guided_activities  
					WHERE completed_at BETWEEN $1 AND $2
					AND user_id = $3
					ORDER BY %s %s, activity_id DESC
					LIMIT $4 OFFSET $5
				`, filter.sortColumn(), filter.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	totalRecords := 0
	completedSelfGuideActivities := []*UserCompletedSelfGuideActivity{}

	rows, err := a.DB.QueryContext(ctx, query, fromTime, toTime, userID, filter.limit(), filter.offset())

	if err != nil {
		return nil, Metadata{}, err
	}
	for rows.Next() {
		var completedSelfGuideActivity UserCompletedSelfGuideActivity
		err = rows.Scan(
			&totalRecords,
			&completedSelfGuideActivity.UserId, &completedSelfGuideActivity.ExerciseId,
			&completedSelfGuideActivity.DurationMinutes, &completedSelfGuideActivity.CompletedAt,
			&completedSelfGuideActivity.Notes,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		completedSelfGuideActivities = append(completedSelfGuideActivities, &completedSelfGuideActivity)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)
	return completedSelfGuideActivities, metadata, nil
}
