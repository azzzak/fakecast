package api

import (
	"encoding/xml"
	"net/http"

	"github.com/azzzak/fakecast/feed"
	"github.com/go-chi/chi"
)

func (cfg *Cfg) genFeed(w http.ResponseWriter, r *http.Request) error {
	c := chi.URLParam(r, "channel")

	cid, err := cfg.Store.SwapAliasForCID(c)
	if err != nil {
		return err
	}

	channel, err := cfg.Store.ChannelInfo(cid)
	if err != nil {
		return err
	}

	setCoverURL(cfg, channel)

	podcasts, err := cfg.Store.ListFullPodcastsFrom(cid)
	if err != nil {
		return err
	}

	podcasts = checkPodcasts(cfg, channel.Alias, podcasts)
	rss := feed.GenerateFeed(channel, podcasts, cfg.Host)

	w.Write([]byte(xml.Header))

	encoder := xml.NewEncoder(w)
	if err := encoder.Encode(rss); err != nil {
		return err
	}

	return nil
}
