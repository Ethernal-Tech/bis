package db

import (
	"bisgo/common"
	"database/sql"
	"log"
	"time"
)

func (wrapper *DBWrapper) InsertTransaction(t common.Transaction) uint64 {
	query := `INSERT INTO [dbo].[Transaction] OUTPUT INSERTED.Id VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`

	var insertedID uint64

	err := wrapper.db.QueryRow(query,
		sql.Named("p1", t.OriginatorBank),
		sql.Named("p2", t.BeneficiaryBank),
		sql.Named("p3", t.Sender),
		sql.Named("p4", t.Receiver),
		sql.Named("p5", t.Currency),
		sql.Named("p6", t.Amount),
		sql.Named("p7", t.TypeId),
		sql.Named("p8", t.LoanId)).Scan(&insertedID)

	if err != nil {
		log.Fatal(err)
	}

	polices := wrapper.GetPolices(t.BeneficiaryBank, t.TypeId)

	for _, policy := range polices {
		query = `INSERT INTO [dbo].[TransactionPolicyStatus] VALUES (@p1, @p2, @p3)`
		_, err := wrapper.db.Exec(query,
			sql.Named("p1", insertedID),
			sql.Named("p2", policy.Id),
			sql.Named("p3", 0))

		if err != nil {
			log.Fatal(err)
		}
	}

	return insertedID
}

