package manager

import (
	"bisgo/errlog"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SwiftManager struct {
	msgID uint64
}

type TokenResponse struct {
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             string `json:"expires_in"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func InitSwiftManager() *SwiftManager {
	return &SwiftManager{
		msgID: 2143657698,
	}
}

func (s *SwiftManager) GenerateAsserionToken() (string, error) {
	// Sample private key and certificate in PEM format
	privateKeyPEM := os.Getenv("PRIVATE_KEY_PEM")
	certPEM := os.Getenv("CERT_PEM")

	// Parse the private key
	privateKey, err := parsePEMPrivateKey(privateKeyPEM)
	if err != nil {
		errlog.Println(err)
		return "", err
	}

	// Parse the certificate
	cert, err := parsePEMCert(certPEM)
	if err != nil {
		errlog.Println(err)
		return "", err
	}

	// Create the JWT header
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": "RS256",
		"x5c": []string{base64.StdEncoding.EncodeToString(cert.Raw)},
	}

	// Create the JWT header
	claims := jwt.RegisteredClaims{
		Issuer:    "https://sandbox.swift.com/oauth2/v1/token",
		Audience:  jwt.ClaimStrings{"https://sandbox.swift.com/oauth2/v1/token"},
		Subject:   "CN=desktop, O=sandbox, O=swift",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        generateRandomString(12),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header = header
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		errlog.Println(err)
		return "", err
	}

	return signedToken, nil
}

// Helper function to generate a random string
func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Helper function to parse PEM encoded certificate
func parsePEMCert(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// Helper function to parse PEM encoded private key
func parsePEMPrivateKey(keyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *SwiftManager) GetJwtToken(assertion string) (string, error) {
	url := "https://sandbox.swift.com/oauth2/v1/token"
	method := "POST"

	payload := strings.NewReader("assertion=" + assertion + "&scope=swift.cbdc&grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		errlog.Println(err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set up Basic Authentication
	username := os.Getenv("CONSUMER_KEY")
	password := os.Getenv("CONSUMER_SECRET")
	req.SetBasicAuth(username, password)

	res, err := client.Do(req)
	if err != nil {
		errlog.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errlog.Println(err)
		return "", err
	}

	if res.StatusCode != 200 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			errlog.Println(err)
			return "", err
		}

		return "", errors.New(errorResponse.Error + ": " + errorResponse.ErrorDescription)
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		errlog.Println(err)
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func (s *SwiftManager) RevokeJwtToken(accessToken string) error {
	url := "https://sandbox.swift.com/oauth2/v1/revoke"
	method := "POST"

	expectedBody := "token=" + accessToken
	payload := strings.NewReader(expectedBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		errlog.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set up Basic Authentication
	username := os.Getenv("CONSUMER_KEY")
	password := os.Getenv("CONSUMER_SECRET")
	req.SetBasicAuth(username, password)

	res, err := client.Do(req)
	if err != nil {
		errlog.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errlog.Println(err)
		return err
	}

	if res.StatusCode != 200 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			errlog.Println(err)
			return err
		}

		return errors.New(errorResponse.Error + ": " + errorResponse.ErrorDescription)
	}

	if expectedBody == string(body) {
		return nil
	} else {
		return errors.New("token revoke failed, expected: " + expectedBody + " but got: " + string(body))
	}
}

// pacs.008 message struct
type PACS008Message struct {
	XMLNS             string            `json:"@xmlns"`
	FiToFICstmrCdtTrf FiToFICstmrCdtTrf `json:"fiToFICstmrCdtTrf"`
}

type FiToFICstmrCdtTrf struct {
	GrpHdr      GrpHdr      `json:"grpHdr"`
	CdtTrfTxInf CdtTrfTxInf `json:"cdtTrfTxInf"`
}

type GrpHdr struct {
	MsgId    string   `json:"msgId"`
	CreDtTm  string   `json:"creDtTm"`
	NbOfTxs  string   `json:"nbOfTxs"`
	SttlmInf SttlmInf `json:"sttlmInf"`
}

type SttlmInf struct {
	SttlmMtd string `json:"sttlmMtd"`
}

type CdtTrfTxInf struct {
	PmtId           PmtId           `json:"pmtId"`
	IntrBkSttlmAmt  IntrBkSttlmAmt  `json:"intrBkSttlmAmt"`
	IntrBkSttlmDt   string          `json:"intrBkSttlmDt"`
	InstdAmt        InstdAmt        `json:"instdAmt"`
	ChrgBr          string          `json:"chrgBr"`
	InstgAgt        InstgAgt        `json:"instgAgt"`
	InstdAgt        InstdAgt        `json:"instdAgt"`
	Dbtr            Dbtr            `json:"dbtr"`
	DbtrAcct        DbtrAcct        `json:"dbtrAcct"`
	DbtrAgt         DbtrAgt         `json:"dbtrAgt"`
	CdtrAgt         CdtrAgt         `json:"cdtrAgt"`
	Cdtr            Cdtr            `json:"cdtr"`
	CdtrAcct        CdtrAcct        `json:"cdtrAcct"`
	InstrForCdtrAgt InstrForCdtrAgt `json:"instrForCdtrAgt"`
}

type PmtId struct {
	InstrId    string `json:"instrId"`
	EndToEndId string `json:"endToEndId"`
	Uetr       string `json:"uetr"`
}

type IntrBkSttlmAmt struct {
	Ccy   string `json:"ccy"`
	Value string `json:"value"`
}

type InstdAmt struct {
	Ccy   string `json:"ccy"`
	Value string `json:"value"`
}

type InstgAgt struct {
	FinInstnId FinInstnId `json:"finInstnId"`
}

type InstdAgt struct {
	FinInstnId FinInstnId `json:"finInstnId"`
}

type FinInstnId struct {
	Bicfi string `json:"bicfi"`
}

type Dbtr struct {
	Nm      string  `json:"nm"`
	PstlAdr PstlAdr `json:"pstlAdr"`
	Id      Id      `json:"id"`
}

type DbtrAcct struct {
	Id Id `json:"id"`
}

type Id struct {
	Othr Othr `json:"othr"`
}

type Othr struct {
	Id string `json:"id"`
}

type PstlAdr struct {
	TwnNm string `json:"twnNm"`
	Ctry  string `json:"ctry"`
}

type DbtrAgt struct {
	FinInstnId FinInstnId `json:"finInstnId"`
}

type CdtrAgt struct {
	FinInstnId FinInstnId `json:"finInstnId"`
}

type Cdtr struct {
	Nm      string  `json:"nm"`
	PstlAdr PstlAdr `json:"pstlAdr"`
	Id      Id      `json:"id"`
}

type CdtrAcct struct {
	Id Id `json:"id"`
}

type InstrForCdtrAgt struct {
	InstrInf string `json:"instrInf"`
}

func (s *SwiftManager) GeneratePACS008Message(currency, amount, originatorName, beneficiaryName, complianceCheckID string) PACS008Message {
	msgID := fmt.Sprintf("%d", s.msgID+1)

	currentUTC := time.Now().UTC()

	dateAndTime := currentUTC.Format("2006-01-02T15:04:05+02:00")
	date := currentUTC.Format("2006-01-02")

	return PACS008Message{
		XMLNS: "urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08",
		FiToFICstmrCdtTrf: FiToFICstmrCdtTrf{
			GrpHdr: GrpHdr{
				MsgId:   msgID,
				CreDtTm: dateAndTime,
				NbOfTxs: "1",
				SttlmInf: SttlmInf{
					SttlmMtd: "INGA",
				},
			},
			CdtTrfTxInf: CdtTrfTxInf{
				PmtId: PmtId{
					InstrId:    msgID,
					EndToEndId: "pacs008EndToEndId-001",
					Uetr:       "8a562c67-ca16-48ba-b074-65581be6f001",
				},
				IntrBkSttlmAmt: IntrBkSttlmAmt{
					Ccy:   currency,
					Value: amount,
				},
				IntrBkSttlmDt: date,
				InstdAmt: InstdAmt{
					Ccy:   currency,
					Value: amount,
				},
				ChrgBr: "SHAR",
				InstgAgt: InstgAgt{
					FinInstnId: FinInstnId{
						Bicfi: "KRWBKRSEXXX",
					},
				},
				InstdAgt: InstdAgt{
					FinInstnId: FinInstnId{
						Bicfi: "AUDBAU2SXXX",
					},
				},
				Dbtr: Dbtr{
					Nm: originatorName,
					PstlAdr: PstlAdr{
						TwnNm: "Seoul",
						Ctry:  "KR",
					},
					Id: Id{
						Othr: Othr{
							Id: "34567890",
						},
					},
				},
				DbtrAcct: DbtrAcct{
					Id: Id{
						Othr: Othr{
							Id: "34567890",
						},
					},
				},
				DbtrAgt: DbtrAgt{
					FinInstnId: FinInstnId{
						Bicfi: "KRWBKRSEXXX",
					},
				},
				CdtrAgt: CdtrAgt{
					FinInstnId: FinInstnId{
						Bicfi: "AUDBAU2SXXX",
					},
				},
				Cdtr: Cdtr{
					Nm: beneficiaryName,
					PstlAdr: PstlAdr{
						TwnNm: "Sydney",
						Ctry:  "AU",
					},
					Id: Id{
						Othr: Othr{
							Id: "34567890",
						},
					},
				},
				CdtrAcct: CdtrAcct{
					Id: Id{
						Othr: Othr{
							Id: "34567890",
						},
					},
				},
				InstrForCdtrAgt: InstrForCdtrAgt{
					InstrInf: complianceCheckID,
				},
			},
		},
	}
}

func (s *SwiftManager) SendPACS008ToSwiftNetwork(message PACS008Message, accessToken, url string) error {
	// url := "https://sandbox.swift.com/cdbc-connector/customer-credit-transfers/abcd"
	method := "POST"

	expectedBody := "token=" + accessToken
	payload := strings.NewReader(expectedBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		errlog.Println(err)
		return err
	}

	// Set up Bearer Authentication
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		errlog.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errlog.Println(err)
		return err
	}

	fmt.Println(string(body))

	return nil
}
