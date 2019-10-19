package zookeeper

type RetryStrategy interface {
	Retry(func() error) error
}

type discontinue struct {}

func (discontinue) Retry(func() error) error {

}