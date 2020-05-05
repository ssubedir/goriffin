package queue

// Task - Task interface
type Task interface {
	Run()
	CheckHeartbeat()
}
