package DB

import (
	"bisgo/app/models"
	"bisgo/errlog"
	"database/sql"
	"errors"
)

// GetTransactionTypeById returns transaction type with the given id.
func (h *DBHandler) GetTransactionTypeById(transactionTypeId int) (models.NewTransactionType, error) {
	query := `SELECT Id, Code, Name FROM TransactionType WHERE Id = @p1`

	var transactionType models.NewTransactionType
	err := h.db.QueryRow(query,
		sql.Named("p1", transactionTypeId)).Scan(&transactionType.Id, &transactionType.Code, &transactionType.Name)
	if err != nil {
		errlog.Println(err)
		return models.NewTransactionType{}, errors.New("unsuccessful obtainance of transaction type")
	}

	return transactionType, nil
}

// GetTransactionTypeByCode returns transaction type with the given code.
func (h *DBHandler) GetTransactionTypeByCode(transactionTypeCode string) (models.NewTransactionType, error) {
	query := `SELECT Id, Code, Name FROM TransactionType WHERE Code = @p1`

	var transactionType models.NewTransactionType
	err := h.db.QueryRow(query,
		sql.Named("p1", transactionTypeCode)).Scan(&transactionType.Id, &transactionType.Code, &transactionType.Name)
	if err != nil {
		errlog.Println(err)
		return models.NewTransactionType{}, errors.New("unsuccessful obtainance of transaction type")
	}

	return transactionType, nil
}
