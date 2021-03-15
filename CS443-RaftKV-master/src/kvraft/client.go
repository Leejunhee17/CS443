package raftkv

import "labrpc"
import "crypto/rand"
import "math/big"
import (
	// "fmt"
	"sync"
	"time"
)

type Clerk struct {
	servers  []*labrpc.ClientEnd
	// You will have to modify this struct.
	exLeader int
	mu			 sync.Mutex
	xid			 int64
	reqNum   int64
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	// You'll have to add code here.
	ck.exLeader = 0
	ck.xid = nrand()
	ck.reqNum = 0
	return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) Get(key string) string {

	// You will have to modify this function.
	// fmt.Printf("C/GET: function start\n")
	for{
		args := &GetArgs{}
		args.Key = key
		args.Xid = ck.xid
		args.ReqNum = ck.reqNum
		ck.exLeader %= len(ck.servers)
		reply := &GetReply{}
		chanOk := make(chan bool, 1)
		go func(){
			ok := ck.servers[ck.exLeader].Call("KVServer.Get", args, reply)
			chanOk <- ok
		}()
		select {
		case <-time.After(550 * time.Millisecond):
			ck.exLeader++
			continue
		case ok := <-chanOk:
			// fmt.Printf("C/PUTAPPEND: ok = %v, reply = %v\n", ok, reply)
			if ok && !reply.WrongLeader {
				// fmt.Printf("C/GET: function end\n")
				ck.reqNum = nrand()
				return reply.Value
			}
			ck.exLeader++
		}
	}
}

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) PutAppend(key string, value string, op string) {
	// You will have to modify this function.
	// fmt.Printf("C/PUTAPPEND: function start\n")
	for{
		ck.exLeader %= len(ck.servers)
		args := &PutAppendArgs{}
		args.Key = key
		args.Value = value
		args.Op = op
		args.Xid = ck.xid
		args.ReqNum = ck.reqNum
		reply := &PutAppendReply{}
		chanOk := make(chan bool, 1)
		go func(){
			ok := ck.servers[ck.exLeader].Call("KVServer.PutAppend", args, reply)
			chanOk <- ok
		}()
		select {
		case <-time.After(550 * time.Millisecond):
			ck.exLeader++
			continue
		case ok := <-chanOk:
			// fmt.Printf("C/PUTAPPEND: ok = %v, reply = %v\n", ok, reply)
			if ok && !reply.WrongLeader {
				// fmt.Printf("C/PUTAPPEND: function end\n")
				ck.reqNum = nrand()
				return
			}
			ck.exLeader++
		}
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
