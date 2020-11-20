package data

import "sync"

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

func (lbbft *LBBFTCore) GetEntryBySeq(seq uint32) *data.Entry {
	if ent, ok := lbbft.Log.Get(eid); ok {
		return ent
	} else {
		lbbft.waitEntry.Wait()
		return lbbft.GetEntry(eid) // unsafe note: 一直等待，直到对应seq的entry到来。
	}
}
