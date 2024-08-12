package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func NewRaftNode(id string, peers []string) *RaftNode {
	node := &RaftNode{
		ID:          id,
		CurrentTerm: 0,
		VotedFor:    "",
		Log:         []LogEntry{},
		Type:        Follower,
		Peers:       peers,
		CommitIndex: 0,
		Leader:      "",
		State:       make([]int, 4),
		Timer:       NewTimer(),
	}

	node.Timer.Start(getTimeout())
	node.Timer.OnTimeout(node.handleTimeout)

	return node
}

func (n *RaftNode) handleTimeout() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.Timer.Stop()
	n.Type = Candidate
	n.votingProcedure()
}

func (n *RaftNode) votingProcedure() {
	fmt.Printf("Voting process started by %s\n", n.ID)

	n.CurrentTerm++
	n.VotedFor = n.ID

	var lastLogIndex int
	var lastLogTerm int

	if len(n.Log) > 0 {
		lastLogIndex = len(n.Log) - 1
		lastLogTerm = n.Log[lastLogIndex].Term
	} else {
		lastLogIndex = -1
		lastLogTerm = -1
	}

	voteResults := make(chan bool, len(n.Peers))

	for _, peer := range n.Peers {
		go func(peer string) {
			body := RequestVoteReqBody{
				Term:         n.CurrentTerm,
				CandidateId:  n.ID,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			resp, err := postRequest(fmt.Sprintf("http://localhost:%s/requestVote", peer), body)
			if err != nil {
				fmt.Printf("Error requesting vote from %s:\n", peer)
				voteResults <- false
				return
			}
			voteResults <- resp.VoteGranted
		}(peer)
	}

	// Collect votes
	trueCount := 1
	falseCount := 0

	for i := 0; i < len(n.Peers); i++ {
		voteGranted := <-voteResults
		if voteGranted {
			trueCount++
		} else {
			falseCount++
		}
	}

	if n.Type == Candidate && trueCount > len(n.Peers)/2 {
		fmt.Printf("%s becomes the leader with %d votes\n", n.ID, trueCount)
		n.Type = Leader
		n.sendHeartBeats()
	} else {
		fmt.Printf("%s remains a follower with %d votes\n", n.ID, trueCount)
		n.Type = Follower
		n.VotedFor = ""
		n.Timer.Reset()
	}
}

func (n *RaftNode) sendHeartBeats() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if n.Type != Leader {
			return
		}

		var prevLogIndex int
		var prevLogTerm int

		// Check if the log is empty before accessing elements
		if len(n.Log) > 0 {
			prevLogIndex = len(n.Log) - 1
			prevLogTerm = n.Log[prevLogIndex].Term
		} else {
			prevLogIndex = -1
			prevLogTerm = -1
		}

		for _, peer := range n.Peers {
			go func(peer string) {
				body := AppendEntriesReqBody{
					From:         n.ID,
					Term:         n.CurrentTerm,
					LeaderId:     n.ID,
					PrevLogIndex: prevLogIndex,
					PrevLogTerm:  prevLogTerm,
					Entries:      []LogEntry{},
					LeaderCommit: n.CommitIndex,
				}
				_, err := postRequest(fmt.Sprintf("http://localhost:%s/appendEntries", peer), body)
				if err != nil {

				}
			}(peer)
		}

		fmt.Printf("Leader %s sent heartbeats\n", n.ID)
	}
}

func postRequest(url string, body interface{}) (*Vote, error) {
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vote Vote
	err = json.NewDecoder(resp.Body).Decode(&vote)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Port number is required as an argument")
	}

	port := os.Args[1]
	peers := []string{"3000", "4000", "5000", "6000", "7000"}

	// Remove the current node's port from the peers list
	for i, p := range peers {
		if p == port {
			peers = append(peers[:i], peers[i+1:]...)
			break
		}
	}

	raftNode := NewRaftNode(port, peers)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("working"))
	})

	http.HandleFunc("/requestVote", func(w http.ResponseWriter, r *http.Request) {
		var req RequestVoteReqBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		raftNode.mu.Lock()
		defer raftNode.mu.Unlock()

		if req.Term < raftNode.CurrentTerm {
			json.NewEncoder(w).Encode(Vote{Term: raftNode.CurrentTerm, VoteGranted: false})
			return
		}

		if raftNode.VotedFor == "" || raftNode.VotedFor == req.CandidateId {
			raftNode.VotedFor = req.CandidateId
			json.NewEncoder(w).Encode(Vote{Term: req.Term, VoteGranted: true})
			return
		}

		json.NewEncoder(w).Encode(Vote{Term: req.Term, VoteGranted: false})
	})

	http.HandleFunc("/appendEntries", func(w http.ResponseWriter, r *http.Request) {
		var req AppendEntriesReqBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		raftNode.mu.Lock()
		defer raftNode.mu.Unlock()

		if len(req.Entries) == 0 {
			fmt.Printf("Received HeartBeat from ID %s\n", req.From)
			raftNode.Leader = req.From
			raftNode.Timer.Reset()
			return
		}

		fmt.Println("Log entries: ", req.Entries)
		fmt.Println("State before: ", raftNode.State)

		for _, entry := range req.Entries {
			executeCommand(raftNode, entry.Command)
		}

		fmt.Println("State after: ", raftNode.State)

		if raftNode.Type == Candidate {
			raftNode.Type = Follower
		}

		if req.Term < raftNode.CurrentTerm {
			json.NewEncoder(w).Encode(Vote{Term: raftNode.CurrentTerm, VoteGranted: false})
			return
		}

		if raftNode.Log[req.PrevLogIndex].Term != req.PrevLogTerm {
			json.NewEncoder(w).Encode(Vote{Term: req.Term, VoteGranted: false})
			return
		}

		if req.LeaderCommit > raftNode.CommitIndex {
			raftNode.CommitIndex = min(req.LeaderCommit, req.PrevLogIndex)
		}

		json.NewEncoder(w).Encode(Vote{Term: req.Term, VoteGranted: true})
	})

	http.ListenAndServe(":"+port, nil)
}
