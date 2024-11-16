package repository

import (
	"database/sql"
	"errors"
	"tender-managment/internal/model"
)

type TenderRepository struct {
	db *sql.DB
}

func NewTenderRepository(db *sql.DB) *TenderRepository {
	return &TenderRepository{db: db}
}

func (r *TenderRepository) CreateTender(tender *model.Tender) (*model.Tender, error) {
	query := `
        INSERT INTO tenders (client_id, title, description, deadline, budget, status , attachment_path)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		tender.ClientID,
		tender.Title,
		tender.Description,
		tender.Deadline,
		tender.Budget,
		tender.Status,
		tender.Attachment,
	).Scan(&tender.ID, &tender.CreatedAt, &tender.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return tender, nil
}

func (r *TenderRepository) ListTendersByClientID(clientID int) ([]model.Tender, error) {
	query := `
        SELECT id, client_id, title, description, deadline, budget, status, created_at, updated_at
        FROM tenders
        WHERE client_id = $1`

	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []model.Tender
	for rows.Next() {
		var tender model.Tender
		err := rows.Scan(
			&tender.ID,
			&tender.ClientID,
			&tender.Title,
			&tender.Description,
			&tender.Deadline,
			&tender.Budget,
			&tender.Status,
			&tender.CreatedAt,
			&tender.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (r *TenderRepository) GetTenderByID(tenderID int) (*model.Tender, error) {
	query := `
        SELECT id, client_id, title, description, deadline, budget, status, created_at, updated_at
        FROM tenders
        WHERE id = $1`

	var tender model.Tender
	err := r.db.QueryRow(query, tenderID).Scan(
		&tender.ID,
		&tender.ClientID,
		&tender.Title,
		&tender.Description,
		&tender.Deadline,
		&tender.Budget,
		&tender.Status,
		&tender.CreatedAt,
		&tender.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("tender not found")
	}
	if err != nil {
		return nil, err
	}

	return &tender, nil
}

func (r *TenderRepository) UpdateTenderStatus(tenderID int, status string) error {
	query := `
        UPDATE tenders
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	result, err := r.db.Exec(query, status, tenderID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("tender not found")
	}

	return nil
}

func (r *TenderRepository) DeleteTender(tenderID int) error {
	query := `DELETE FROM tenders WHERE id = $1`

	result, err := r.db.Exec(query, tenderID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("tender not found")
	}

	return nil
}
