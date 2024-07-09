package DB

import (
	"bisgo/app/models"
	"bisgo/common"
	"database/sql"
	"log"
	"time"
)

func (wrapper *DBHandler) InsertTransaction(t models.NewTransaction) string {
	query := `INSERT INTO [dbo].[Transaction] (Id, OriginatorBankId, BeneficiaryBankId, SenderId, ReceiverId, Currency, Amount, TransactionTypeId, LoanId)
          VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`

	var insertedID string

	if t.Id != "" {
		insertedID = t.Id
	} else {
		var err error
		insertedID, err = common.GenerateHash(t)
		if err != nil {
			log.Fatal(err)
		}
	}

	row := wrapper.db.QueryRow(query,
		sql.Named("p1", insertedID),
		sql.Named("p2", t.OriginatorBankId),
		sql.Named("p3", t.BeneficiaryBankId),
		sql.Named("p4", int(t.SenderId)),
		sql.Named("p5", int(t.ReceiverId)),
		sql.Named("p6", t.Currency),
		sql.Named("p7", t.Amount),
		sql.Named("p8", t.TransactionTypeId),
		sql.Named("p9", t.LoanId))

	if row.Err() != nil {
		log.Fatal(row.Err())
	}

	//polices := wrapper.GetPolices(t.BeneficiaryBankId, t.TransactionTypeId)

	// for _, policy := range polices {
	// 	query = `INSERT INTO [dbo].[TransactionPolicy] VALUES (@p1, @p2, @p3)`
	// 	_, err := wrapper.db.Exec(query,
	// 		sql.Named("p1", insertedID),
	// 		sql.Named("p2", policy.Id),
	// 		sql.Named("p3", 0))

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	return insertedID
}

func (wrapper *DBHandler) InsertTransactionProof(transactionId string, value string) {
	query := `INSERT INTO [dbo].[TransactionProof] VALUES (@p1, @p2)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", value))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBHandler) UpdateTransactionState(transactionId string, state int) {
	query := `INSERT INTO [dbo].[TransactionHistory] VALUES (@p1, @p2, @p3)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", state),
		sql.Named("p3", time.Now()))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBHandler) GetTransactionTypeId(transactionType string) int {
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

func (wrapper *DBHandler) GetTransactionTypes() []models.NewTransactionType {
	query := `SELECT Id, Code, Name From TransactionType`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	types := []models.NewTransactionType{}
	for rows.Next() {
		var tType models.NewTransactionType
		if err := rows.Scan(&tType.Id, &tType.Code, &tType.Name); err != nil {
			log.Fatal(err)
		}
		types = append(types, tType)
	}
	return types
}

func (wrapper *DBHandler) GetCommercialBankTransactions(bankId string, searchModel models.SearchModel) []models.TransactionModel {
	query := `WITH CreatedTransactions AS (
					SELECT	th.TransactionId,
							th.Date AS CreatedDate
					FROM [dbo].[TransactionHistory] th
					WHERE th.StatusId = 1`

	if searchModel.From != "" {
		query += " AND th.Date >= '" + searchModel.From + "'"
	}
	if searchModel.To != "" {
		query += " AND th.Date <= '" + searchModel.To + "'"
	}

	query += `), LatestStatus AS (
					SELECT  th.TransactionId,
							th.StatusId,
							th.Date,
							ROW_NUMBER() OVER (PARTITION BY th.TransactionId ORDER BY th.Date DESC) AS rn
					FROM [dbo].[TransactionHistory] th)

				SELECT   t.Id
						,ob.Name
						,bb.Name
						,bcs.GlobalIdentifier
						,bcr.GlobalIdentifier
						,bcs.Name
						,bcr.Name
						,t.Currency
						,t.Amount
						,s.Name
					FROM [dbo].[Transaction] t
					JOIN CreatedTransactions ct ON t.Id = ct.TransactionId
					JOIN LatestStatus ls ON t.Id = ls.TransactionId AND ls.rn = 1
					JOIN [dbo].[Status] s ON ls.StatusId = s.Id
					JOIN Bank as ob ON ob.GlobalIdentifier = t.OriginatorBankId
					JOIN Bank as bb ON bb.GlobalIdentifier = t.BeneficiaryBankId
					JOIN BankClient as bcs ON bcs.Id = t.SenderId
					JOIN BankClient as bcr ON bcr.Id = t.ReceiverId
				WHERE (t.OriginatorBankId = @p1 OR t.BeneficiaryBankId = @p1) and (ob.Name like @p2 OR bb.Name like @p2 OR bcs.Name like @p2 OR bcr.Name like @p2)`

	if searchModel.StatusId != "" {
		query += ` AND ls.StatusId = ` + searchModel.StatusId
	}

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId),
		sql.Named("p2", "%"+searchModel.Value+"%"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	transactions := []models.TransactionModel{}
	for rows.Next() {
		var trnx models.TransactionModel
		if err := rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdentifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.Status); err != nil {
			log.Println("Error scanning row:", err)
			return nil
		}
		trnx = *convertTxStatusDBtoPR(&trnx)
		transactions = append(transactions, trnx)
	}
	return transactions
}

