package raftkv

import (
	"labgob"
	"labrpc"
	"log"
	"raft"
	"sync"
	// "fmt"
	"time"
)

const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	Optype 	 string
	Key 		 string
	Value 	 string
	Xid 		 int64
	ReqNum 	 int64
}

type KVServer struct {
	mu      		 sync.Mutex
	me      		 int
	rf      		 *raft.Raft
	applyCh 		 chan raft.ApplyMsg

	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	kvDB				 map[string]string
	chanCommits	 map[int]chan Op
	xidDB 			 map[int64]int64
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	// fmt.Printf("S/GET: start!!\n")
	_, isLeader := kv.rf.GetState()
	if !isLeader {
		reply.WrongLeader = true
		return
	}
	// fmt.Printf("S/GET: server%v leader? %v!!\n", kv.me, isLeader)
	op := Op{}
	op.Optype = "Get"
	op.Key = args.Key
	op.Xid = args.Xid
	op.ReqNum = args.ReqNum
	var index int
	index, _, isLeader = kv.rf.Start(op)
	if !isLeader {
		// fmt.Printf("S/GET: server%v is not leader after rf.Start\n", kv.me)
		reply.WrongLeader = true
		return
	}
	kv.mu.Lock()
	chanCommit, exists := kv.chanCommits[index]
	if !exists {
		chanCommit = make(chan Op, 1)
		kv.chanCommits[index] = chanCommit
	}
	kv.mu.Unlock()
	select{
	case committed := <-chanCommit:
		if committed == op {
			reply.WrongLeader = false
			kv.mu.Lock()
			val, exists := kv.kvDB[args.Key]
			if !exists {
				reply.Err = ErrNoKey
			} else {
				reply.Value = val
				kv.xidDB[args.Xid] = args.ReqNum
				// fmt.Printf("S/GET: server%v key %v -> value %v\n", kv.me, args.Key, reply.Value)
			}
			kv.mu.Unlock()
		} else {
			// fmt.Printf("S/GET: committed op is not equal to start op\n")
			reply.WrongLeader = true
		}
	case <-time.After((550 + 150 + 100 + 100) * time.Millisecond):
		// fmt.Printf("S/GET: server%v time out\n", kv.me)
		reply.WrongLeader = true
	}
}

func (kv *KVServer) PutAppend(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	// fmt.Printf("S/PUTAPPEND: start!!\n")
	_, isLeader := kv.rf.GetState()
	if !isLeader {
		reply.WrongLeader = true
		return
	}
	// fmt.Printf("S/PUTAPPEND: server%v start with %v!!\n", kv.me, *args)
	op := Op{}
	op.Optype = args.Op
	op.Key = args.Key
	op.Value = args.Value
	op.Xid = args.Xid
	op.ReqNum = args.ReqNum
	var index int
	index, _, isLeader = kv.rf.Start(op)
	if !isLeader {
		// fmt.Printf("S/PUTAPPEND: server%v is not leader after rf.Start\n", kv.me)
		reply.WrongLeader = true
		return
	}
	kv.mu.Lock()
	chanCommit, exists := kv.chanCommits[index]
	if !exists {
		chanCommit = make(chan Op, 1)
		kv.chanCommits[index] = chanCommit
	}
	kv.mu.Unlock()
	select{
	case committed := <-chanCommit:
		if committed == op {
			reply.WrongLeader = false
		} else {
			// fmt.Printf("S/PUTAPPEND: committed op is not equal to start op\n")
			reply.WrongLeader = true
		}
		// fmt.Printf("S/PUTAPPEND: server%v key %v -> value %v\n", kv.me, args.Key, kv.kvDB[args.Key])
		return
	case <-time.After((550 + 150 + 100 + 100) * time.Millisecond):
		// fmt.Printf("S/PUTAPPEND: server%v time out\n", kv.me)
		reply.WrongLeader = true
	}
}

//
// the tester calls Kill() when a KVServer instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (kv *KVServer) Kill() {
	kv.rf.Kill()
	// Your code here, if desired.
}

//
// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots with persister.SaveSnapshot(),
// and Raft should save its state (including log) with persister.SaveRaftState().
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
//
func StartKVServer(servers []*labrpc.ClientEnd, me int, persister *raft.Persister, maxraftstate int) *KVServer {
	// call labgob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	labgob.Register(Op{})

	kv := new(KVServer)
	kv.me = me
	kv.maxraftstate = maxraftstate

	// You may need initialization code here.

	kv.applyCh = make(chan raft.ApplyMsg)
	kv.rf = raft.Make(servers, me, persister, kv.applyCh)
	kv.kvDB = make(map[string]string)
	kv.chanCommits = make(map[int]chan Op)
	kv.xidDB = make(map[int64]int64)

	// You may need initialization code here.
	// fmt.Printf("S/STARTKV: function start\n")
	go func(){
		for {
			msg := <-kv.applyCh
			// fmt.Printf("S/STARTKV: consensus end with %v\n", msg)
			op := msg.Command.(Op)
			index := msg.CommandIndex
			kv.mu.Lock()
			// fmt.Printf("S/STARTKV: xidDB = %v\n", kv.xidDB)
			reqNum, exists := kv.xidDB[op.Xid]
			if !exists || reqNum != op.ReqNum {
				if op.Optype == "Put" {
					kv.kvDB[op.Key] = op.Value
				} else if op.Optype == "Append" {
					kv.kvDB[op.Key] += op.Value
				}
				kv.xidDB[op.Xid] = op.ReqNum
				// fmt.Printf("S/STARTKV: kvDB changed to %v\n", kv.kvDB)
			}
			chanCommit, exists := kv.chanCommits[index]
			if exists {
				chanCommit <- op
			} else {
				chanCommit = make(chan Op, 1)
			}
			kv.mu.Unlock()
		}
	}()

	return kv
}
