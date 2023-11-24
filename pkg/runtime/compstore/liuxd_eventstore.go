package compstore

import (
	"github.com/liuxd6825/dapr-components-contrib/liuxd/eventstore"
	"strings"
)

func (c *ComponentStore) AddEventStore(name string, item eventstore.EventStore) {
	c.lock.Lock()
	defer c.lock.Unlock()
	compName := strings.ToLower(name)
	c.eventStorages[compName] = item
}

func (c *ComponentStore) GetEventStore(name string) (eventstore.EventStore, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	compName := strings.ToLower(name)
	item, ok := c.eventStorages[compName]
	return item, ok
}

func (c *ComponentStore) ListEventStore() map[string]eventstore.EventStore {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.eventStorages
}

func (c *ComponentStore) EventStoresLen() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.eventStorages)
}

func (c *ComponentStore) DeleteEventStore(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	compName := strings.ToLower(name)
	delete(c.eventStorages, compName)
}
