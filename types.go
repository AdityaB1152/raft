package main

import (
	"sync"
)

// NodeType represents the type of a node in the Raft cluster.
type NodeType int

const (
	Follower NodeType = iota
	Candidate
	Leader
)

// LogEntry represents a single entry in the Raft log.
type LogEntry struct {
	Term    int
	Command string
}

type RequestVoteReqBody struct {
	Term         int
	CandidateId  string
	LastLogIndex int
	LastLogTerm  int
}

type AppendEntriesReqBody struct {
	From         string
	Term         int
	LeaderId     string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type Vote struct {
	Term        int
	VoteGranted bool
}

// RaftNode represents a node in the Raft cluster.
type RaftNode struct {
	ID          string
	CurrentTerm int
	VotedFor    string
	Log         []LogEntry
	Type        NodeType
	Peers       []string
	CommitIndex int
	LastApplied int
	Leader      string
	State       []int
	Timer       *Timer
	mu          sync.Mutex
}
