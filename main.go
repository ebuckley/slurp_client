package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	hostPort    string
	hostAddress string
	requestFile string
	serverConn  net.Conn
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("ERROR: we expect 3 arguments, invocation of client should be \"slurp_client server_host server_port file_name directory\"")
	}

	hostAddress = os.Args[1]
	hostPort = os.Args[2]
	requestFile = os.Args[3]
	if len(requestFile) > 255 {
		log.Fatal("WE can't even, your filename is too long!")
		os.Exit(1)
		return
	}

	//set up for interupts
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//application signal listener
	go func() {
		sig := <-sigs
		log.Println("recieved signal", sig)
		done <- true
	}()

	go func() {
		serverConn, err := net.Dial("tcp", hostAddress+":"+hostPort)
		if err != nil {
			log.Fatal("connection was not succesfully made")
			os.Exit(1)
			done <- true
		}

		var buf bytes.Buffer
		buf.WriteString(requestFile)
		req := make([]byte, 255)

		_, err = buf.Read(req)
		if err != nil {
			log.Fatal("WTF? could not read create array of bytes")
			os.Exit(1)
			done <- true
		}

		_, err = serverConn.Write(req)
		if err != nil {
			log.Fatal("error writing request")
			os.Exit(1)
			done <- true
			return
		}

		//recieve filesize
		var fileSize int64

		err = binary.Read(serverConn, binary.BigEndian, &fileSize)
		if err != nil {
			log.Fatal("error reading filesize from socket")
			os.Exit(1)
			done <- true
			return
		}

		log.Println("start recieving n bytes: %s", fileSize)
		fd, err := os.Create(requestFile)
		if err != nil {
			log.Fatal("error creating file")
			os.Exit(1)
			done <- true
			return
		}

		//read from socket, write to file
		n, err := io.CopyN(fd, serverConn, fileSize)
		if err != nil {
			log.Fatal("error copying from socket", err)
			os.Exit(1)
			done <- true
			return
		}

		if n != fileSize {
			log.Fatal("we didn't read all the data.. something went wrong... ")
			os.Exit(1)
			done <- true
			return
		}

		log.Println("Recieved", requestFile)
		done <- true
	}()

	//cleanup before exiting
	<-done
	log.Println("gracefull shutdown, cya")

	if serverConn != nil {
		serverConn.Close()
	}
}
