package main

import (
	"fmt"
	"time"
)

type Subject struct {
	name  string
	grade uint
}
type Student struct {
	name       string
	no_subject int
	subjects   []Subject
}

func (s *Student) AddSubject(subject Subject) {
	s.subjects = append(s.subjects, subject)
}

func (s *Student) PrintSubjects() {	
	fmt.Println("You have the following subjects:")
	fmt.Println("Subject\t\t Grade")
	for i, subject := range s.subjects {
		fmt.Printf("%d. %s\t\t %d\n", i+1, subject.name, subject.grade)
	}

	average := s.CalculateAverageGrade()
	if average != -1 {
		fmt.Printf("\n##Average grade for %s is: %.2f\n", s.name, average)
	}
}
func (s *Student) CalculateAverageGrade() float64 {
	if len(s.subjects) == 0 {
		return 0.0
	}
	var total uint
	for _, subject := range s.subjects {
		total += subject.grade
	}
	return float64(total) / float64(len(s.subjects))
}

func IsValidGrade(grade int) bool {
	return grade >= 0 && grade <= 100
}
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
	time.Sleep(100 * time.Millisecond)
}

func WelcomeMessage() (name string, no_subject int) {
	fmt.Println("\n\n***********Welcome to the Student Grade Calculator Console App!! ************")
	fmt.Printf("Please Enter Your name:\n")
	fmt.Scan(&name)
	isValidNoSubject := false
	for !isValidNoSubject {
		fmt.Printf("Please Enter Your number of subjects:\n")
		fmt.Scan(&no_subject)
		isValidNoSubject = no_subject > 0
		if !isValidNoSubject {
			fmt.Println("Invalid number of subjects. Please enter a positive integer.")
		}
	}
	ClearScreen()
	return
}

func main() {

	name, no_subject := WelcomeMessage()
	fmt.Printf("Welcome %s, you have chosen to enter %d subjects.\n", name, no_subject)
	student := Student{name: name, no_subject: no_subject}
	for i := range no_subject {
		var subject string
		var grade int
		fmt.Printf("Subject %d Name:\n", i+1)
		fmt.Scan(&subject)
		fmt.Printf("Subject %d Grade:\n", i+1)
		fmt.Scan(&grade)
		for !IsValidGrade(grade) {
			fmt.Println("Invalid grade. Please enter a grade between 0 and 100.")
			fmt.Printf("Subject %d Grade:\n", i+1)
			fmt.Scan(&grade)
		}

		student.AddSubject(Subject{name: subject, grade: uint(grade)})
	}
	ClearScreen()
	student.PrintSubjects()
}
