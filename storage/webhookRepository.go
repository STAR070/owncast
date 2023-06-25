package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/data"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type WebhookRepository struct {
	datastore *data.Store
}

func NewWebhookRepository(datastore *data.Store) *WebhookRepository {
	return &WebhookRepository{datastore: datastore}
}

var temporaryGlobalWebhooksInstance *WebhookRepository

// GetWebhookRepository returns the shared instance of the owncast datastore.
func GetWebhookRepository() *WebhookRepository {
	if temporaryGlobalWebhooksInstance == nil {
		temporaryGlobalWebhooksInstance = NewWebhookRepository(data.GetDatastore())
	}
	return temporaryGlobalWebhooksInstance
}

// InsertWebhook will add a new webhook to the database.
func (w *WebhookRepository) InsertWebhook(url string, events []models.EventType) (int, error) {
	log.Traceln("Adding new webhook")

	eventsString := strings.Join(events, ",")

	tx, err := w.datastore.DB.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO webhooks(url, events) values(?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(url, eventsString)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	newID, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(newID), err
}

// DeleteWebhook will delete a webhook from the database.
func (w *WebhookRepository) DeleteWebhook(id int) error {
	log.Traceln("Deleting webhook")

	tx, err := w.datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM webhooks WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		_ = tx.Rollback()
		return errors.New(fmt.Sprint(id) + " not found")
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetWebhooksForEvent will return all of the webhooks that want to be notified about an event type.
func (w *WebhookRepository) GetWebhooksForEvent(event models.EventType) []models.Webhook {
	webhooks := make([]models.Webhook, 0)

	query := `SELECT * FROM (
		WITH RECURSIVE split(id, url, event, rest) AS (
		  SELECT id, url, '', events || ',' FROM webhooks
		   UNION ALL
		  SELECT id, url,
				 substr(rest, 0, instr(rest, ',')),
				 substr(rest, instr(rest, ',')+1)
			FROM split
		   WHERE rest <> '')
		SELECT id, url, event
		  FROM split
		 WHERE event <> ''
	  ) AS webhook WHERE event IS "` + event + `"`

	rows, err := w.datastore.DB.Query(query)
	if err != nil || rows.Err() != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string

		if err := rows.Scan(&id, &url, &event); err != nil {
			log.Debugln(err)
			log.Error("There is a problem with the database.")
			break
		}

		singleWebhook := models.Webhook{
			ID:  id,
			URL: url,
		}

		webhooks = append(webhooks, singleWebhook)
	}
	return webhooks
}

// GetWebhooks will return all the webhooks.
func (w *WebhookRepository) GetWebhooks() ([]models.Webhook, error) { //nolint
	webhooks := make([]models.Webhook, 0)

	query := "SELECT * FROM webhooks"

	rows, err := w.datastore.DB.Query(query)
	if err != nil {
		return webhooks, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var events string
		var timestampString string
		var lastUsedString *string

		if err := rows.Scan(&id, &url, &events, &timestampString, &lastUsedString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return webhooks, err
		}

		timestamp, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return webhooks, err
		}

		var lastUsed *time.Time
		if lastUsedString != nil {
			lastUsedTime, _ := time.Parse(time.RFC3339, *lastUsedString)
			lastUsed = &lastUsedTime
		}

		singleWebhook := models.Webhook{
			ID:        id,
			URL:       url,
			Events:    strings.Split(events, ","),
			Timestamp: timestamp,
			LastUsed:  lastUsed,
		}

		webhooks = append(webhooks, singleWebhook)
	}

	if err := rows.Err(); err != nil {
		return webhooks, err
	}

	return webhooks, nil
}

// SetWebhookAsUsed will update the last used time for a webhook.
func (w *WebhookRepository) SetWebhookAsUsed(webhook models.Webhook) error {
	tx, err := w.datastore.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE webhooks SET last_used = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(webhook.ID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}