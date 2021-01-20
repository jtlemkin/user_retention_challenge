package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// startTime represents the timestamp of the earliest date we are considering
const startTime int64 = 1451606400

// numDays is the maximum number of days that we're considering
const numDays = 14

type usersSet = map[int64]bool

// getDayID converts a timestamp to a day ID
func getDayID(timestamp int64) int {
	if timestamp < startTime {
		return -1
	}

	var day int64 = 86400

	id := (timestamp - startTime) / day

	if id > numDays {
		return -1
	}

	return int(id)
}

// getDisjointUsers returns the users that are in first but not in second
func getDisjointUsers(first usersSet, second usersSet) usersSet {
	disjoint := make(usersSet)

	for user := range first {
		if !second[user] {
			disjoint[user] = true
		}
	}

	return disjoint
}

// parseActivities creates an array of day structs from an input csv file
func parseActivities(file *os.File) [numDays]usersSet {
	reader := csv.NewReader(file)

	var days [numDays]usersSet
	for i := range days {
		days[i] = make(usersSet)
	}

	for {
		activity, err := reader.Read()
		if err == io.EOF {
			break
		}

		timestamp, _ := strconv.ParseInt(activity[0], 10, 32)
		user, _ := strconv.ParseInt(activity[1], 10, 64)

		// TODO: error checking for Atoi conversions
		start := getDayID(timestamp)
		days[start][user] = true
	}

	return days
}

// composeLine generates the line to write to the CSV file
// start is the index of the first day
func composeLine(start int, activities [numDays]usersSet) string {
	line := strconv.Itoa(start + 1)

	var continuingUsers usersSet
	if start == 0 {
		continuingUsers = activities[start]
	} else {
		continuingUsers = getDisjointUsers(activities[start], activities[start-1])
	}

	for i := 0; i < numDays; i++ {
		line += ","

		if start+i >= numDays {
			// Day not in data set
			line += "0"
		} else if start+i == numDays-1 {
			// Last day in data set
			// All continuing users terminate here
			line += strconv.Itoa(len(continuingUsers))
		} else {
			// Terminating users are current users not present in the next day's users
			terminatingUsers := getDisjointUsers(continuingUsers, activities[start+i+1])
			// The new current users are the current users who are not terminating
			continuingUsers = getDisjointUsers(continuingUsers, terminatingUsers)

			line += strconv.Itoa(len(terminatingUsers))
		}
	}

	return line
}

// WriteCSV writes the output CSV file
func writeCSV(activities [numDays]usersSet) {
	f, _ := os.Create("output.csv")
	// TODO: Error checking on file creation
	defer f.Close()

	for i := range activities {
		line := composeLine(i, activities)
		f.WriteString(line + "\n")
	}
}

func main() {
	path := os.Args[1]
	file, _ := os.Open(path)

	// TODO: Check err variable from os.Open

	activities := parseActivities(file)
	writeCSV(activities)
}
