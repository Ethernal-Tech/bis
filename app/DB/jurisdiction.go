package DB

import (
	"bisgo/app/models"
	"database/sql"
	"log"
)

func (wrapper *DBHandler) GetJurisdiction(jurisdictionId string) models.Jurisdiction {
	query := `SELECT j.Id, j.Name
					From Jurisdiction j
					WHERE j.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", jurisdictionId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var jurisdiction models.Jurisdiction

	for rows.Next() {
		if err := rows.Scan(&jurisdiction.Id, &jurisdiction.Name); err != nil {
			log.Fatal(err)
		}
	}

	return jurisdiction
}

func (wrapper *DBHandler) GetJurisdictionOfBank(bankId string) models.Jurisdiction {
	query := `SELECT j.Id, j.Name FROM Bank as b
					LEFT JOIN Jurisdiction as j on b.JurisdictionId = j.Id
					Where GlobalIdentifier = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var jurisdiction models.Jurisdiction

	for rows.Next() {
		if err := rows.Scan(&jurisdiction.Id, &jurisdiction.Name); err != nil {
			log.Fatal(err)
		}
	}
	return jurisdiction
}
