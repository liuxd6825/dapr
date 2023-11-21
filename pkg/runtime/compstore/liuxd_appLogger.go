package compstore

import "github.com/liuxd6825/dapr-components-contrib/liuxd/applog"

func (c *ComponentStore) AddAppLogger(name string, appLog applog.Logger) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.appLoggers[name] = appLog
}

func (c *ComponentStore) GetAppLogger(name string) (applog.Logger, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	appLog, ok := c.appLoggers[name]
	return appLog, ok
}

func (c *ComponentStore) ListAppLogger() map[string]applog.Logger {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.appLoggers
}

func (c *ComponentStore) AppLoggersLen() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.appLoggers)
}

func (c *ComponentStore) DeleteAppLogger(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.appLoggers, name)
}
