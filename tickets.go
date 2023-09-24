/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * tickets.go
 */

package bond

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Tickets struct {
	ID       string      `bson:"id"`
	Versions [][2]string `bson:"versions"`
}

var ticketsIns *map[string]Tickets

// GetTickets returns *map[string]Tickets instance
func GetTickets() *map[string]Tickets {
	if ticketsIns == nil {
		filename := "tickets.json"
		data, err := os.ReadFile(filename)
		if err != nil {
			url := "https://raw.githubusercontent.com/simagix/bond/main/tickets.json"
			log.Println("download driver manifest from", url)
			resp, err := http.Get(url)
			if err != nil {
				return nil
			}
			defer resp.Body.Close()

			if data, err = io.ReadAll(resp.Body); err != nil {
				return nil
			}
			fname := "tickets.temp"
			log.Println("write tickets manifest to", fname)
			_ = os.WriteFile(fname, data, 0644)
		}

		// Parse JSON data into a map
		var m map[string]Tickets
		if err := json.Unmarshal(data, &m); err != nil {
			return nil
		}
		ticketsIns = &m
	}
	return ticketsIns
}

func CheckUpgradeRecommendation(mongoVersion string) string {
	tickets := GetTickets()
	list := []string{}
	for key, ticket := range *tickets {
		for _, versions := range ticket.Versions {
			if mongoVersion >= versions[0] && mongoVersion <= versions[1] {
				list = append(list, fmt.Sprintf("<a href='%s'>%s</a>", ticket.ID, key))
			}
		}
	}
	if len(list) <= 0 {
		return ""
	} else if len(list) == 1 {
		return list[0]
	} else if len(list) == 2 {
		return fmt.Sprintf("%s and %s", list[0], list[1])
	}
	length := len(list)
	return fmt.Sprintf("%s, and %s", strings.Join(list[:length-1], ", "), list[length-1])
}
