package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"hwraid/web"

	"github.com/gorilla/mux"
)

var privKey rsa.PrivateKey

func main() {
	err := loadPrivKey()
	if err != nil {
		fmt.Printf("load priv key error: %v\n", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/license", web.JsonResponse(createLicense)).Methods("POST")
	l, _ := net.Listen("tcp", "192.168.2.197:8089")
	http.Serve(l, router)
}

func loadPrivKey() error {
	var file *os.File
	file, err := os.Open("privkey")
	if err != nil {
		return err
	}
	defer file.Close()

	var buf []byte = make([]byte, 4096)
	var n int

	n, err = file.Read(buf)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf[:n], &privKey)
	if err != nil {
		return err
	}

	return nil
}

func sign(sn string) ([]byte, error) {
	hash := md5.New()
	io.WriteString(hash, sn)
	hashed := hash.Sum(nil)

	var h crypto.Hash
	return rsa.SignPKCS1v15(rand.Reader, &privKey, h, hashed)
}

func createLicense(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	sn := r.FormValue("sn")

	fmt.Printf("sn: %v\n", sn)
	signature, err := sign(sn)
	if err != nil {
		return nil, err
	}
	return map[string]string{"signature": fmt.Sprintf("%x", signature)}, nil
}
