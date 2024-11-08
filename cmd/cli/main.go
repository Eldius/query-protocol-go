package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

func main() {

	hostaddr := "192.168.0.42:25565"

	addr, err := net.ResolveUDPAddr("udp4", hostaddr)
	if err != nil {
		err = fmt.Errorf("error resolving host: %w", err)
		panic(err)
	}

	con, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		err = fmt.Errorf("error connecting to server: %w", err)
		panic(err)
	}

	var buf [4]byte

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Read(buf[0:])

	// make sessionID 'minecraft-safe'
	for i := 0; i < 4; i++ {
		buf[i] = buf[i] & 0x0F
	}

	sessionID := buf

	// Build challenge request packet and write to socket
	reqBuf := [7]byte{0xFE, 0xFD, 0x09}
	copy(reqBuf[3:], sessionID[0:])
	if _, err := con.Write(reqBuf[:]); err != nil {
		err = fmt.Errorf("writing to server: %w", err)
		panic(err)
	}

	// read full response from socket
	//var buf [2048]byte
	var res = &bytes.Buffer{}
	defer func() {
		_ = con.SetDeadline(time.Time{})
	}()
	// A simple read loop, this function handles multi-packet responses until EOF
	for {
		if err := con.SetDeadline(time.Now().Add(10 * time.Millisecond)); err != nil {
			err = fmt.Errorf("setting deadline: %w", err)
			panic(err)
		}
		b, err := con.Read(buf[0:])
		if b > 0 {
			res.Write(buf[:b])
		}
		if err == io.EOF || b < 2048 {
			break
		}
		if b == 0 && err != io.EOF {
			err = fmt.Errorf("error reading from server: %w", err)
			panic(err)
		}

	}
	fmt.Printf("%s\n", res.String())

	//// ensure our response header is good2go
	//err = req.verifyResponseHeader(resBuf)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// read until end delimiter
	//res, err := resBuf.ReadBytes(0x00)
	//if err != nil {
	//	return nil, errors.New("malformed challenge response")
	//}
	//
	//// chop off tailing null byte and convert to string, then to int
	//tokenString := string(res[:len(res)-1])
	//tokenInt, err := strconv.ParseInt(tokenString, 10, 32)
	//if err != nil {
	//	return nil, errors.New("malformed challenge response")
	//}
	//
	//// Convert our integer to byte array and return
	//tokenBuf := &bytes.Buffer{}
	//binary.Write(tokenBuf, binary.BigEndian, tokenInt)
	//tokenBytes := tokenBuf.Bytes()
	//return tokenBytes[len(tokenBytes)-4:], nil
}
