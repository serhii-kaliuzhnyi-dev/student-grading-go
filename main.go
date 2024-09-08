package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the header and ignore it
	_, err = reader.Read() // Skip the first row (header)
	if err != nil {
		panic(err)
	}

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Create a slice to hold the parsed students
	var students []student

	// Loop through the records and parse each student
	for _, record := range records {
		// Parse the test scores (convert from string to int)
		test1, _ := strconv.Atoi(record[3])
		test2, _ := strconv.Atoi(record[4])
		test3, _ := strconv.Atoi(record[5])
		test4, _ := strconv.Atoi(record[6])

		// Create a student struct and append it to the list
		students = append(students, student{
			firstName:   record[0],
			lastName:    record[1],
			university:  record[2],
			test1Score:  test1,
			test2Score:  test2,
			test3Score:  test3,
			test4Score:  test4,
		})
	}

	return students
}

func calculateGrade(students []student) []studentStat {
	var gradedStudents []studentStat

	for _, s := range students {
		// Calculate the final score as the average of the test scores
		finalScore := float32(s.test1Score+s.test2Score+s.test3Score+s.test4Score) / 4

		// Determine the grade based on the final score
		var grade Grade
		switch {
		case finalScore >= 70:
			grade = A
		case finalScore >= 50:
			grade = B
		case finalScore >= 35:
			grade = C
		default:
			grade = F
		}

		// Append the studentStat to the gradedStudents slice
		gradedStudents = append(gradedStudents, studentStat{
			student:    s,
			finalScore: finalScore,
			grade:      grade,
		})
	}

	return gradedStudents
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	if len(gradedStudents) == 0 {
		return studentStat{} // return empty if no students
	}

	topStudent := gradedStudents[0]
	for _, s := range gradedStudents[1:] {
		if s.finalScore > topStudent.finalScore {
			topStudent = s
		}
	}

	return topStudent
}

func findTopperPerUniversity(gradedStudents []studentStat) map[string]studentStat {
	toppers := make(map[string]studentStat)

	for _, s := range gradedStudents {
		// Check if the university already has a topper
		top, exists := toppers[s.university]
		if !exists || s.finalScore > top.finalScore {
			// Update the topper for the university if the current student has a higher score
			toppers[s.university] = s
		}
	}

	return toppers
}
