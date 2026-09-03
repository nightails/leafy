package tui

type state struct {
	devices []device // mounted devices
	media   []medium // media to be copied
	tasks   []task   // tasks to be executed
}

type device struct {
	name       string
	path       string
	mountpoint string
}

type medium struct {
	name            string
	format          string
	root, src, dest string
	copied, total   int64
}

type task struct {
	id          int
	media       medium
	done, total int64
	err         error
}
