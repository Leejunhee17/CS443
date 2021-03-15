package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import "labrpc"

import "math/rand"
import "time"
// import "fmt"

// import "bytes"
// import "labgob"

func min(x int, y int) int{
	if x > y {
		return y
	}
	return x
}

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type Log struct {
	Command interface{}
	Term 	 	int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (Lab1, Lab2, Challenge1).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	// persistent state on all
	currentTerm int
	votedFor		int
	log					[]Log
	// volatile state on all
	commitIndex	int
	lastApplied	int
	// volatile state on leader
	nextIndex 	[]int
	matchIndex	[]int

	state						 	 int // 0 = Follwer, 1 = Candidate, 2 = Leader
	voteCount				 	 int

	chanAppendEntries  chan bool
	chanRequestVote	   chan bool
	chanVoteResolution chan bool
	chanCommit				 chan bool

	isDead						 bool
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (Lab1).
	// rf.mu.Lock()
	// defer rf.mu.Unlock()
	term = rf.currentTerm
	isleader = false
	if rf.state == 2 {
		isleader = true
	}
	// fmt.Printf("GetState: peer %v, term = %v, state = %v, isleader = %v\n", rf.me, rf.currentTerm, rf.state, isleader)
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (Challenge1).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (Challenge1).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (Lab1, Lab2).
	TERM				 int
	CANDIDATEID	 int
	LASTLOGINDEX int
	LASTLOGTERM  int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (Lab1).
	TERM				int
	VOTEGRANTED bool
}

type AppendEntriesArgs struct {
	TERM				 int
	LEADERID		 int
	PREVLOGINDEX int
	PREVLOGTERM  int
	ENTRIES 		 []Log
	LEADERCOMMIT int
}

