package payment

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

var (
	msgId = uuid.New().String()
)

type Document struct {
	XMLName xml.Name `xml:"Document"`
	Xmlns   string   `xml:"xmlns,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Schema  string   `xml:"xsi:schemaLocation,attr"`
	GrpHdr  GrpHdr   `xml:"GrpHdr"`
	Cdtr    Cdtr     `xml:"Cdtr"`
	Dbtr    Dbtr     `xml:"Dbtr"`
	Amt     Amt      `xml:"Amt"`
}

type GrpHdr struct {
	MsgId   string `xml:"MsgId"`
	CreDtTm string `xml:"CreDtTm"`
}

type Cdtr struct {
	Nm       string   `xml:"Nm"`
	CdtrAcct CdtrAcct `xml:"CdtrAcct"`
}

type CdtrAcct struct {
	Id IBAN `xml:"Id"`
}

type IBAN struct {
	IBAN string `xml:"IBAN"`
}

type Dbtr struct {
	Nm       string   `xml:"Nm"`
	CdtrAcct DbtrAcct `xml:"CdtrAcct"`
}

type DbtrAcct struct {
	Id IBAN `xml:"Id"`
}

type Amt struct {
	Ccy   string  `xml:"Ccy,attr"`
	Value float64 `xml:",chardata"`
}

func generateXml(payment Payment) []byte {
	creDtTm := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	doc := Document{
		Xmlns:  "urn:iso:std:iso:20022:tech:xsd:pain.008.002.02",
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Schema: "urn:iso:std:iso:20022:tech:xsd:pain.008.002.02 pain.008.002.02.xsd",
		GrpHdr: GrpHdr{MsgId: msgId, CreDtTm: creDtTm},
		Cdtr:   Cdtr{Nm: payment.CreditorName, CdtrAcct: CdtrAcct{Id: IBAN{IBAN: payment.CreditorIban}}},
		Dbtr:   Dbtr{Nm: payment.DebtorName, CdtrAcct: DbtrAcct{Id: IBAN{IBAN: payment.DebtorIban}}},
		Amt:    Amt{Ccy: "EUR", Value: payment.Ammount},
	}

	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	xmlData = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), xmlData...)

	fmt.Println(string(xmlData))

	return xmlData
}

func WritePaymentToBank(payment Payment, filename string) error {
	data := generateXml(payment)
	err := os.WriteFile(filename+string(filepath.Separator)+msgId+".xml", data, 0644)
	return err
}
