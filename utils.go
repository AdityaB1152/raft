package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getTimeout() time.Duration {
	return time.Duration(getRandomIntInclusive(5000, 10000)) * time.Millisecond
}

func getRandomIntInclusive(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func countVotes(votes []Vote) (int, int) {
	trueCount, falseCount := 0, 0

	for _, vote := range votes {
		if vote.VoteGranted {
			trueCount++
		} else {
			falseCount++
		}
	}

	return trueCount, falseCount
}

func executeCommand(node *RaftNode, command string) {
	switch command {
	case "Set x=1":
		node.State[0] = 1
	case "Set x=2":
		node.State[0] = 2
	case "Set x=3":
		node.State[0] = 3
	case "Set y=1":
		node.State[1] = 1
	case "Set y=2":
		node.State[1] = 2
	case "Set y=3":
		node.State[1] = 3
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateRandomCommand generates a random command to change the value of x or y.
func GenerateRandomCommand() string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Define possible commands for x and y with their values
	commands := []string{
		"Set x=1",
		"Set x=2",
		"Set x=3",
		"Set y=1",
		"Set y=2",
		"Set y=3",
	}

	// Randomly select a command
	randomIndex := rand.Intn(len(commands))
	return commands[randomIndex]
}
