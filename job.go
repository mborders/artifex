package queue

type Job struct {
	Run func()
}