type AppendEntriesReply struct {
	TERM		int
	SUCCESS bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (Lab1, Lab2).
	// fmt.Printf("RequestVote: peer%v args = (term %v, candidateid %v, lastlogindex %v, lastlogterm %v)\n", rf.me, args.TERM, args.CANDIDATEID, args.LASTLOGINDEX, args.LASTLOGTERM)
	reply.VOTEGRANTED = false
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.TERM < rf.currentTerm {
		reply.TERM = rf.currentTerm
		// fmt.Printf("RequestVote: VoteReject (I)!!\n")
		return
	}
	if args.TERM > rf.currentTerm {
		// fmt.Printf("RequestVote: RULE I for peer%v\n", rf.me)
		rf.currentTerm = args.TERM
		rf.state = 0
		rf.votedFor = -1
	}
	reply.TERM = rf.currentTerm
	isUptoDate := func() bool {
		if rf.log[len(rf.log) - 1].Term < args.LASTLOGTERM {
			return true
		}
		if rf.log[len(rf.log) - 1].Term == args.LASTLOGTERM {
			if len(rf.log) - 1 <= args.LASTLOGINDEX {
				return true
			}
		}
		return false
	}
	if (rf.votedFor == -1 || rf.votedFor == args.CANDIDATEID) && isUptoDate() {
		// fmt.Printf("RequestVote: VoteGrant!!\n")
		go func() {
			rf.chanRequestVote <- true
		}()
		rf.votedFor = args.CANDIDATEID
		reply.VOTEGRANTED = true
	} else {
		// fmt.Printf("RequestVote: VoteReject (II)!!\n")
	}
	return
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		if rf.state != 1 {
			return ok
		}
		if reply.VOTEGRANTED {
			rf.voteCount++
			if rf.state == 1 && rf.voteCount > len(rf.peers) / 2 {
				go func() {
					rf.chanVoteResolution <- true
				}()
			}
		} else if reply.TERM > rf.currentTerm {
			// fmt.Printf("sendRequestVote: RULE I for peer%v\n", rf.me)
			rf.currentTerm = reply.TERM
			rf.state = 0
			rf.votedFor = -1
		}
	}
	return ok
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	// Your code here (Lab1, Lab2).
	// fmt.Printf("AppendEntries: peer%v args = (term %v, leaderid %v, prevlogindex %v, prevlogterm %v, entries %v, leadercommit %v)\n", rf.me, args.TERM, args.LEADERID, args.PREVLOGINDEX, args.PREVLOGTERM, args.ENTRIES, args.LEADERCOMMIT)
	reply.SUCCESS = false
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.TERM < rf.currentTerm {
		reply.TERM = rf.currentTerm
		// fmt.Printf("AppendEntries: Fail (I)!!\n")
		return
	}
	go func() {
		rf.chanAppendEntries <- true
	}()
	if args.TERM > rf.currentTerm {
		// fmt.Printf("AppendEntries: RULE I for peer%v\n", rf.me)
		rf.currentTerm = args.TERM
		rf.state = 0
		rf.votedFor = -1
	}
	reply.TERM = args.TERM
	if args.PREVLOGINDEX > len(rf.log) - 1 {
		// fmt.Printf("AppendEntries: Fail (II)!!\n")
		return
	}
	if rf.log[args.PREVLOGINDEX].Term != args.PREVLOGTERM {
		// fmt.Printf("AppendEntries: Fail (III)!!\n")
		return
	}
	rf.log = rf.log[:args.PREVLOGINDEX + 1]
	if len(args.ENTRIES) != 0 {
		rf.log = append(rf.log, args.ENTRIES...)
	}
	reply.SUCCESS = true
	if args.LEADERCOMMIT > rf.commitIndex {
		rf.commitIndex = min(args.LEADERCOMMIT, len(rf.log))
		go func() {
			rf.chanCommit <- true
		}()
	}
	// fmt.Printf("AppendEntries: Success!!\n")
	return
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		// fmt.Printf("sendAppendEntries: Peer%v to peer%v with args term = %v, reply term = %v, my term = %v\n", rf.me, server, args.TERM, reply.TERM, rf.currentTerm)
		if rf.state != 2 || args.TERM != rf.currentTerm {
			return ok
		}
		if reply.SUCCESS {
			if len(args.ENTRIES) > 0 {
					rf.nextIndex[server] += len(args.ENTRIES)
					rf.matchIndex[server] = rf.nextIndex[server] - 1
			}
		} else if reply.TERM > rf.currentTerm {
			// fmt.Printf("sendAppendEntries: RULE I for peer%v's term %v < %v\n", rf.me, rf.currentTerm, reply.TERM)
			rf.currentTerm = reply.TERM
			rf.state = 0
			rf.votedFor = -1
		} else {
			rf.nextIndex[server]--
		}
	}
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (Lab2).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state != 2 {
		// fmt.Printf("Start: peer%v with term = %v is not leader\n", rf.me, rf.currentTerm)
		isLeader = false
		return index, term, isLeader
	}
	// fmt.Printf("Start: peer%v with term = %v is leader\n", rf.me, rf.currentTerm)
	index = len(rf.log)
	term = rf.currentTerm
	log := Log{Command: command, Term: term}
	rf.log = append(rf.log, log)
	// fmt.Printf("Start: peer%v log %v start\n", rf.me, rf.log)
	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
	// fmt.Printf("Kill\n")
	rf.isDead = true
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (Lab1, Lab2, Challenge1).
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.log = append(rf.log, Log{Term: 0})
	rf.commitIndex = 0
	rf.lastApplied = 0

	rf.state = 0 // start Follower state
	rf.chanAppendEntries = make(chan bool)
	rf.chanRequestVote = make(chan bool)
	rf.chanVoteResolution = make(chan bool)
	rf.chanCommit = make(chan bool)
	rf.isDead= false

	timeout := 550 + (rand.Int63() % 150)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// fmt.Printf("Make: Start\n")
	go func() {
		cnt := 0
		for {
			if rf.isDead {
				// fmt.Printf("Make %v: peer%v dead\n", cnt, rf.me)
				break
			}
			if rf.state == 0 { // Follower
				// fmt.Printf("Make %v: I'm peer%v, follower (currentTerm %v, log %v, commitIndex %v, lastApplied %v)\n", cnt, rf.me, rf.currentTerm, rf.log, rf.commitIndex, rf.lastApplied)
				select {
				case <-rf.chanAppendEntries:
					// fmt.Printf("Follower %v: Leader alive\n", cnt)
				case <-rf.chanRequestVote:
					// fmt.Printf("Follower %v: Vote is already started\n", cnt)
				case <-time.After(time.Duration(timeout) * time.Millisecond):
					// fmt.Printf("Follower %v: Leader dead, peer%v become candidate\n", cnt, rf.me)
					rf.state = 1
				}
			} else if rf.state == 1 { // Candidate
				rf.mu.Lock()
				rf.currentTerm++
				rf.votedFor = rf.me
				rf.voteCount = 1
				rf.mu.Unlock()
				// fmt.Printf("Make %v: I'm peer%v, my state is Candidate and current term is %v\n", cnt, rf.me, rf.currentTerm)
				go func() {
					rf.mu.Lock()
					defer rf.mu.Unlock()
					args := &RequestVoteArgs{}
					args.TERM = rf.currentTerm
					args.CANDIDATEID = rf.me
					args.LASTLOGINDEX = len(rf.log) - 1
					args.LASTLOGTERM = rf.log[args.LASTLOGINDEX].Term
					for i := 0; i < len(rf.peers); i++ {
						if i != rf.me && rf.state == 1 {
							reply := &RequestVoteReply{}
							// fmt.Printf("Candidate %v: send RequestVote RPC to peer%v\n",cnt, i)
							go rf.sendRequestVote(i, args, reply)
						}
					}
				}()
				select {
				case <-time.After(time.Duration(timeout) * time.Millisecond):
					// fmt.Printf("Candidate %v: Oh! vote failed T.T\n", cnt)
				case <-rf.chanAppendEntries:
					// fmt.Printf("Candidate %v: leader is already elected\n", cnt)
					rf.mu.Lock()
					rf.state = 0
					rf.mu.Unlock()
				case <-rf.chanVoteResolution:
					// fmt.Printf("Candidate %v: peer%v become leader with %vvotes\n", cnt, rf.me, rf.voteCount)
					rf.mu.Lock()
					rf.state = 2
					rf.nextIndex = make([]int, len(rf.peers))
					rf.matchIndex	= make([]int, len(rf.peers))
					for i := 0; i < len(rf.peers); i++ {
						rf.nextIndex[i] = len(rf.log)
						rf.matchIndex[i] = 0
					}
					rf.mu.Unlock()
				}
			} else if rf.state == 2 { // Leader
				// fmt.Printf("Make %v: I'm peer%v, leader (currentTerm %v, log %v, commitIndex %v, lastApplied %v, nextIndex[] %v, matchindex[] %v)\n", cnt, rf.me, rf.currentTerm, rf.log, rf.commitIndex, rf.lastApplied, rf.nextIndex, rf.matchIndex)
				rf.mu.Lock()
				N := rf.commitIndex
				for i := len(rf.log) - 1; i > rf.commitIndex; i-- {
					matchCount := 1
					for j := 0; j < len(rf.peers); j++ {
						if j != rf.me && rf.matchIndex[j] >= i && rf.log[i].Term == rf.currentTerm {
							matchCount++
						}
					}
					if matchCount > len(rf.peers) / 2 {
						N = i
						break
					}
				}
				if N > rf.commitIndex {
					rf.commitIndex = N
					go func() {
						rf.chanCommit <- true
					}()
				}
				for i := 0; i < len(rf.peers); i++ {
					if i != rf.me && rf.state == 2 {
						args := &AppendEntriesArgs{}
						args.TERM = rf.currentTerm
						args.LEADERID = rf.me
						args.PREVLOGINDEX = rf.nextIndex[i] - 1
						args.PREVLOGTERM = rf.log[args.PREVLOGINDEX].Term
						args.ENTRIES = make([]Log, len(rf.log) - args.PREVLOGINDEX - 1)
						copy(args.ENTRIES, rf.log[args.PREVLOGINDEX + 1:])
						args.LEADERCOMMIT = rf.commitIndex
						reply := &AppendEntriesReply{}
						// fmt.Printf("Leader %v: send AppendEntries RPC %v to peer%v\n",cnt, args, i)
						go rf.sendAppendEntries(i, args, reply)
					}
				}
				rf.mu.Unlock()
				time.Sleep(time.Duration(100) * time.Millisecond) // Heartbeat = 10 times per second
			}
			// rf.GetState()
			cnt++
		}
	}()

	go func() {
		for {
			<-rf.chanCommit
			// fmt.Printf("Make: peer%v commit!!\n", rf.me)
			rf.mu.Lock()
			for i := rf.lastApplied + 1; i <= rf.commitIndex; i++ {
				command := rf.log[i].Command
				msg := ApplyMsg{CommandValid: true, CommandIndex: i, Command: command}
				// fmt.Printf("Make: msg %v to channel\n", msg)
				applyCh <- msg
				rf.lastApplied = i
			}
			rf.mu.Unlock()
		}
	}()
	// fmt.Printf("Make: End")
	return rf
}
