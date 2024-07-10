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
