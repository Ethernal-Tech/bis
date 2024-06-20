package DB

import (
	"bisgo/app/models"
	"database/sql"
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

func (wrapper *DBHandler) GetBankId(bankName string) string {
	query := `SELECT GlobalIdentifier FROM [Bank] WHERE name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankName))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bankId string
	for rows.Next() {
		if err := rows.Scan(&bankId); err != nil {
			log.Fatal(err)
		}
	}

	return bankId
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

// GetOrCreateClient method to get client Id by GlobalIdentifier and Name or create a new client if not exists
func (wrapper *DBHandler) GetOrCreateClient(globalIdentifier, name, address, bankId string) uint {
	// Check if client exists
	query := `SELECT Id FROM BankClient WHERE GlobalIdentifier = @p1 AND Name = @p2`

	var clientId uint
	err := wrapper.db.QueryRow(query, sql.Named("p1", globalIdentifier), sql.Named("p2", name)).Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Client does not exist, insert new client
			insertQuery := `INSERT INTO [dbo].[BankClient] OUTPUT INSERTED.Id VALUES (@p1, @p2, @p3, @p4)`
			err = wrapper.db.QueryRow(insertQuery, sql.Named("p1", globalIdentifier), sql.Named("p2", name), sql.Named("p3", address), sql.Named("p4", bankId)).Scan(&clientId)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Other errors
			log.Fatal(err)
		}
	}

	return clientId
}
