package main

import (
	overload "github.com/helmutkemper/iotmaker.network.util.overload"
	"log"
	"os"
	"strconv"
	"time"
)

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
	var delayMin = time.Millisecond * time.Duration(minDelayInt64)

	// (English): Maximal delay between packages, 5 seconds
	// (Português): Atraso máximo inserido entre os pacotes, 5 segundos
	var delayMax = time.Millisecond * time.Duration(maxDelayInt64)

	var over = &overload.NetworkOverload{
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

	err = over.Listen()
	if err != nil {
		log.Print("listen error", err.Error())
		return
	}
}
