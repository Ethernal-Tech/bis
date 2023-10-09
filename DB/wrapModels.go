package DB

type TransactionModel struct {
	Id              uint64
	OriginatorBank  string
	BeneficiaryBank string
	Sender          uint64
	Receiver        uint64
	Curency         string
	Amount          int
	Type            string
	Status          string
}
