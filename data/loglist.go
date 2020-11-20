package data

import (
	"sync"
)

type LogList struct {
	Log []*Entry

	mut sync.Mutex
}

func NewLogList() *LogList {

	dPP := PrePrepareArgs{
		View:     1,
		Seq:      0,
		Commands: nil,
	}

	nonce := &Entry{
		PP:            &dPP,
		PreparedCert:  nil,
		Prepared:      true,
		CommittedCert: nil,
		Committed:     true,
		PreEntryHash:  &EntryHash{},
		Digest:        &EntryHash{},
		PrepareHash:   &EntryHash{},
		CommitHash:    &EntryHash{},
	}

	return &LogList{
		Log: append([]*Entry{}, nonce),
	}
}

func (logList *LogList) Put(ent *Entry) (ok bool) {
	logList.mut.Lock()
	defer logList.mut.Unlock()

	listLen := uint32(len(logList.Log))

	if listLen < ent.PP.Seq {
		ok = false
	} else if listLen == ent.PP.Seq {
		ent.PreEntryHash = logList.Log[listLen-1].Digest
		logList.Log = append(logList.Log, ent)
		ok = true
	} else { // listLen > ent.PP.Seq
		panic(`try to overwrite an entry`)
	}
	return
}

func (logList *LogList) Get(eid *EntryID) (ent *Entry, ok bool) {
	logList.mut.Lock()
	defer logList.mut.Unlock()

	if uint32(len(logList.Log)) <= eid.N {
		ent = nil
		ok = false
	} else {
		ent = logList.Log[eid.N]
		if ent.PP.View != eid.V {
			panic(`view not equal with the same seq...`)
		}
		ok = true
	}
	return
}

func (logList *LogList) GetEntryBySeq(seq uint32) (ent *Entry, ok bool) {
	logList.mut.Lock()
	defer logList.mut.Unlock()

	if uint32(len(logList.Log)) <= seq {
		ent = nil
		ok = false
	} else {
		ent = logList.Log[seq]
		ok = true
	}
	return
}

func (logList *LogList) GetApplyCmds(seq uint32) (*[]Command, uint32) {
	logList.mut.Lock()
	defer logList.mut.Unlock()

	commands := make([]Command, 0)

	for seq < uint32(len(logList.Log)) {
		ent := logList.Log[seq]
		if ent.Committed {
			commands = append(commands, *ent.PP.Commands...)
			seq++
		} else {
			break
		}
	}
	return &commands, seq - 1
}
