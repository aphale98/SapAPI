package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	baseURL = "https://sandbox.api.sap.com/s4hanacloud/sap/opu/odata4/sap/api_cabillingrequest/srvd_a2x/sap/cabillingrequest/0001/CABillgRequest(%27100000000001%27)"
	apiKey  = "rzAifWlYUCoWRwYjowGEu32wdIAqu2Nn"
	// certFile = "C:/Windows/System32/cert.pem"
	// keyFile  = "C:/Windows/System32/key.pem"
)

type BillingRequest struct {
	CABillgReqDocument string `json:"CABillgReqDocument"`
	// CABillgReqCategory            string
	// CABillgReqType                string
	// CABillgReqReason              string
	// CAApplicationArea             string
	// ContractAccount               int
	BusinessPartner string `json:"BusinessPartner"`
	// CAInvoicingDocument           string
	// CABillgReqDescription         string
	// CABillgReqStatus              string
	CABillgReqTotalAmount         float64 `json:"CABillgReqTotalAmount"`
	CABillgReqTotalAmountCurrency string  `json:"CABillgReqTotalAmountCurrency"`
	// CABillgReqCreationUser        string
	// CABillgReqCreationDate        string
	// CABillgReqCreationTime        string
	// CABillgReqChangeUser          string
	// CABillgReqChangeDate          string
	// CABillgReqChangeTime          string
	// CAClrfctnExist                string
	// CABillgReqReference           string
	// CABillgReqNumberOfItems       string
	// LogicalSystem                 string
	// CABllbleItmListId             string
	// CADeletionDate                string
	// CABillgReqCompletionDate      string
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/bills", getBillingRequest)

	// srv := &http.Server{
	// 	Addr: ":8080",
	// 	TLSConfig: &tls.Config{
	// 		Certificates: getCertificates(),
	// 	},
	// }

	// go func() {
	// 	if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()

	router.Run(":8080")
	// router.RunTLS(":8080", certFile, keyFile)
}

// func getCertificates() []tls.Certificate {
// 	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
// 	if err != nil {
// 		log.Fatalf("server: loadkeys: %s", err)
// 	}

// 	return []tls.Certificate{cert}
// }

func getBillingRequest(c *gin.Context) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("APIKey", apiKey)
	req.Header.Set("DataServiceVersion", "2.0")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting bills: %s", resp.Status)})
		return
	}

	bill := BillingRequest{}
	err = json.Unmarshal(body, &bill)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	bills, err := json.Marshal(bill)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	strArr := strings.Split(string(bills), " ")

	c.JSON(http.StatusOK, strArr)
}
