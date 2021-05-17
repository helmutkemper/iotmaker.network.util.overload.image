package main

import (
	overload "github.com/helmutkemper/iotmaker.network.util.overload"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var delayMin time.Duration
var delayMax time.Duration
var over *overload.NetworkOverload
var wg sync.WaitGroup

func main() {
	var err error

	var inAddress = os.Getenv("IN_ADDRESS")
	var outAddress = os.Getenv("OUT_ADDRESS")
	var minDelayString = os.Getenv("MIN_DELAY")
	var maxDelayString = os.Getenv("MAX_DELAY")

	if inAddress == "" {
		log.Print("in address is not set")
		return
	}

	if outAddress == "" {
		log.Print("out address is not set")
		return
	}

	if minDelayString == "" {
		log.Print("min delay is not set")
		return
	}

	if maxDelayString == "" {
		log.Print("max delay is not set")
		return
	}

	var minDelayInt64 int64
	var maxDelayInt64 int64

	minDelayInt64, err = strconv.ParseInt(minDelayString, 10, 64)
	if err != nil {
		log.Print("min delay parse error")
		return
	}

	maxDelayInt64, err = strconv.ParseInt(maxDelayString, 10, 64)
	if err != nil {
		log.Print("max delay parse error")
		return
	}

	// (English): Minimal delay between packages, 0.5 seconds
	// (Português): Atraso mínimo inserido entre os pacotes, 0.5 segundos
	delayMin = time.Millisecond * time.Duration(minDelayInt64)

	// (English): Maximal delay between packages, 5 seconds
	// (Português): Atraso máximo inserido entre os pacotes, 5 segundos
	delayMax = time.Millisecond * time.Duration(maxDelayInt64)

	over = &overload.NetworkOverload{
		ProtocolInterface: &overload.TCPConnection{},
	}

	// (English): Enables the TCP protocol and the input and output addresses
	// (Português): Habilita o protocolo TCP e os endereços de entrada e saída
	err = over.SetAddress(overload.KTypeNetworkTcp, inAddress, outAddress)
	if err != nil {
		log.Print("set address error")
		return
	}

	// (English): [optional] Points to the custom function for data processing
	// (Português): [opcional] Aponta a função personalizada para tratamento dos dados
	//over.ParserAppendTo(binaryDump)

	// (English): Determines the maximum and minimum times between packages
	// (Português): Determina os tempos máximo e mínimos entre os pacotes
	over.SetDelay(delayMin, delayMax)

	log.Printf("overloading...")

	wg.Add(1)
	go func() {
		go func() {
			http.HandleFunc("/maxDelay", maxDelay)
			http.HandleFunc("/minDelay", minDelay)
			err = http.ListenAndServe(":8080", nil)
			if err != nil {
				log.Print("command server error", err.Error())
				return
			}
		}()
	}()
	go func() {
		err = over.Listen()
		if err != nil {
			log.Print("listen error", err.Error())
			return
		}
	}()
	wg.Wait()
}

func maxDelay(w http.ResponseWriter, req *http.Request) {
	var err error
	var delayAsInt64 int64

	var delayAsString = req.RequestURI
	delayAsString = strings.Replace(delayAsString, "/maxDelay?", "", -1)
	delayAsInt64, err = strconv.ParseInt(delayAsString, 10, 64)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	delayMax = time.Duration(delayAsInt64) * time.Millisecond
	over.SetDelay(delayMin, delayMax)

	w.Write([]byte("ok"))
}

func minDelay(w http.ResponseWriter, req *http.Request) {
	var err error
	var delayAsInt64 int64

	var delayAsString = req.RequestURI
	delayAsString = strings.Replace(delayAsString, "/minDelay?", "", -1)
	delayAsInt64, err = strconv.ParseInt(delayAsString, 10, 64)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	delayMin = time.Duration(delayAsInt64) * time.Millisecond
	over.SetDelay(delayMin, delayMax)

	w.Write([]byte("ok"))
}
