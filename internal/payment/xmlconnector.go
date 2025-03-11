package payment

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type GrpHdr struct {
	MsgId   string `xml:"MsgId"`
	CreDtTm string `xml:"CreDtTm"`
}

type IBAN struct {
	IBAN string `xml:"IBAN"`
}

type Cdtr struct {
	Nm       string   `xml:"Nm"`
	CdtrAcct CdtrAcct `xml:"CdtrAcct"`
}

type CdtrAcct struct {
	Id IBAN `xml:"Id"`
}

type DbtrAcct struct {
	Id IBAN `xml:"Id"`
}

type Dbtr struct {
	Nm       string   `xml:"Nm"`
	CdtrAcct DbtrAcct `xml:"CdtrAcct"`
}

type Amt struct {
	Ccy   string `xml:"Ccy,attr"`
	Value int64  `xml:",chardata"`
}

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

func generateXml(payment Payment) ([]byte, error) {
	creDtTm := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	doc := Document{
		Xmlns:  "urn:iso:std:iso:20022:tech:xsd:pain.008.002.02",
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Schema: "urn:iso:std:iso:20022:tech:xsd:pain.008.002.02 pain.008.002.02.xsd",
		GrpHdr: GrpHdr{MsgId: payment.IdempotencyUniqueKey, CreDtTm: creDtTm},
		Cdtr:   Cdtr{Nm: payment.CreditorName, CdtrAcct: CdtrAcct{Id: IBAN{IBAN: payment.CreditorIban}}},
		Dbtr:   Dbtr{Nm: payment.DebtorName, CdtrAcct: DbtrAcct{Id: IBAN{IBAN: payment.DebtorIban}}},
		Amt:    Amt{Ccy: "EUR", Value: payment.Ammount},
	}

	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Sprintln("Error compiling XML document:", err)
		return nil, err
	}

	xmlData = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), xmlData...)

	return xmlData, nil
}

func WritePaymentToBank(payment Payment, filename string) error {
	data, err := generateXml(payment)
	if err != nil {
		return err
	}

	msgId := uuid.New().String()
	filename = filename + string(filepath.Separator) + msgId + ".xml"

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write payment XML: %w", err)
	}

	return nil
}
