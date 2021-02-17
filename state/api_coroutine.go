package state

import "github.com/iglev/glua/api"

// NewThread - lua_newthread
func (self *luaState) NewThread() api.LuaState {
	t := &luaState{registry: self.registry}
	t.pushLuaStack(newLuaStack(api.LUA_MINSTACK, t))
	self.stack.push(t)
	return t
}

// Resume - lua_resume
func (self *luaState) Resume(from api.LuaState, nArgs int) int {
	lsFrom := from.(*luaState)
	if lsFrom.coChan == nil {
		lsFrom.coChan = make(chan int)
	}

	if self.coChan == nil {
		// start coroutine
		self.coChan = make(chan int)
		self.coCaller = lsFrom
		go func() {
			self.coStatus = self.PCall(nArgs, -1, 0)
			lsFrom.coChan <- 1
		}()
	} else {
		// resume coroutine
		if self.coStatus != api.LUA_YIELD { // todo
			self.stack.push("cannot resume non-suspended coroutine")
			return api.LUA_ERRRUN
		}
		self.coStatus = api.LUA_OK
		self.coChan <- 1
	}

	<-lsFrom.coChan // wait coroutine to finish or yield
	return self.coStatus
}

// Yield - lua_yield
func (self *luaState) Yield(nResults int) int {
	if self.coCaller == nil { // todo
		panic("attempt to yield from outside a coroutine")
	}
	self.coStatus = api.LUA_YIELD
	self.coCaller.coChan <- 1
	<-self.coChan
	return self.GetTop()
}

// IsYieldable - lua_isyieldable
func (self *luaState) IsYieldable() bool {
	if self.isMainThread() {
		return false
	}
	return self.coStatus != api.LUA_YIELD // todo
}

// Status - lua_status
func (self *luaState) Status() int {
	return self.coStatus
}

// debug
func (self *luaState) GetStack() bool {
	return self.stack.prev != nil
}
