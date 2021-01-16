module github.com/chinmobi/modlib

go 1.15

require (
	github.com/chinmobi/modlib/evt v0.1.0 // indirect
	github.com/chinmobi/modlib/grpool v0.1.0 // indirect
)

replace (
	github.com/chinmobi/modlib/evt => "./evt"
	github.com/chinmobi/modlib/grpool => "./grpool"
)
