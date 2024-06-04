package tlshandler

import (
	"crypto/x509/pkix"
	"net"
)

type CertData struct {
	Organization  string         `json:"organisation"`
	Country       string         `json:"country"`
	Province      string         `json:"province"`
	Locality      string         `json:"locality"`
	StreetAddress string         `json:"street_address"`
	PostalCode    string         `json:"postal_code"`
	NotAfter      NotAfterStruct `json:"add_date"`
	IPAddrresses  []net.IP       `json:"ip_addresses"`
}

type NotAfterStruct struct {
	Years  int `json:"years"`
	Months int `json:"months"`
	Days   int `json:"days"`
}

func GetSubject(cd CertData) pkix.Name {

	return pkix.Name{
		Organization:  []string{cd.Organization},
		Country:       []string{cd.Country},
		Province:      []string{cd.Province},
		Locality:      []string{cd.Locality},
		StreetAddress: []string{cd.StreetAddress},
		PostalCode:    []string{cd.PostalCode},
	}
}
