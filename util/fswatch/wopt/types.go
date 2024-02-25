package wopt

type FileWatchType int32

const (
	WatchInit    FileWatchType = iota
	WatchChanged               //文件修改
	WatchCreated               //文件被创建
	WatchDelete                //文件被删除
)

type FileWatchEvent struct {
	WatchType FileWatchType
	Dir       string
	FileName  string
}
