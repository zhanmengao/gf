package fclient

import (
	"context"
	"fmt"
	framework "github.com/zhanmengao/gf/context"
	"github.com/zhanmengao/gf/proto/go/gerror"
	"github.com/zhanmengao/gf/wscli/conn/conntyp"
)

func (c *Client) GetUID() string {
	if c.Sess == nil {
		return ""
	}
	return c.Sess.GetUID()
}

func (c *Client) SetUID(ctx context.Context, uid string) (ret context.Context, err error) {
	if c.Sess, ret, err = framework.SessPool.NewUserSession(ctx, uid); err != nil {
		return
	}
	if c.Sess == nil {
		err = gerror.ErrServerInternalUnknown().Format("sess is nil")
		return
	}
	c.Sess.SetSessValue(framework.SessionKeyClient, c)
	return
}

func (c *Client) SetAesKey(aesKey string) {
	c.Sess.SetSessValue(framework.SessionKeyAes, aesKey)
}

func (c *Client) GetAesKey() string {
	if c.Sess == nil {
		return ""
	}
	if iKey := c.Sess.GetSessValue(framework.SessionKeyAes); iKey != nil {
		return iKey.(string)
	}
	return ""
}

func (c *Client) GetRouteAddr(srvName string) (addr string, ok bool) {
	iAddr := c.Sess.GetSessValue(fmt.Sprintf(framework.SessionKeyRouteAddr, srvName))
	if iAddr != nil {
		addr, ok = iAddr.(string)
	}
	return
}

func (c *Client) SetRouteAddr(srvName, addr string) {
	c.Sess.SetSessValue(fmt.Sprintf(framework.SessionKeyRouteAddr, srvName), addr)
}

func (c *Client) ClearRouteAddr(srvName string) {
	c.SetRouteAddr(srvName, "")
}

func (c *Client) GetConnect() conntyp.IConnect {
	return c.conn
}

func (c *Client) GetLastSeq() int32 {
	return c.lastSeq
}
