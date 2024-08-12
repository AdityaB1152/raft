package main

import (
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
	case "0":
		node.State[0]++
	case "1":
		node.State[1]++
	case "2":
		node.State[2]++
	case "3":
		node.State[3]++
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
