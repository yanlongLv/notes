package group

import "sync"

//Group ...
type Group struct {
	new  func() interface{}
	objs map[string]interface{}
	sync.RWMutex
}

//NewGroup ..
func NewGroup(new func() interface{}) *Group {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group{
		new:  new,
		objs: make(map[string]interface{}),
	}
}

//Get ...
func (g *Group) Get(key string) interface{} {
	g.RLock()
	obj, ok := g.objs[key]
	if ok {
		g.RUnlock()
		return obj
	}
	g.RUnlock()
	g.Lock()
	defer g.Unlock()
	obj, ok = g.objs[key]
	if ok {
		return obj
	}
	obj = g.new()
	g.objs[key] = obj
	return obj
}

//Reset ...
func (g *Group) Reset(new func() interface{}) {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	g.Lock()
	g.new = new
	g.Unlock()
	g.Clear()
}

//Clear ...
func (g *Group) Clear() {
	g.Lock()
	g.objs = make(map[string]interface{})
	g.Unlock()
}
