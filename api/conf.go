package api

import (
	"github.com/h2object/rpc"
	"path"
)

func (h2o *Client) SetConfig(l Logger, auth Auth, section string, ret interface{}) error {
	URL := rpc.BuildHttpURL(h2o.addr, path.Join(section, ".conf"), nil)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	if err := h2o.conn.PutJson(l, URL, ret, nil); err != nil {
		return err
	}
	return nil
}

func (h2o *Client) GetConfig(l Logger, auth Auth, section string, ret interface{}) error {
	URL := rpc.BuildHttpURL(h2o.addr, path.Join(section, ".conf"), nil)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	if err := h2o.conn.Get(l, URL, ret); err != nil {
		return err
	}
	return nil
}


