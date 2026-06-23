package queue

type Job struct {
	InstallationID string
	Owner          string
	RepoName       string
	PRNumber       int
	CommitID       string
}

type jobQueue struct {
	queue chan *Job
}

func NewJobQueue(bufferSize int) *jobQueue {
	return &jobQueue{
		queue: make(chan *Job, bufferSize),
	}
}

func (jobBuffer *jobQueue) JobEnqueue(job *Job) bool {
	select {
	case jobBuffer.queue <- job:
		return true
	default:
		return false
	}
}
