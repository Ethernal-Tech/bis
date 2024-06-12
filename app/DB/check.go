package DB

import (
	"database/sql"
	"log"
	"strings"
)

func (wrapper *DBHandler) CheckCFM(receiverId string, countryId int) int64 {
	query := `SELECT b.GlobalIdentifier FROM BankClient as bc 
				JOIN (SELECT * FROM Bank Where CountryId = @p2) as b ON b.GlobalIdentifier = bc.BankId
				Where bc.GlobalIdentifier = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", receiverId),
		sql.Named("p2", countryId))

	if err != nil {
		log.Fatal(err)
	}

	var bankIds []string
	for rows.Next() {
		var bankId string
		if err := rows.Scan(&bankId); err != nil {
			log.Fatal(err)
		}
		bankIds = append(bankIds, bankId)
	}
	rows.Close()

	c := strings.Join(bankIds, `','`)

	query = `SELECT
			(SELECT ISNULL(SUM(Amount), 0)
			FROM [Transaction] t
            JOIN (SELECT TransactionId, StatusId FROM [TransactionHistory] WHERE StatusId = 7) as th on th.TransactionId = t.Id
			Where ReceiverId = @p1 and BeneficiaryBankId IN ('` + c + `') and TransactionTypeId IN (1))
			-
			((SELECT ISNULL(SUM(Amount), 0)
			FROM [Transaction] t
            JOIN (SELECT TransactionId, StatusId FROM [TransactionHistory] WHERE StatusId = 7) as th on th.TransactionId = t.Id
			Where SenderId = @p1 and OriginatorBankId IN ('` + c + `') and TransactionTypeId IN (2)))
			as difference`

	rows, err = wrapper.db.Query(query,
		sql.Named("p1", receiverId))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var amount int64
	for rows.Next() {
		if err := rows.Scan(&amount); err != nil {
			log.Fatal(err)
		}
	}

	return amount
}
