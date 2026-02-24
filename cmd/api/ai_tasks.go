package main

import (
	"time"

	"github.com/google/uuid"
	"tranquara.net/internal/data"
	"tranquara.net/internal/pubsub"
)

// AITaskMessage is the envelope for messages published to the ai_tasks queue.
type AITaskMessage struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}

// JournalIndexPayload is the payload sent to AI service for Qdrant indexing.
type JournalIndexPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	MoodScore *int      `json:"mood_score,omitempty"`
	MoodLabel *string   `json:"mood_label,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// publishJournalToAI publishes a journal to the ai_tasks queue for Qdrant indexing.
// This is non-blocking — failures are logged but don't affect the HTTP response.
func (app *application) publishJournalToAI(journal *data.UserJournal) {
	if app.rabbitchannel == nil {
		app.logger.PrintInfo("RabbitMQ not connected, skipping AI journal publish", nil)
		return
	}

	payload := JournalIndexPayload{
		ID:        journal.ID,
		UserID:    journal.UserID,
		Title:     journal.Title,
		Content:   journal.Content,
		MoodScore: journal.MoodScore,
		MoodLabel: journal.MoodLabel,
		CreatedAt: journal.CreatedAt,
	}

	message := AITaskMessage{
		Event:     "journal.index",
		Timestamp: time.Now(),
		Payload:   payload,
	}

	err := pubsub.PublishJson(app.rabbitchannel, "", "ai_tasks", message)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"action":     "publish_journal_to_ai",
			"journal_id": journal.ID.String(),
		})
		return
	}

	app.logger.PrintInfo("published journal to AI service", map[string]string{
		"journal_id": journal.ID.String(),
		"event":      "journal.index",
	})
}

// JournalDeletePayload is the payload sent to AI service to remove a journal from Qdrant.
type JournalDeletePayload struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

// publishJournalDeleteToAI publishes a journal delete event to the ai_tasks queue.
// This removes the journal's vector from Qdrant so deleted content is no longer used in RAG.
func (app *application) publishJournalDeleteToAI(journalID, userID uuid.UUID) {
	if app.rabbitchannel == nil {
		app.logger.PrintInfo("RabbitMQ not connected, skipping AI journal delete publish", nil)
		return
	}

	message := AITaskMessage{
		Event:     "journal.delete",
		Timestamp: time.Now(),
		Payload: JournalDeletePayload{
			ID:     journalID,
			UserID: userID,
		},
	}

	err := pubsub.PublishJson(app.rabbitchannel, "", "ai_tasks", message)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"action":     "publish_journal_delete_to_ai",
			"journal_id": journalID.String(),
		})
		return
	}

	app.logger.PrintInfo("published journal delete to AI service", map[string]string{
		"journal_id": journalID.String(),
		"event":      "journal.delete",
	})
}
