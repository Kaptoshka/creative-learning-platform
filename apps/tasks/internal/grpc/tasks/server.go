package auth

type Tasks interface {
}

type serverAPI struct {
	tasksv1.UnimplementedTasksServer
	tasks Tasks
}
