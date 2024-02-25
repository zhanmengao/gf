package antsqueue

type Option struct {
	NonBlock  bool  //非阻塞
	Size      int   //队列大小
	CloseWait int64 //关闭时等待毫秒，-1为死等
}

var DefaultOption = NewOption()

func NewOption() *Option {
	return &Option{
		NonBlock:  false,
		Size:      5000,
		CloseWait: 1000,
	}
}

func (o *Option) WithNonBlock(nonBlock bool) *Option {
	o.NonBlock = nonBlock
	return o
}

func (o *Option) WithSize(sz int) *Option {
	o.Size = sz
	return o
}

func (o *Option) WithCloseWait(ms int64) *Option {
	o.CloseWait = ms
	return o
}
