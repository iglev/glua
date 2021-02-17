package stdlib

import "github.com/iglev/glua/api"

var coFuncs = map[string]api.GoFunction{
	"create":      coCreate,
	"resume":      coResume,
	"yield":       coYield,
	"status":      coStatus,
	"isyieldable": coYieldable,
	"running":     coRunning,
	"wrap":        coWrap,
}

func OpenCoroutineLib(ls api.LuaState) int {
	ls.NewLib(coFuncs)
	return 1
}

// coroutine.create (f)
func coCreate(ls api.LuaState) int {
	ls.CheckType(1, api.LUA_TFUNCTION)
	ls2 := ls.NewThread()
	ls.PushValue(1)  /* move function to top */
	ls.XMove(ls2, 1) /* move function from ls to ls2 */
	return 1
}

// coroutine.resume (co [, val1, ···])
func coResume(ls api.LuaState) int {
	co := ls.ToThread(1)
	ls.ArgCheck(co != nil, 1, "thread expected")

	if r := _auxResume(ls, co, ls.GetTop()-1); r < 0 {
		ls.PushBoolean(false)
		ls.Insert(-2)
		return 2 /* return false + error message */
	} else {
		ls.PushBoolean(true)
		ls.Insert(-(r + 1))
		return r + 1 /* return true + 'resume' returns */
	}
}

func _auxResume(ls, co api.LuaState, narg int) int {
	if !ls.CheckStack(narg) {
		ls.PushString("too many arguments to resume")
		return -1 /* error flag */
	}
	if co.Status() == api.LUA_OK && co.GetTop() == 0 {
		ls.PushString("cannot resume dead coroutine")
		return -1 /* error flag */
	}
	ls.XMove(co, narg)
	status := co.Resume(ls, narg)
	if status == api.LUA_OK || status == api.LUA_YIELD {
		nres := co.GetTop()
		if !ls.CheckStack(nres + 1) {
			co.Pop(nres) /* remove results anyway */
			ls.PushString("too many results to resume")
			return -1 /* error flag */
		}
		co.XMove(ls, nres) /* move yielded values */
		return nres
	} else {
		co.XMove(ls, 1) /* move error message */
		return -1       /* error flag */
	}
}

// coroutine.yield (···)
func coYield(ls api.LuaState) int {
	return ls.Yield(ls.GetTop())
}

// coroutine.status (co)
func coStatus(ls api.LuaState) int {
	co := ls.ToThread(1)
	ls.ArgCheck(co != nil, 1, "thread expected")
	if ls == co {
		ls.PushString("running")
	} else {
		switch co.Status() {
		case api.LUA_YIELD:
			ls.PushString("suspended")
		case api.LUA_OK:
			if co.GetStack() { /* does it have frames? */
				ls.PushString("normal") /* it is running */
			} else if co.GetTop() == 0 {
				ls.PushString("dead")
			} else {
				ls.PushString("suspended")
			}
		default: /* some error occurred */
			ls.PushString("dead")
		}
	}

	return 1
}

// coroutine.isyieldable ()
func coYieldable(ls api.LuaState) int {
	ls.PushBoolean(ls.IsYieldable())
	return 1
}

// coroutine.running ()
func coRunning(ls api.LuaState) int {
	isMain := ls.PushThread()
	ls.PushBoolean(isMain)
	return 2
}

// coroutine.wrap (f)
func coWrap(ls api.LuaState) int {
	panic("todo: coWrap!")
}
