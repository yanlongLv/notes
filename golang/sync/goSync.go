package golang

import "sync"

// UserInfo ...
type UserInfo struct {
	Name string
}

var (
	lock     sync.Mutex
	instance *UserInfo
)

func getInstance() (*UserInfo, error) {
	if instance == nil {
		//---Lock
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &UserInfo{
				Name: "fan",
			}
		}
	} //---Unlock()
	return instance, nil
}
