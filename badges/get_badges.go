package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type List_Item struct {
	Name  string
	Badge string
}
type Badge_Section struct {
	Title string      `json:"title"`
	Items []List_Item `json:"badgeList`
}

func getBadges() {
	data, err := os.ReadFile("available_badges.md")
	if err != nil {
		panic(err)
	}

	badgeTables := string(data)

	//regex to get section headers after ###
	sectionRegex := regexp.MustCompile(`(?m)^###\s(.+)$`)

	// Regular expression to match the name and badge pairs in the table rows
	rowRegex := regexp.MustCompile(`\|\s(.+?)\s+\|\s(.+?)\s+\|`)

	// Initialize a slice to store sections
	var sections []Badge_Section

	// Find all section titles
	sectionMatches := sectionRegex.FindAllStringSubmatch(badgeTables, -1)

	// Split the data into sections
	sectionsData := sectionRegex.Split(badgeTables, -1)

	for i, section := range sectionMatches {
		sectionTitle := strings.TrimSpace(section[1])

		// Extract rows in the current section
		rows := rowRegex.FindAllStringSubmatch(sectionsData[i+1], -1)

		// Initialize a slice to store items for the current section
		var items []List_Item

		for _, row := range rows {
			name := strings.TrimSpace(row[1])
			badge := strings.TrimSpace(row[2])

			if name == "Name" || strings.Contains(name, "---") {
				continue
			}

			// Create an Item struct for each row
			item := List_Item{
				Name:  name,
				Badge: badge,
			}
			items = append(items, item)
		}

		// Create a Section struct and add it to the sections slice
		sectionObj := Badge_Section{
			Title: sectionTitle,
			Items: items,
		}
		sections = append(sections, sectionObj)
	}

	// Print the sections and items
	for _, section := range sections {
		fmt.Println("Section:", section.Title)
		for _, item := range section.Items {
			fmt.Printf("Name: %s, Badge: %s\n", item.Name, item.Badge)
		}
	}
}

func main() {
	getBadges()
}
