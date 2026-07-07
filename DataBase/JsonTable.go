package DataBase

import (
	"Server/Helpers/File"
	"encoding/json"
	"errors"
	"sync"
)

type JsonTable[T Indexable] struct {
	mutex     sync.Mutex
	directory string
	path      string
	lastId    uint32
	data      []T
	def       T
}

/* public */

func NewJsonTable[T Indexable](name string, d T) *JsonTable[T] {
	result := &JsonTable[T]{}

	result.directory = File.CurrentDirectory() + "/db"
	result.path = result.directory + "/" + name + ".json"
	result.lastId = 0
	result.data = make([]T, 0)
	result.def = d

	if File.Exists(result.path) {
		bytes := File.ReadBytes(result.path)
		json.Unmarshal(bytes, &result.data)
	} else {
		File.CreateDirectory(result.directory)
		File.WriteText(result.path, "[]")
	}

	result.updateLastId()

	return result
}

func (table *JsonTable[T]) Add(item T) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	table.data = append(table.data, item)

	table.updateLastId()
	table.save()
}

func (table *JsonTable[T]) AddOrSet(item T) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	index := table.findIndex(item)

	if index >= 0 {
		table.data[index] = item
	} else {
		table.data = append(table.data, item)
	}

	table.updateLastId()
	table.save()
}

func (table *JsonTable[T]) Count() uint32 {
	return uint32(len(table.data))
}

func (table *JsonTable[T]) Get(predicate func(T) bool) (T, error) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	for _, item := range table.data {
		if predicate(item) {
			return item, nil
		}
	}

	return table.def, errors.New("Not found")
}

func (table *JsonTable[T]) GetAt(index uint32) (T, error) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	if index >= table.Count() {
		return table.def, errors.New("Not found")
	}

	return table.data[index], nil
}

func (table *JsonTable[T]) NextId(predicate func(T) bool) uint32 {
	if predicate == nil {
		return table.lastId + 1
	}

	for _, item := range table.data {
		if predicate(item) {
			return item.GetId()
		}
	}

	return table.lastId + 1
}

func (table *JsonTable[T]) Remove(predicate func(T) bool) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	if predicate == nil {
		table.data = table.data[:0]
	} else {
		filtered := table.data[:0]

		for _, item := range table.data {
			if !predicate(item) {
				filtered = append(filtered, item)
			}
		}

		table.data = filtered
	}

	table.updateLastId()
	table.save()
}

func (table *JsonTable[T]) RemoveAll() {
	table.Remove(nil)
}

/* private */

func (table *JsonTable[T]) findIndex(item T) int {
	for i, existing := range table.data {
		if existing.GetId() == item.GetId() {
			return i
		}
	}
	return -1
}

func (table *JsonTable[T]) save() {
	data, _ := json.Marshal(table.data)

	File.WriteBytes(table.path, data)
}

func (table *JsonTable[T]) updateLastId() {
	for _, item := range table.data {
		if item.GetId() > table.lastId {
			table.lastId = item.GetId()
		}
	}
}
