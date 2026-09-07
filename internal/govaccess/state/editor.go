package state

import "sync"

// StateEditor is the private detached editor interface used by finalized
// ContextBoundChange.ApplyTo. Only a state transaction may obtain one.
type StateEditor interface {
	PutResource(uid string, version string, payload []byte) error
	DeleteResource(uid string) error
	GetResource(uid string) (version string, payload []byte, ok bool)
	ShadowDigest() []byte
	IndexDigest() []byte
}

type memoryEditor struct {
	mu          sync.Mutex
	resources   map[string]resourceEntry
	shadow      []byte
	index       []byte
	deltaShadow []byte
	deltaIndex  []byte
}

type resourceEntry struct {
	version string
	payload []byte
}

func newMemoryEditor(root map[string]resourceEntry) *memoryEditor {
	cp := make(map[string]resourceEntry, len(root))
	for k, v := range root {
		cp[k] = resourceEntry{version: v.version, payload: copyBytes(v.payload)}
	}
	return &memoryEditor{resources: cp, shadow: []byte("shadow-empty"), index: []byte("index-empty")}
}

func (e *memoryEditor) PutResource(uid string, version string, payload []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if uid == "" || version == "" {
		return errInvalid("editor_put_invalid")
	}
	e.resources[uid] = resourceEntry{version: version, payload: copyBytes(payload)}
	e.deltaShadow = append(e.deltaShadow, []byte(uid+":"+version)...)
	e.deltaIndex = append(e.deltaIndex, []byte(uid)...)
	e.shadow = append([]byte(nil), e.deltaShadow...)
	e.index = append([]byte(nil), e.deltaIndex...)
	return nil
}

func (e *memoryEditor) DeleteResource(uid string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.resources, uid)
	e.deltaShadow = append(e.deltaShadow, []byte("del:"+uid)...)
	e.deltaIndex = append(e.deltaIndex, []byte("del:"+uid)...)
	e.shadow = append([]byte(nil), e.deltaShadow...)
	e.index = append([]byte(nil), e.deltaIndex...)
	return nil
}

func (e *memoryEditor) GetResource(uid string) (string, []byte, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	r, ok := e.resources[uid]
	if !ok {
		return "", nil, false
	}
	return r.version, copyBytes(r.payload), true
}

func (e *memoryEditor) ShadowDigest() []byte {
	e.mu.Lock()
	defer e.mu.Unlock()
	return copyBytes(e.shadow)
}

func (e *memoryEditor) IndexDigest() []byte {
	e.mu.Lock()
	defer e.mu.Unlock()
	return copyBytes(e.index)
}

func (e *memoryEditor) snapshot() map[string]resourceEntry {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make(map[string]resourceEntry, len(e.resources))
	for k, v := range e.resources {
		out[k] = resourceEntry{version: v.version, payload: copyBytes(v.payload)}
	}
	return out
}

func (e *memoryEditor) deltaEmpty() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.deltaShadow) == 0 && len(e.deltaIndex) == 0
}
