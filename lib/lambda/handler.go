package lambda

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

type MessageTpl struct {
	VarName string
	VarVal  interface{}
}

type StorageEntry struct {
	Namespace   string
	writeChan   chan<- MessageTpl
	persistence map[string]interface{}
	lambdas     map[string]string
}

var store = map[string]*StorageEntry{}

func NewEntry() *StorageEntry {
	return &StorageEntry{
		persistence: make(map[string]interface{}),
	}
}

func Add(name string, entry *StorageEntry) error {
	entry.lambdas = make(map[string]string)
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

func (se *StorageEntry) SetupPersist(varName string) {
	se.persistence[varName] = ""
}

func (se *StorageEntry) SetWriterChannel(wc chan<- MessageTpl) {
	se.writeChan = wc
}

func (se *StorageEntry) RunLambda(name string) {
	state := createBaseState(se)
	lua.LoadString(state, se.lambdas[name])
	state.Call(0, 0)
	// let gc handle state?
}
