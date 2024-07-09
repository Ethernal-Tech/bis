package DB

import (
	"bisgo/app/models"
	"bisgo/errlog"
	"database/sql"
	"errors"
	"log"
)

func (wrapper *DBHandler) Login(username, password string) *models.BankEmployeeModel {
	query := `SELECT be.Name, be.Username, be.Password, be.BankId, b.Name BankName, j.Name Jurisdiction
				FROM [dbo].[BankEmployee] be
				JOIN [dbo].[Bank] b ON be.BankId = b.GlobalIdentifier
				JOIN [dbo].[Jurisdiction] j ON b.JurisdictionId = j.Id
			  	WHERE be.Username = @p1 AND be.Password = @p2`

	rows, err := wrapper.db.Query(query, sql.Named("p1", username), sql.Named("p2", password))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var user models.BankEmployeeModel
		if err := rows.Scan(&user.Name, &user.Username, &user.Password, &user.BankId, &user.BankName, &user.Country); err != nil {
			log.Println("Error scanning row:", err)
			return nil
		}
		return &user
	}

	return nil
}

func (wrapper *DBHandler) IsCentralBankEmployee(username string) bool {
	query := `SELECT CASE
                  WHEN EXISTS (
                    SELECT 1
                    FROM BankEmployee AS be
                    JOIN Bank AS b ON b.GlobalIdentifier = be.BankId
                    WHERE be.Username = @p1 AND b.BankTypeId = @p2
                  )
                  THEN 'true'
                  ELSE 'false'
              END AS CentralBankEmployee;`

	rows, err := wrapper.db.Query(query, sql.Named("p1", username), sql.Named("p2", models.CentralBank))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var CentralBankEmployee bool = false
	for rows.Next() {
		if err = rows.Scan(&CentralBankEmployee); err != nil {
			log.Println("Error scanning row:", err)
			return false
		}
	}

	return CentralBankEmployee
}

// GetBankIdByName returns the id (global identifier) of the selected (bankName) bank.
func (h *DBHandler) GetBankIdByName(bankName string) (string, error) {
	query := `SELECT GlobalIdentifier FROM Bank WHERE Name = @p1`

	var bankId string
	err := h.db.QueryRow(query).Scan(&bankId)
	if err != nil {
		errlog.Println(err)
		return "", errors.New("unsuccessful obtainance of bank id")
	}

	return bankId, nil
}

func (wrapper *DBHandler) GetBankClientId(bankClientName string) string {
	query := `SELECT GlobalIdentifier FROM [BankClient] WHERE Name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankClientName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bankClientId string
	for rows.Next() {
		if err := rows.Scan(&bankClientId); err != nil {
			log.Fatal(err)
		}
	}

	return bankClientId
}

func (wrapper *DBHandler) GetBank(bankId string) models.NewBank {
	query := `SELECT b.GlobalIdentifier, b.Name, b.JurisdictionId
          FROM Bank b
          WHERE b.GlobalIdentifier = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bank models.NewBank

	for rows.Next() {
		if err := rows.Scan(&bank.GlobalIdentifier, &bank.Name, &bank.JurisdictionId); err != nil {
			log.Fatal(err)
		}
	}

	return bank
}

func (wrapper *DBHandler) GetBanks() []models.BankModel {
	query := `SELECT b.GlobalIdentifier, b.Name, j.Name JurisdictionName
          FROM Bank b
          JOIN Jurisdiction j ON b.JurisdictionId = j.Id`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	banks := []models.BankModel{}
	for rows.Next() {
		var bank models.BankModel
		if err := rows.Scan(&bank.Id, &bank.Name, &bank.Country); err != nil {
			log.Fatal(err)
		}
		banks = append(banks, bank)
	}
	return banks
}

func (wrapper *DBHandler) GetClientNameByID(clientID uint) string {
	query := `SELECT Name FROM BankClient WHERE Id = @p1`

	var clientName string
	err := wrapper.db.QueryRow(query, sql.Named("p1", clientID)).Scan(&clientName)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Fatal(err)
	}

	return clientName
}

func (wrapper *DBHandler) GetClientByID(clientID uint) models.NewBankClient {
	query := `SELECT Id, GlobalIdentifier, Name, Address, BankId FROM BankClient WHERE Id = @p1`

	var client models.NewBankClient
	err := wrapper.db.QueryRow(query, sql.Named("p1", clientID)).Scan(&client.Id, &client.GlobalIdentifier, &client.Name, &client.Address, &client.BankId)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.NewBankClient{}
		}
		log.Fatal(err)
	}

	return client
}

// CreateOrGetBankClient creates a new bank client and returns its id.
// If bank client already exists, method only returns id.
func (h *DBHandler) CreateOrGetBankClient(globalIdentifier, name, address, bankId string) (int, error) {
	returnErr := errors.New("unsuccessful creation/obtainance of bank client/its id")

	// query to check if the bank client already exists in the system
	query := `SELECT Id FROM BankClient WHERE GlobalIdentifier = @p1 AND Name = @p2`

	var bankClientId int
	err := h.db.QueryRow(query,
		sql.Named("p1", globalIdentifier),
		sql.Named("p2", name)).Scan(&bankClientId)

	// if policy type (row) doesn't exist, error [sql.ErrNoRows] appears
	if err != nil {
		// check for potentially some other type of error
		if err != sql.ErrNoRows {
			errlog.Println(err)
			return -1, returnErr
		}

		// query to instert a new bank client and get its id
		query := `INSERT INTO BankClient (GlobalIdentifier, Name, Adress, BankId) OUTPUT INSERTED.Id VALUES (@p1, @p2, @p3, @p4)`
		err = h.db.QueryRow(query,
			sql.Named("p1", globalIdentifier),
			sql.Named("p2", name),
			sql.Named("p3", address),
			sql.Named("p4", bankId)).Scan(&bankClientId)
		if err != nil {
			errlog.Println(err)
			return -1, returnErr
		}
	}

	return bankClientId, nil
}
