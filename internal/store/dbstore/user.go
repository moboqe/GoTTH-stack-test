package dbstore

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"goth/internal/config"
	"goth/internal/hash"
	"goth/internal/store"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserStore struct {
	db           *gorm.DB
	passwordhash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db:           params.DB,
		passwordhash: params.PasswordHash,
	}
}

func (s *UserStore) CreateUser(email string, password string) (string, error) {
	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return "", err
	}
	res1, err := reg("bnetaccount create", email, password)
	if err != nil {

	}
	res2, err := reg("account set gmlevel", email, password)
	if err != nil {
	}

	res := strings.Join([]string{res1, res2}, "\r\n")

	return res, s.db.Create(&store.User{
		Email:    email,
		Password: hashedPassword,
	}).Error
}

func (s *UserStore) GetUser(email string) (*store.User, error) {

	var user store.User
	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

type Envelope struct {
	XMLName    xml.Name `xml:"Envelope"`
	XMLNSsenv  string   `xml:"SOAP-ENV,attr"`
	XMLNSns1   string   `xml:"ns1,attr"`
	XMLNSxsd   string   `xml:"xsd,attr"`
	XMLNSxsi   string   `xml:"xsi,attr"`
	XMLNSxsenc string   `xml:"SOAP-ENC,attr"`
	Body       Body     `xml:"Body"`
}

type Body struct {
	Response Response `xml:"executeCommandResponse,omitempty"`
	Fault    Fault    `xml:"Fault,omitempty"`
}
type Response struct {
	Result string `xml:"result,omitempty"`
}
type Fault struct {
	FaultCode   string `xml:"faultcode,omitempty"`
	FaultString string `xml:"faultstring,omitempty"`
	Detail      string `xml:"detail,omitempty"`
}

func reg(initialCommand string, email string, password string) (result string, err error) {
	trinityCfg := config.MustLoadTrinity()
	mySQLCfg := mysql.Config{User: trinityCfg.MySQL.User,
		Passwd: trinityCfg.MySQL.Passwd,
		Net:    trinityCfg.MySQL.Net,
		Addr:   trinityCfg.MySQL.Addr,
		DBName: trinityCfg.MySQL.DBName}

	db, err := sql.Open("mysql", mySQLCfg.FormatDSN())
	if err != nil {
		log.Fatalf(err.Error())
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("DB connected!!!")

	var command string
	if initialCommand == "bnetaccount create" {
		command = strings.Join([]string{initialCommand, email, password}, " ")
	} else {
		var username string
		row := db.QueryRow("SELECT username FROM account WHERE email = ?", email)
		if err := row.Scan(&username); err != nil {
			if err == sql.ErrNoRows {
				return result, fmt.Errorf("hashByEmail %s: no such hash", email)
			}
			return result, fmt.Errorf("hashByEmail %s: %v", email, err)
		}
		command = strings.Join([]string{initialCommand, username, "3", "-1"}, " ")
	}

	rootUsername := trinityCfg.WorldServer.RootUsername
	rootPassword := trinityCfg.WorldServer.RootPassword
	url := trinityCfg.WorldServer.WorldServerURL

	method := "POST"

	payload := strings.NewReader(`<?xml version="1.0" encoding="utf-8"?>` +
		"" +
		`<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="urn:TC" xmlns:xsd="http://www.w3.org/1999/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">` +
		"" +
		`    <SOAP-ENV:Body>` +
		"" + `
			<ns1:executeCommand>` +
		"" + `<command>` + command + `</command>` +
		"" + ` </ns1:executeCommand>` +
		"" + `
		</SOAP-ENV:Body>` + "" + `
	</SOAP-ENV:Envelope>`)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
	defer cancel()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	stringToEncode := strings.Join([]string{rootUsername, rootPassword}, ":")
	origBytes := []byte(stringToEncode)
	encodedText := base64.StdEncoding.EncodeToString(origBytes)
	auth := strings.Join([]string{"Basic", encodedText}, " ")

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Authorization", auth)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	soapResponse := Envelope{}
	if err = xml.Unmarshal(body, &soapResponse); err != nil {
		fmt.Println(err)
		return
	}

	empJSON, err := json.MarshalIndent(soapResponse, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("MarshalIndent function output %s\n", string(empJSON))

	return soapResponse.Body.Response.Result, nil
}
