package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	ip          = ""
	port        int64
	InitDone    = false
	NewDone     = false
)

// init initializes the logger and parses CMD args.
func init() {

	logFile, err := os.OpenFile("ELN_installer.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(3)
	}()

}

func main() {
	server_folder_info, err := os.Stat("/var/eln_file_server")
	defer os.Exit(0)
	defer cleanup()

	if !is_root() {
		handle_error(errors.New("Please run as root (Hint: sudo su)!!"))
	}

	if err != nil {
		on_cleanup(func() {
			if !InitDone {
				InfoLogger.Println("Remove /var/eln_file_server")
				cmd("rm -R /var/eln_file_server")
			}
		})
		run_int_setup()
		server_folder_info, err = os.Stat("/var/eln_file_server")
		handle_error(err)
		InfoLogger.Printf("\n--------------------------------\n--------------------------------\nGo to:\nhttps://%s\n--------------------------------\n--------------------------------\n\n", get_ip())
	}

	if !server_folder_info.IsDir() {
		ErrorLogger.Fatal("/var/eln_file_server is not a directory")
	}

	if contains(os.Args, "-u") {
		InfoLogger.Println("------Running updates!!!!------")
		run_update()
	} else {
		InfoLogger.Println("------Start new Instance!!!!------")
		run_insatnce_setup()
	}
}