func (wrapper *DBHandler) GetCentralBankTransactions(bankId string, searchModel models.SearchModel) ([]models.TransactionModel, string) {
	query := `SELECT JurisdictionId FROM Bank
			Where GlobalIdentifier = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var centralBankJurisdictionId string
	for rows.Next() {
		if err = rows.Scan(&centralBankJurisdictionId); err != nil {
			log.Println("Error scanning row:", err)
			return []models.TransactionModel{}, ""
		}
	}

	query = `WITH CreatedTransactions AS (
		SELECT	th.TransactionId,
				th.Date AS CreatedDate
		FROM [dbo].[TransactionHistory] th
		WHERE th.StatusId = 1`

	if searchModel.From != "" {
		query += " AND th.Date >= '" + searchModel.From + "'"
	}
	if searchModel.To != "" {
		query += " AND th.Date <= '" + searchModel.To + "'"
	}

	query += `), LatestStatus AS (
			SELECT  th.TransactionId,
					th.StatusId,
					th.Date,
					ROW_NUMBER() OVER (PARTITION BY th.TransactionId ORDER BY th.Date DESC) AS rn
			FROM [dbo].[TransactionHistory] th)

		SELECT   t.Id
				,ob.Name
				,ob.JurisdictionId
				,bb.Name
				,bcs.GlobalIdentifier
				,bcr.GlobalIdentifier
				,bcs.Name
				,bcr.Name
				,t.Currency
				,t.Amount
				,s.Name
			FROM [dbo].[Transaction] t
			JOIN CreatedTransactions ct ON t.Id = ct.TransactionId
			JOIN LatestStatus ls ON t.Id = ls.TransactionId AND ls.rn = 1
			JOIN [dbo].[Status] s ON ls.StatusId = s.Id
			JOIN Bank as ob ON ob.GlobalIdentifier = t.OriginatorBankId
			JOIN Bank as bb ON bb.GlobalIdentifier = t.BeneficiaryBankId
			JOIN BankClient as bcs ON bcs.Id = t.SenderId
			JOIN BankClient as bcr ON bcr.Id = t.ReceiverId
		WHERE (ob.JurisdictionId = @p1 OR bb.JurisdictionId = @p1) and (ob.Name like @p2 OR bb.Name like @p2 OR bcs.Name like @p2 OR bcr.Name like @p2)`

	if searchModel.StatusId != "" {
		query += ` AND ls.StatusId = ` + searchModel.StatusId
	}

	rows, err = wrapper.db.Query(query,
		sql.Named("p1", centralBankJurisdictionId),
		sql.Named("p2", "%"+searchModel.Value+"%"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	transactions := []models.TransactionModel{}
	for rows.Next() {
		var trnx models.TransactionModel
		if err = rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.OriginatorBankCountryId, &trnx.BeneficiaryBank,
			&trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdentifier, &trnx.SenderName, &trnx.ReceiverName,
			&trnx.Currency, &trnx.Amount, &trnx.Status); err != nil {
			log.Println("Error scanning row:", err)
			return []models.TransactionModel{}, ""
		}

		trnx = *convertTxStatusDBtoPR(&trnx)
		transactions = append(transactions, trnx)
	}
	return transactions, centralBankJurisdictionId
}

func (wrapper *DBHandler) GetTransactionHistory(transactionId string) models.TransactionModel {
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
				JOIN Bank as ob ON ob.GlobalIdentifier = t.OriginatorBankId
				JOIN Bank as bb ON bb.GlobalIdentifier = t.BeneficiaryBankId
				JOIN BankClient as bcs ON bcs.Id = t.SenderId
				JOIN BankClient as bcr ON bcr.Id = t.ReceiverId
				JOIN [TransactionType] as ty ON ty.Id = t.TransactionTypeId
				WHERE t.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}

	var trnx models.TransactionModel
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
		var statusHistory models.StatusHistoryModel
		if err := rows.Scan(&statusHistory.Name, &statusHistory.Date); err != nil {
			log.Fatal(err)
		}
		statusHistory = *convertHistoryStatusDBtoPR(&statusHistory)
		statusHistory.DateString = statusHistory.Date.Format("2006-01-02 15:04:05")
		trnx.StatusHistory = append(trnx.StatusHistory, statusHistory)
	}
	trnx.Status = trnx.StatusHistory[len(trnx.StatusHistory)-1].Name

	query = `SELECT pt.Name
				FROM TransactionPolicy tp
				JOIN Policy as p ON tp.PolicyId = p.Id
				JOIN PolicyType pt ON p.PolicyTypeId = pt.Id
				Where p.TransactionTypeId = (SELECT t.TransactionTypeId FROM [Transaction] as t
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

func (wrapper *DBHandler) GetComplianceCheckByID(checkID string) models.NewTransaction {
	query := `SELECT t.Id, t.OriginatorBankId, t.BeneficiaryBankId, t.SenderId, t.ReceiverId, t.Currency, t.Amount, t.TransactionTypeId, t.LoanId
              FROM [Transaction] t
              WHERE t.Id = @p1`

	row := wrapper.db.QueryRow(query, sql.Named("p1", checkID))

	var transaction models.NewTransaction
	if err := row.Scan(&transaction.Id, &transaction.OriginatorBankId, &transaction.BeneficiaryBankId, &transaction.SenderId, &transaction.ReceiverId, &transaction.Currency, &transaction.Amount, &transaction.TransactionTypeId, &transaction.LoanId); err != nil {
		if err == sql.ErrNoRows {
			return models.NewTransaction{}
		}
		log.Fatal(err)
	}

	return transaction
}
