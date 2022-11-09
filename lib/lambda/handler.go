package lambda

import (
	"fmt"
	"time"

	"github.com/Shopify/go-lua"
	"github.com/sirupsen/logrus"
)

type MessageTpl struct {
	VarName string
	VarVal  interface{}
}

type StorageEntry struct {
	session     string
	lastUsed    time.Time
	writeChan   chan<- MessageTpl
	persistence map[string]interface{}
	lambdas     map[string]string
}

var store = map[string]*StorageEntry{}

func NewEntry() *StorageEntry {
	return &StorageEntry{
		lastUsed:    time.Now(),
		persistence: make(map[string]interface{}),
		lambdas:     make(map[string]string),
	}
}

func Add(name string, entry *StorageEntry) error {
	entry.session = name
	store[name] = entry
	return nil
}

func Get(name string) (*StorageEntry, error) {
	se, ok := store[name]
	if !ok {
		return nil, fmt.Errorf("missing")
	}

	return se, nil
}

func (se *StorageEntry) AddLambda(name, code string) {
	se.lambdas[name] = code
}

func (se *StorageEntry) SetupPersist(varName string, varVal string) {
	se.persistence[varName] = varVal
}

func (se *StorageEntry) SetWriterChannel(wc chan<- MessageTpl) {
	se.writeChan = wc
}

func (se *StorageEntry) RunLambda(name string, args map[string]string) {
	defer func() {
		err := recover()
		if err != nil {
			se.writeChan <- MessageTpl{VarName: "error", VarVal: "fatal error in lua code execution"}
		}
	}()

	state := createBaseState(se)

	// make sure we keep track usage so the session does not get deleted by cleanup worker
	se.lastUsed = time.Now()

	// inject argument map
	state.NewTable()
	tableIndex := state.Top()

	for key, val := range args {
		state.PushString(key)
		state.PushString(val)
		state.SetTable(tableIndex)
	}

	state.SetGlobal("args")

	lua.LoadString(state, se.lambdas[name])
	state.Call(0, 0)
	// let gc handle state?
}

func (se *StorageEntry) ToTemplate() *Template {
	strArgMap := make(map[string]string)
	for varName, varVal := range se.persistence {
		strArgMap[varName] = fmt.Sprintf("%v", varVal)
	}

	return &Template{
		Arguments: strArgMap,
		Lambdas:   se.lambdas,
	}
}

func CleanupWorker(doneChan <-chan bool) {
	for {
		select {
		case <-doneChan:
			break
		case <-time.After(time.Minute):
			for key, se := range store {
				if se.lastUsed.Before(time.Now().Add(-(time.Hour * 12))) {
					delete(store, key)
					logrus.WithField("namespace", key).Info("deleted session")
				}
			}
		}
	}
}
