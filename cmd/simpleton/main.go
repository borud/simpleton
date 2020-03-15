package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/borud/simpleton/pkg/model"
	"github.com/borud/simpleton/pkg/store"
	"github.com/borud/simpleton/pkg/web"
	"github.com/jessevdk/go-flags"
)

// Options contains the command line options
//
type Options struct {
	// Webserver options
	WebServerListenAddress string `short:"w" long:"webserver-listen-address" description:"Listen address for webserver" default:":8008" value-name:"[<host>]:<port>"`
	WebServerStaticDir     string `short:"s" long:"webserver-static-dir" description:"Static dir for files served through webserver" default:"static" value-name:"<directory>"`

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

	go func() {
		buffer := make([]byte, parsedOptions.UDPBufferSize)
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				log.Printf("Error reading, exiting: %v", err)
			}

			data := model.Data{
				Timestamp:  time.Now(),
				FromAddr:   addr.String(),
				PacketSize: n,
				Payload:    buffer[:n],
			}

			id, err := db.PutData(&data)
			if err != nil {
				log.Printf("Error storing data: %v", err)
				continue
			}

			// Update the assigned id
			data.ID = id

			if parsedOptions.Verbose {
				json, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshalling to JSON: %v", err)
					continue
				}
				log.Printf("DATA> %s", json)
			}
		}
	}()
	log.Printf("Started UDP listener on %s", parsedOptions.UDPListenAddress)
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

	// Listen to UDP socket
	listenUDP(db)

	// Set up webserver
	webServer := web.New(db, parsedOptions.WebServerListenAddress, parsedOptions.WebServerStaticDir)
	webServer.ListenAndServe()
}
