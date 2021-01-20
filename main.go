package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// StartTime represents the timestamp of the earliest date we are considering
const StartTime int64 = 1451606400

// NumDays is the maximum number of days that we're considering
const NumDays = 14

type usersSet = map[int64]bool

// GetDayID converts a timestamp to a day ID
func GetDayID(timestamp int64) int {
	if timestamp < StartTime {
		return -1
	}

	var day int64 = 86400

	id := (timestamp - StartTime) / day

	if id > NumDays {
		return -1
	}

	return int(id)
}

// GetDisjointUsers returns the users that are in first but not in second
func GetDisjointUsers(first usersSet, second usersSet) usersSet {
	disjoint := make(usersSet)

	for user := range first {
		if !second[user] {
			disjoint[user] = true
		}
	}

	return disjoint
}

// ParseActivities creates an array of day structs from an input csv file
func ParseActivities(file *os.File) [NumDays]usersSet {
	reader := csv.NewReader(file)

	var days [NumDays]usersSet
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
		start := GetDayID(timestamp)
		days[start][user] = true
	}

	return days
}

// ComposeLine generates the line to write to the CSV file
// start is the index of the first day
func ComposeLine(start int, activities [NumDays]usersSet) string {
	line := strconv.Itoa(start + 1)

	var continuingUsers usersSet
	if start == 0 {
		continuingUsers = activities[start]
	} else {
		continuingUsers = GetDisjointUsers(activities[start], activities[start-1])
	}

	for i := 0; i < NumDays; i++ {
		line += ","

		if start+i >= NumDays {
			// Day not in data set
			line += "0"
		} else if start+i == NumDays-1 {
			// Last day in data set
			// All continuing users terminate here
			line += strconv.Itoa(len(continuingUsers))
		} else {
			// Terminating users are current users not present in the next day's users
			terminatingUsers := GetDisjointUsers(continuingUsers, activities[start+i+1])
			// The new current users are the current users who are not terminating
			continuingUsers = GetDisjointUsers(continuingUsers, terminatingUsers)

			line += strconv.Itoa(len(terminatingUsers))
		}
	}

	return line
}

// WriteCSV writes the output CSV file
func WriteCSV(activities [NumDays]usersSet) {
	f, _ := os.Create("output.csv")
	// TODO: Error checking on file creation
	defer f.Close()

	for i := range activities {
		line := ComposeLine(i, activities)
		f.WriteString(line + "\n")
	}
}

func main() {
	path := os.Args[1]
	file, _ := os.Open(path)

	// TODO: Check err variable from os.Open

	activities := ParseActivities(file)
	WriteCSV(activities)
}
