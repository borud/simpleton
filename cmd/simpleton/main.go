package main

import (
	"log"
	"net"
	"os"

	"github.com/borud/simpleton/pkg/store"
	"github.com/jessevdk/go-flags"
)

// Options contains the command line options
//
type Options struct {
	// UDP listener
	UDPListenAddress string `short:"u" long:"udp-listener" description:"Listen address for UDP listener" default:":7000" value-name:"<[host]:port>"`
	UDPBufferSize    int    `short:"b" long:"udp-buffer-size" description:"Size of UDP read buffer" default:"1024" value-name:"<num bytes>"`

	// Database options
	DBFilename string `short:"d" long:"db" description:"Data storage file" default:"simpleton.db" value-name:"<file>"`

	// Verbose
	Verbose bool `short:"v" long:"verbose" description:"Turn on verbose logging"`
}

var parsedOptions Options

// listenUDP listens for incoming UDP packets and passes them off to
// the database storage.
//
func listenUDP(db *store.SqliteStore) {
	pc, err := net.ListenPacket("udp", parsedOptions.UDPListenAddress)
	if err != nil {
		log.Fatalf("Failed to listen to %s: %v", parsedOptions.UDPListenAddress, err)
	}

	buffer := make([]byte, parsedOptions.UDPBufferSize)
	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading, exiting: %v", err)
		}
		db.PutData(addr, n, buffer[:n])

		if parsedOptions.Verbose {
			log.Printf("DATA> from='%v' packetSize=%d payload'%x'", addr, n, buffer[:n])
		}
	}
}

func main() {
	// Parse command line options
	parser := flags.NewParser(&parsedOptions, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Open database
	db, err := store.New(parsedOptions.DBFilename)
	if err != nil {
		log.Fatalf("Unable to open or create database: %v", err)
	}

	listenUDP(db)
}
