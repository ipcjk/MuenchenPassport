package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

)

func main() {
	if len(os.Args) <= 1  {
		log.Fatal("Zu wenig Argumente")
	}

	data := url.Values{}
	data.Add("ifNummer", os.Args[1])
	data.Add("pbAbfragen", "Abfragen")

	req, err := http.NewRequest("POST", "https://www17.muenchen.de/Passverfolgung/", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Ein Dokument mit dieser Nummer ist nicht vorhanden") {
			fmt.Println(os.Args[1], " - diese Passnummer ist falsch!")
		} else if strings.Contains(scanner.Text(), " noch nicht zur Abholung bereit.")  {
			fmt.Println(os.Args[1], "liegt noch nicht zur Abholung bereit.")
		} else if strings.Contains(scanner.Text(), `liegt zur<B STYLE="color: green"> Abholung bereit.</B></TD>`) {
			fmt.Println(os.Args[1],  "liegt zur Abholung bereit, yes!")
		}
	}

}
