package DB

import (
	"bisgo/app/models"
	"bisgo/errlog"
	"database/sql"
	"errors"
)

// GetBankClientById returns the bank client with the given id.
func (h *DBHandler) GetBankClientById(id int) (models.NewBankClient, error) {
	query := `SELECT Id, GlobalIdentifier, Name, Address, BankId FROM BankClient WHERE Id = @p1`

	var bankClient models.NewBankClient
	err := h.db.QueryRow(query,
		sql.Named("p1", id)).Scan(&bankClient.Id,
		&bankClient.GlobalIdentifier,
		&bankClient.Name,
		&bankClient.Address,
		&bankClient.BankId)
	if err != nil {
		errlog.Println(err)
		return models.NewBankClient{}, errors.New("unsuccessful obtainance of compliance check")
	}

	return bankClient, nil
}

func (h *DBHandler) GetOrAddCumulativeAmount(clientID int) (int64, error) {
	returnErr := errors.New("retrieving cumulative amount failed")
	// First, try to retrieve the existing cumulative amount
	var cumulativeAmount int64

	query := "SELECT CumulativeAmount FROM AdditionalClientInfo WHERE ClientId = @p1"
	err := h.db.QueryRow(query, sql.Named("p1", clientID)).Scan(&cumulativeAmount)

	if err == nil {
		// Record found, return the cumulative amount
		return cumulativeAmount, nil
	}

	if err != sql.ErrNoRows {
		// An error occurred that wasn't just "no rows found"
		errlog.Println(err)
		return 0, returnErr
	}

	// Record not found, so it needs to be added
	query = "INSERT INTO AdditionalClientInfo (ClientId, CumulativeAmount) VALUES (@p1, 0)"
	_, err = h.db.Exec(query, sql.Named("p1", clientID))
	if err != nil {
		errlog.Println(err)
		return 0, returnErr
	}

	// Return the default value of 0 for the new record
	return 0, nil
}

func (h *DBHandler) UpdateCumulativeAmount(clientID int, newCumulativeAmount int64) error {
	query := `UPDATE [AdditionalClientInfo] SET CumulativeAmount = @p2 WHERE ClientId = @p1`

	_, err := h.db.Exec(query,
		sql.Named("p1", clientID),
		sql.Named("p2", newCumulativeAmount))
	if err != nil {
		errlog.Println(err)
		return errors.New("cumulative amount update failed")
	}

	return nil
}
