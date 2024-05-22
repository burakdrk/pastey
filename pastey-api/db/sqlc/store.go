package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, options *sql.TxOptions, fn func(*Queries) error) error {
	const maxRetries = 5
	var err error

	for i := 0; i < maxRetries; i++ {
		tx, err := store.db.BeginTx(ctx, options)
		if err != nil {
			return err
		}

		q := New(tx)
		err = fn(q)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
			}

			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "40001" {
				continue
			}

			return err
		}

		err = tx.Commit()
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "40001" {
				continue
			}

			return err
		}

		return nil
	}

	return fmt.Errorf("transaction failed after %d retries: %w", maxRetries, err)
}

type SaveCopyParams struct {
	UserID       int64 `json:"user_id"`
	FromDeviceID int64 `json:"from_device_id"`
	Copies       []struct {
		ToDeviceID    int64  `json:"to_device_id"`
		EncryptedData string `json:"encrypted_data"`
	} `json:"copies"`
}

func (store *Store) SaveCopy(ctx context.Context, arg SaveCopyParams) ([]ClipboardEntry, error) {
	var result []ClipboardEntry

	err := store.execTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}, func(q *Queries) error {
		new_id := uuid.New()
		new_time := time.Now().UTC()

		entries, err := q.GetEntryByUserForUpdate(ctx, arg.UserID)
		if err != nil {
			return err
		}

		// fmt.Printf("Transaction %d: Found %d entries\n", tran, len(entries))
		// for _, entry := range entries {
		// 	fmt.Printf("Transaction %d: Entry %v\n", tran, entry.EntryID)
		// }

		user, err := q.GetUserById(ctx, arg.UserID)
		if err != nil {
			return err
		}

		var LIMIT int

		if user.Ispremium {
			LIMIT = 15
		} else {
			LIMIT = 2
		}

		// Delete the oldest entries.
		// ASSUMES ENTRIES ARE ORDERED BY CREATED_AT AND CREATED_AT IS THE SAME FOR DUPLICATES.
		entry_ids_to_keep := make(map[uuid.UUID]bool)
		deleted_entry_ids := make(map[uuid.UUID]bool)
		count := 0
		for _, entry := range entries {
			if _, ok := entry_ids_to_keep[entry.EntryID]; ok {
				continue
			}

			if count < LIMIT-1 {
				entry_ids_to_keep[entry.EntryID] = true
				count++
			} else {
				if _, ok := deleted_entry_ids[entry.EntryID]; ok {
					continue
				}

				err = q.DeleteEntry(ctx, entry.EntryID)

				//fmt.Printf("Transaction %d: Deleted entry %v\n", tran, entry.EntryID)

				if err != nil {
					return err
				}

				deleted_entry_ids[entry.EntryID] = true
			}
		}

		for _, copy := range arg.Copies {
			entry, err := q.CreateEntry(ctx, CreateEntryParams{
				EntryID:       new_id,
				UserID:        arg.UserID,
				FromDeviceID:  arg.FromDeviceID,
				ToDeviceID:    copy.ToDeviceID,
				EncryptedData: copy.EncryptedData,
				CreatedAt:     new_time,
			})
			if err != nil {
				return err
			}

			result = append(result, entry)
		}

		return nil
	})

	return result, err
}
