package compstore

import "github.com/liuxd6825/components-contrib/liuxd/applog"

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
