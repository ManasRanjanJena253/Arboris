package queue

type Job struct {
	InstallationID int64
	Owner          string
	RepoName       string
	PRNumber       int
	CommitID       string
}

type JobQueue struct {
	Queue chan *Job
}

func NewJobQueue(bufferSize int) *JobQueue {
	return &JobQueue{
		Queue: make(chan *Job, bufferSize),
	}
}

func (jobBuffer *JobQueue) JobEnqueue(job *Job) bool {
	select {
	case jobBuffer.Queue <- job:
		return true
	default:
		return false
	}
}
