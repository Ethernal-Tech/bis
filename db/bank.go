package db

import (
	"bisgo/common"
	"database/sql"
	"log"
)

func (wrapper *DBWrapper) Login(username, password string) *BankEmployeeModel {
	query := `SELECT [BankEmployee].Name Name, Username, Password, BankId, [Bank].Name BankName
			  FROM [dbo].[BankEmployee], [dbo].[Bank]
			  WHERE BankId = [Bank].Id AND Username = @p1 AND Password = @p2`

	rows, err := wrapper.db.Query(query, sql.Named("p1", username), sql.Named("p2", password))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var user BankEmployeeModel
		if err := rows.Scan(&user.Name, &user.Username, &user.Password, &user.BankId, &user.BankName); err != nil {
			log.Println("Error scanning row:", err)
			return nil
		}
		return &user
	}

	return nil
}

func (wrapper *DBWrapper) IsCentralBankEmployee(username string) bool {
	query := `SELECT CASE
					WHEN EXISTS (
						SELECT 1
						FROM BankEmployee AS be
						LEFT JOIN Bank AS b ON b.Id = be.BankId
						WHERE be.username = @p1 AND b.BankTypeId = @p2
					)
					THEN 'true'
					ELSE 'false'
				END AS CentralBankEmployee;`

	rows, err := wrapper.db.Query(query, sql.Named("p1", username), sql.Named("p2", common.CentralBank))
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

func (wrapper *DBWrapper) GetBankId(bankName string) uint64 {
	query := `SELECT Id FROM [Bank] WHERE name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankName))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bankId uint64
	for rows.Next() {
		if err := rows.Scan(&bankId); err != nil {
			log.Fatal(err)
		}
	}

	return bankId
}

func (wrapper *DBWrapper) GetBankClientId(bankClientName string) uint64 {
	query := `SELECT Id FROM [BankClient] WHERE Name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankClientName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bankClientId uint64
	for rows.Next() {
		if err := rows.Scan(&bankClientId); err != nil {
			log.Fatal(err)
		}
	}

	return bankClientId
}

func (wrapper *DBWrapper) GetBank(bankId uint64) common.Bank {
	query := `SELECT b.Id, b.Name, b.CountryId
					From Bank b
					WHERE b.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bank common.Bank

	for rows.Next() {
		if err := rows.Scan(&bank.Id, &bank.Name, &bank.CountryId); err != nil {
			log.Fatal(err)
		}
	}

	return bank
}

func (wrapper *DBWrapper) GetBanks() []BankModel {
	query := `SELECT b.Id, b.Name, c.Name
					From Bank b
					JOIN Country as c ON c.Id = b.CountryId`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	banks := []BankModel{}
	for rows.Next() {
		var bank BankModel
		if err := rows.Scan(&bank.Id, &bank.Name, &bank.Country); err != nil {
			log.Fatal(err)
		}
		banks = append(banks, bank)
	}
	return banks
}

func (wrapper *DBWrapper) GetCountry(countryId uint) common.Country {
	query := `SELECT c.Id, c.Name, c.CountryCode
					From Country c
					WHERE c.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", countryId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var country common.Country

	for rows.Next() {
		if err := rows.Scan(&country.Id, &country.Name, &country.CountryCode); err != nil {
			log.Fatal(err)
		}
	}

	return country
}
