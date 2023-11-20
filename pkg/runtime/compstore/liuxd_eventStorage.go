package compstore

import "github.com/liuxd6825/components-contrib/liuxd/eventstorage"

func (c *ComponentStore) AddEventStorage(name string, store eventstorage.EventStorage) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.eventStorages[name] = store
}

func (c *ComponentStore) GetEventStorage(name string) (eventstorage.EventStorage, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	store, ok := c.eventStorages[name]
	return store, ok
}
