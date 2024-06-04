package tlshandler

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

func LoadCert(name string, directory string, cert *CertSetup) (err ftp_context.LogErr) {
	loc := "func LoadCert(name string,directory string, cert *CertSetup) (err ftp_context.LogErr) "
	cert.Name = name
	cert.Path = directory + "/" + name + ".json"

	pre_ := cert.Open()
	if pre_.Err != nil {
		return ftp_context.NewLogItem(loc, true).Set("after", "cert.Open()").SetMessage("unable to open provided file").AppendParentError(pre_.Err)
	}

	d, err_ := pre_.ReadAll()
	if err_ != nil {
		return ftp_context.NewLogItem(loc, true).Set("after", "pre_.ReadAll()").SetMessage("unable to read file").AppendParentError(err_)
	}
	err_ = json.Unmarshal(d, cert)
	if err_ != nil {
		return ftp_context.NewLogItem(loc, true).Set("after", "json.Unmarshal(d,cert)").SetMessage("unable to unmarshall read data to json").AppendParentError(err_)
	}
	return
}

func GinHandler(router *gin.Engine, server_cert *CertSetup, httpAddr string) (err ftp_context.LogErr) {

	loc := `handler_2(router *gin.Engine) (err ftp_context.LogErr)`
	log.Println("starting", loc)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{server_cert.TlsCert.Certificate},
	}
	server := http.Server{Addr: httpAddr, Handler: router, TLSConfig: tlsConfig}

	if pre_ := server.ListenAndServeTLS("", ""); pre_ != nil {
		return ftp_context.NewLogItem(loc, true).Set("after", "server.ListenAndServeTLS").AppendParentError(pre_)
	}

	return
}

func TLSClient(caPEMs ...[]byte) (tc *http.Client) {

	certpool := x509.NewCertPool()
	for _, caPEM := range caPEMs {
		certpool.AppendCertsFromPEM(caPEM)
	}
	tlc := &tls.Config{
		RootCAs: certpool,
	}

	tp := &http.Transport{
		TLSClientConfig: tlc,
	}
	tc = &http.Client{
		Transport: tp,
	}
	return
}
