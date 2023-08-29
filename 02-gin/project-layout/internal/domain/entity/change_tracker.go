// change tracker 它的作用是标记哪些数据对象变更了，在持久化的时候有选择性的持久化

package entity

type ChangeTracker struct {
	isDirty bool
}

func (c *ChangeTracker) IsDirty() bool {
	return c.isDirty
}

func (c *ChangeTracker) change() {
	c.isDirty = true
}
