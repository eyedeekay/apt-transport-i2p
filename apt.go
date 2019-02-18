package apti2p

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Message struct {
	Status     string
	StatusCode int
	Header     Header
	Exit       int
}

type AptMessage struct {
	Code    int
	Headers map[string]string
}

type Header map[string][]string

func (h Header) Add(key, value string) {
	h[key] = append(h[key], value)
}

func (h Header) Get(key string) string {
	if value, ok := h[key]; ok {
		if len(value) > 0 {
			return value[0]
		}
	}
	return ""
}

func (m *Message) String() string {
	s := []string{m.Status}
	for k, values := range m.Header {
		for _, v := range values {
			s = append(s, k+": "+v)
		}
	}
	return strings.Join(s, "\n") + "\n\n"
}

type AptMethod struct {
}

func (a *AptMethod) Run() {
	c := make(chan *Message)
	go output(c)
	sendCapabilities(c)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		line := stdin.Text()
		if line == "" {
			continue
		}
		fmt.Println(line)
		s := strings.SplitN(line, " ", 2)
		code, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("Malformed message!")
			os.Exit(100)
		}
		request := &Message{
			Status:     line,
			StatusCode: code,
			Header:     Header{},
		}

		for stdin.Scan() {
			line := stdin.Text()

			if line == "" {
				process(c, request)
				break
			}
			s := strings.SplitN(line, ": ", 2)
			request.Header.Add(s[0], s[1])
		}
	}
}

func output(c <-chan *Message) {
	for {
		m := <-c
		os.Stdout.Write([]byte(m.String()))
		if m.Exit != 0 {
			os.Exit(m.Exit)
		}
	}
}

func sendCapabilities(c chan<- *Message) {
	caps := &Message{
		Status: "100 Capabilities",
		Header: Header{},
	}

	caps.Header.Add("Version", "1.2")
	caps.Header.Add("Pipeline", "true")
	caps.Header.Add("Send-Config", "true")

	c <- caps
}

func process(c chan<- *Message, m *Message) {
	switch m.StatusCode {
	case 600:
		go fetch(c, m)
	case 601:
		// TODO: parse config?
	default:
		fail := &Message{
			Status: "401 General Failure",
			Header: Header{},
			Exit:   100,
		}
		fail.Header.Add("Message", "Status code not implemented")

		c <- fail
	}
}

func fetch(c chan<- *Message, m *Message) {
	uri := m.Header.Get("URI")
	filename := m.Header.Get("Filename")

	// TODO: If-Modified-Since

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Could not open file: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}
	defer file.Close()

	// TODO: Fix bug with appending to existing files
	// TODO: implement range requests if file already exists

	realURI := TranslateAddr(uri)

	log.Println("Get: ", realURI)
	aptClient, err := InitClient()
	if err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Client generation error: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}
	resp, err := aptClient.Get(realURI)
	if err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Could not fetch URI: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}
	defer resp.Body.Close()

	started := &Message{
		Status: "200 URI Start",
		Header: Header{
			"URI": []string{uri},
		},
	}
	// TODO: add Last-Modified header

	c <- started

	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()
	sha512Hash := sha512.New()

	if _, err = io.Copy(io.MultiWriter(file, md5Hash, sha1Hash, sha256Hash, sha512Hash), resp.Body); err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Could not write file: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}

	success := &Message{
		Status: "201 URI Done",
		Header: Header{},
	}
	success.Header.Add("URI", uri)
	success.Header.Add("Filename", filename)
	// TODO Size, Last-Modified
	md5Hex := hex.EncodeToString(md5Hash.Sum(nil)[:])
	success.Header.Add("MD5-Hash", md5Hex)
	success.Header.Add("MD5Sum-Hash", md5Hex)
	success.Header.Add("SHA1-Hash", hex.EncodeToString(sha1Hash.Sum(nil)[:]))
	success.Header.Add("SHA256-Hash", hex.EncodeToString(sha256Hash.Sum(nil)[:]))
	success.Header.Add("SHA512-Hash", hex.EncodeToString(sha512Hash.Sum(nil)[:]))

	c <- success
}
