package DB

import (
	"bisgo/app/models"
	"bisgo/errlog"
	"database/sql"
	"errors"
)

// GetBankJurisdiction returns the jurisdiction of the selected (bankId) bank.
func (h *DBHandler) GetBankJurisdiction(bankId string) (models.Jurisdiction, error) {
	query := `SELECT j.Id, j.Name FROM Bank b, Jurisdiction j WHERE b.JurisdictionId = j.Id 
				AND GlobalIdentifier = @p1`

	var jurisdiction models.Jurisdiction
	err := h.db.QueryRow(query, sql.Named("p1", bankId)).Scan(&jurisdiction.Id, &jurisdiction.Name)
	if err != nil {
		errlog.Println(err)
		return models.Jurisdiction{}, errors.New("unsuccessful obtainance of bank jurisdiction")
	}

	return jurisdiction, nil
}