func (wrapper *DBWrapper) InsertTransactionProof(transactionId uint64, value string) {
	query := `INSERT INTO [dbo].[TransactionProof] VALUES (@p1, @p2)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", value))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBWrapper) UpdateTransactionState(transactionId uint64, state int) {
	query := `INSERT INTO [dbo].[TransactionHistory] VALUES (@p1, @p2, @p3)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", state),
		sql.Named("p3", time.Now()))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBWrapper) GetTransactionTypeId(transactionType string) int {
	query := `SELECT Id FROM [TransactionType] WHERE Code = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionType))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var transactionTypeId int
	for rows.Next() {
		if err := rows.Scan(&transactionTypeId); err != nil {
			log.Fatal(err)
		}
	}

	return transactionTypeId
}

func (wrapper *DBWrapper) GetTransactionTypes() []common.TransactionType {
	query := `SELECT Id, Code, Name From TransactionType`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	types := []common.TransactionType{}
	for rows.Next() {
		var tType common.TransactionType
		if err := rows.Scan(&tType.Id, &tType.Code, &tType.Name); err != nil {
			log.Fatal(err)
		}
		types = append(types, tType)
	}
	return types
}

func (wrapper *DBWrapper) GetTransactionsForAddress(address uint64, searchValue string) []TransactionModel {
	query := `SELECT t.Id
					,ob.Name
					,bb.Name
					,bcs.GlobalIdentifier
					,bcr.GlobalIdentifier
					,bcs.Name
					,bcr.Name
					,t.Currency
					,t.Amount
					,s.Name
				FROM [Transaction] as t
				LEFT JOIN (SELECT MAX(StatusId) AS StatusId, Transactionid FROM TransactionHistory GROUP BY Transactionid) as th ON th.Transactionid = t.Id 
				LEFT JOIN [Status] as s ON s.Id = th.StatusId
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				WHERE t.OriginatorBank = @p1 OR t.BeneficiaryBank = @p2`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", address),
		sql.Named("p2", address))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	transactions := []TransactionModel{}
	for rows.Next() {
		var trnx TransactionModel
		if err := rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdentifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.Status); err != nil {
			log.Println("Error scanning row:", err)
			return nil
		}
		trnx = *convertTxStatusDBtoPR(&trnx)
		transactions = append(transactions, trnx)
	}
	return transactions
}

func (wrapper *DBWrapper) GetTransactionsForCentralbank(bankId uint64, searchValue string) ([]TransactionModel, int) {
	query := `SELECT CountryId FROM Bank
			Where Id = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var centralBankCountryId int
	for rows.Next() {
		if err = rows.Scan(&centralBankCountryId); err != nil {
			log.Println("Error scanning row:", err)
			return []TransactionModel{}, 0
		}
	}

	query = `SELECT t.Id
					,ob.Name
					,ob.CountryId
					,bb.Name
					,bcs.GlobalIdentifier
					,bcr.GlobalIdentifier
					,bcs.Name
					,bcr.Name
					,t.Currency
					,t.Amount
					,s.Name
				FROM [Transaction] as t
				LEFT JOIN (SELECT MAX(StatusId) AS StatusId, Transactionid FROM TransactionHistory GROUP BY Transactionid) as th ON th.Transactionid = t.Id 
				LEFT JOIN [Status] as s ON s.Id = th.StatusId
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				WHERE (ob.CountryId = @p1 OR bb.CountryId = @p2) and (ob.Name like @p3 OR bb.Name like @p3 OR bcs.Name like @p3 OR bcr.Name like @p3)`

	rows, err = wrapper.db.Query(query,
		sql.Named("p1", centralBankCountryId),
		sql.Named("p2", centralBankCountryId),
		sql.Named("p3", "%"+searchValue+"%"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	transactions := []TransactionModel{}
	for rows.Next() {
		var trnx TransactionModel
		if err = rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.OriginatorBankCountryId, &trnx.BeneficiaryBank,
			&trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdentifier, &trnx.SenderName, &trnx.ReceiverName,
			&trnx.Currency, &trnx.Amount, &trnx.Status); err != nil {
			log.Println("Error scanning row:", err)
			return []TransactionModel{}, 0
		}

		trnx = *convertTxStatusDBtoPR(&trnx)
		transactions = append(transactions, trnx)
	}
	return transactions, centralBankCountryId
}

func (wrapper *DBWrapper) GetTransactionHistory(transactionId uint64) TransactionModel {
	query := `SELECT t.Id
					,ob.Name
					,bb.Name
					,bcs.GlobalIdentifier
					,bcr.GlobalIdentifier
					,bcs.Name
					,bcr.Name
					,t.Currency
					,t.Amount
					,t.LoanId
					,ty.Code
					,ty.Name
					,ty.Id
				FROM [Transaction] as t
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				JOIN [TransactionType] as ty ON ty.Id = t.TransactionTypeId
				WHERE t.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}

	var trnx TransactionModel
	if rows.Next() {
		if err := rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdentifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.LoanId, &trnx.TypeCode, &trnx.Type, &trnx.TypeId); err != nil {
			log.Fatal(err)
		}
	}

	rows.Close()

	query = `SELECT s.Name
					,th.Date
				FROM TransactionHistory th
				JOIN Status as s ON s.Id = th.StatusId
				Where Transactionid = @p1 Order by StatusId`

	rows, err = wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var statusHistory StatusHistoryModel
		if err := rows.Scan(&statusHistory.Name, &statusHistory.Date); err != nil {
			log.Fatal(err)
		}
		statusHistory = *convertHistoryStatusDBtoPR(&statusHistory)
		statusHistory.DateString = statusHistory.Date.Format("2006-01-02 15:04:05")
		trnx.StatusHistory = append(trnx.StatusHistory, statusHistory)
	}
	trnx.Status = trnx.StatusHistory[len(trnx.StatusHistory)-1].Name

	query = `SELECT p.Name
				FROM TransactionTypePolicy ttp
				JOIN Policy as p ON ttp.PolicyId = p.Id
				Where ttp.TransactionTypeId = (SELECT t.TransactionTypeId FROM [Transaction] as t
												Where t.Id = @p1)
					and ttp.CountryId = (SELECT b.CountryId FROM [Transaction] as t
										JOIN Bank as b ON b.Id = t.BeneficiaryBank
										Where t.Id = @p1)`

	rows, err = wrapper.db.Query(query, sql.Named("p1", transactionId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var policyName string
		if err := rows.Scan(&policyName); err != nil {
			log.Fatal(err)
		}
		trnx.Policies = append(trnx.Policies, policyName)
	}

	return trnx
}
