package main

import (
	"sync"
)

// NodeType represents the type of a node in the Raft cluster.
type NodeType string

// Possible node types in Raft.
const (
	Follower  NodeType = "follower"
	Candidate NodeType = "candidate"
	Leader    NodeType = "leader"
)

// LogEntry represents a single entry in the Raft log.
type LogEntry struct {
	Term    int
	Command string
}

// RequestVoteReqBody represents the body of a requestVote RPC.
type RequestVoteReqBody struct {
	Term         int    `json:"term"`
	CandidateId  string `json:"candidateId"`
	LastLogIndex int    `json:"lastLogIndex"`
	LastLogTerm  int    `json:"lastLogTerm"`
}

// AppendEntriesReqBody represents the body of an appendEntries RPC.
type AppendEntriesReqBody struct {
	From         string     `json:"from"`
	Term         int        `json:"term"`
	LeaderId     string     `json:"leaderId"`
	PrevLogIndex int        `json:"prevLogIndex"`
	PrevLogTerm  int        `json:"prevLogTerm"`
	Entries      []LogEntry `json:"entries"`
	LeaderCommit int        `json:"leaderCommit"`
}

// Vote represents a response to a requestVote or appendEntries RPC.
type Vote struct {
	Term        int  `json:"term"`
	VoteGranted bool `json:"voteGranted"`
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
	Leader      string
	State       []int
	Timer       *Timer
	mu          sync.Mutex
}
