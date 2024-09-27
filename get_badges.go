package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/kyleochata/md_maker/badge"
)

// type Badge_Section struct {
// 	Title string       `json:"title"`
// 	Items []badge.Item `json:"badgeList`
// }

func PullData() []badge.Item {
	data, err := os.ReadFile("available_badges.md")
	if err != nil {
		panic(err)
	}

	badgeTables := string(data)
	sectionRegex := regexp.MustCompile(`(?m)^###\s(.+)$`)
	rowRegex := regexp.MustCompile(`\|\s(.+?)\s+\|\s(.+?)\s+\|`)
	sectionMatches := sectionRegex.FindAllStringSubmatch(badgeTables, -1)

	// Split the data into sections
	sectionsData := sectionRegex.Split(badgeTables, -1)
	var items []badge.Item
	for i, section := range sectionMatches {
		sectionTitle := strings.TrimSpace(section[1])
		items = append(items, badge.Item{Name: sectionTitle, IsSection: true})
		// Extract rows in the current section
		rows := rowRegex.FindAllStringSubmatch(sectionsData[i+1], -1)
		for _, row := range rows {
			name := strings.TrimSpace(row[1])
			badgeText := strings.TrimSpace(row[2])
			if name == "Name" || strings.Contains(name, "---") {
				continue
			}
			item := badge.Item{Name: name, Badge: badgeText}
			items = append(items, item)
		}
	}
	// for i, item := range items {
	// 	fmt.Println(i, item.Name, item.IsSection)
	// }
	return items
}

// func main() {
// 	GetBadges()
// }
