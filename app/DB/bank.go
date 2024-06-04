package DB

import (
	"bisgo/app/models"
	"database/sql"
	"log"
)

func (wrapper *DBHandler) Login(username, password string) *models.BankEmployeeModel {
	query := `SELECT be.Name, be.Username, be.Password, be.BankId, b.Name BankName, c.Name Country
				FROM [dbo].[BankEmployee] be
				JOIN [dbo].[Bank] b ON be.BankId = b.GlobalIdentifier
				JOIN [dbo].[Country] c ON b.CountryId = c.Id
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
	query := `SELECT b.GlobalIdentifier, b.Name, b.CountryId
          FROM Bank b
          WHERE b.GlobalIdentifier = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bank models.NewBank

	for rows.Next() {
		if err := rows.Scan(&bank.GlobalIdentifier, &bank.Name, &bank.CountryId); err != nil {
			log.Fatal(err)
		}
	}

	return bank
}

func (wrapper *DBHandler) GetBanks() []models.BankModel {
	query := `SELECT b.GlobalIdentifier, b.Name, c.Name CountryName
          FROM Bank b
          JOIN Country c ON b.CountryId = c.Id`

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

func (wrapper *DBHandler) GetCountry(countryId uint) models.NewCountry {
	query := `SELECT c.Id, c.Name, c.Code
					From Country c
					WHERE c.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", countryId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var country models.NewCountry

	for rows.Next() {
		if err := rows.Scan(&country.Id, &country.Name, &country.Code); err != nil {
			log.Fatal(err)
		}
	}

	return country
}
