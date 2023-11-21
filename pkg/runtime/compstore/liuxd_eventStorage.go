package compstore

import (
	"github.com/liuxd6825/dapr-components-contrib/liuxd/eventstorage"
)

func (c *ComponentStore) AddEventStorage(name string, item eventstorage.EventStorage) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.eventStorages[name] = item
}

func (c *ComponentStore) GetEventStorage(name string) (eventstorage.EventStorage, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.eventStorages[name]
	return item, ok
}

func (c *ComponentStore) ListEventStorage() map[string]eventstorage.EventStorage {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.eventStorages
}

func (c *ComponentStore) EventStoragesLen() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.eventStorages)
}

func (c *ComponentStore) DeleteEventStorage(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.eventStorages, name)
}
