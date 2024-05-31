package DB

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

func (wrapper *DBHandler) CheckCFM(receiverId uint64, countryId int) int64 {
	query := `SELECT GlobalIdentifier FROM BankClient Where Id = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", receiverId))

	if err != nil {
		log.Fatal(err)
	}

	var globalIdentifier string
	for rows.Next() {
		if err := rows.Scan(&globalIdentifier); err != nil {
			log.Fatal(err)
		}
	}

	rows.Close()
	query = `SELECT b.Id FROM BankClient as bc Join (SELECT * FROM Bank Where CountryId = @p2) as b ON b.Id = bc.BankId	Where bc.GlobalIdentifier = @p1`

	rows, err = wrapper.db.Query(query,
		sql.Named("p1", globalIdentifier),
		sql.Named("p2", countryId))

	if err != nil {
		log.Fatal(err)
	}

	var bankIds []string
	for rows.Next() {
		var bankId uint64
		if err := rows.Scan(&bankId); err != nil {
			log.Fatal(err)
		}
		bankIds = append(bankIds, strconv.Itoa(int(bankId)))
	}
	rows.Close()

	c := strings.Join(bankIds, ",")

	query = `SELECT
			(SELECT ISNULL(SUM(Amount), 0)
			FROM [Transaction] t
            JOIN (SELECT TransactionId, StatusId FROM [TransactionHistory] WHERE StatusId = 7) as th on th.TransactionId = t.Id
			Where Receiver = @p1 and BeneficiaryBank IN (` + c + `) and TransactionTypeId IN (1))
			-
			((SELECT ISNULL(SUM(Amount), 0)
			FROM [Transaction] t
            JOIN (SELECT TransactionId, StatusId FROM [TransactionHistory] WHERE StatusId = 7) as th on th.TransactionId = t.Id
			Where Sender = @p1 and OriginatorBank IN (` + c + `) and TransactionTypeId IN (2)))
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
