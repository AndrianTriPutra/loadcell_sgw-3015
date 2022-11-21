package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	timeout    = 4 * time.Second
	scheduller = 5 * time.Second
)

func main() {
	connection, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed connect to socket:%v", err)
	}
	defer connection.Close()

	var start time.Time
	ticker := time.NewTicker(scheduller)
	for {
		<-ticker.C
		str := ""
		start = time.Now()
		elapsed := 1 * time.Nanosecond
		for elapsed < timeout {
			elapsed = time.Since(start)

			buffer := make([]byte, 1024)
			mLen, err := connection.Read(buffer)
			if err != nil {
				log.Println("Error reading:", err.Error())
			}
			buff := fmt.Sprintf("%s", buffer[:mLen])
			str += buff
		}
		str = strings.ReplaceAll(str, "S", "")
		str = strings.ReplaceAll(str, "N", "")
		str = strings.ReplaceAll(str, "T", "")
		str = strings.ReplaceAll(str, ".", "")
		str = strings.ReplaceAll(str, ",", "")
		str = strings.ReplaceAll(str, "+", "")

		lent := len(str)
		last := str[lent-1 : lent]
		if last != "\n" {
			str += "\n"
		}

		var bstr []string
		for len(str) >= 8 && strings.Contains(str, "\n") {
			data := str[0:strings.Index(str, "\n")]
			bstr = append(bstr, data)
			length := len(str)
			str = str[strings.Index(str, "\n")+1 : length]
		}

		var sload []string
		for _, data := range bstr {
			reg := regexp.MustCompile("[^0-9A-Za-z_]")
			data = reg.ReplaceAllString(data, "")
			if len(data) == 8 && strings.Index(data, "Kg") == 6 {
				data := data[0:strings.Index(data, "Kg")]
				sload = append(sload, data)
			}
		}

		var iload []int
		for _, data := range sload {
			digitCheck := regexp.MustCompile("^[0-9]+$")
			if digitCheck.MatchString(data) {
				intload, err := strconv.Atoi(data)
				if err != nil {
					log.Printf("Error convert str to int->%v", err.Error())
				} else {
					iload = append(iload, intload)
				}
			}
		}

		var load int
		for _, data := range iload {
			//log.Printf("data:%v", data)
			if data > load {
				load = data
			}
		}
		log.Printf("load:%v kg", load)

	}
}
