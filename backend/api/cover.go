package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
	"github.com/go-chi/chi"
)

type cover struct {
	Cover string `json:"cover"`
}

func setCoverURL(cfg *Cfg, c *store.Channel) {
	if c.Cover != "" {
		c.Cover = strings.Join([]string{cfg.Host, "files", c.Alias, fs.CoverDirName, c.Cover}, "/")
	}
}

func (cfg *Cfg) uploadCover(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	cid := r.Context().Value(CID).(int64)

	short, err := cfg.Store.SwapCIDForAlias(cid)
	if err != nil {
		return err
	}

	holder, err := cfg.FS.SaveCover(short, header.Filename)
	if err != nil {
		return err
	}
	defer holder.Close()

	c, err := cfg.Store.ChannelInfo(cid)
	if err != nil {
		return err
	}

	c.Cover = header.Filename
	if err = cfg.Store.UpdateChannel(c); err != nil {
		return err
	}

	setCoverURL(cfg, c)

	output := cover{
		Cover: c.Cover,
	}

	io.Copy(holder, file)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(output); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) deleteCover(w http.ResponseWriter, r *http.Request) error {
	cover := chi.URLParam(r, "cover")

	cid := r.Context().Value(CID).(int64)

	channel, err := cfg.Store.ChannelInfo(cid)
	if err != nil {
		return err
	}

	setCoverURL(cfg, channel)

	channel.Cover = ""
	err = cfg.Store.UpdateChannel(channel)
	if err != nil {
		return err
	}

	err = cfg.FS.RemoveCover(channel.Alias, cover)
	if err != nil {
		return err
	}

	return nil
}
