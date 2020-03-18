package dailytoggl

import (
	"encoding/base64"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

// TwitterClient defines twitter client's behaviors.
type TwitterClient interface {
	post(string, []byte) error
}

type twitterClient struct {
	*anaconda.TwitterApi
}

func newTwitterClient(conf *config) TwitterClient {
	api := anaconda.NewTwitterApiWithCredentials(
		conf.TwitterAccessToken,
		conf.TwitterAccessTokenSecret,
		conf.TwitterConsumerKey,
		conf.TwitterConsumerSecret,
	)
	return &twitterClient{TwitterApi: api}
}

func (c *twitterClient) post(msg string, img []byte) error {
	encoded := base64.StdEncoding.EncodeToString(img)
	media, err := c.UploadMedia(encoded)
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Add("media_ids", media.MediaIDString)
	_, err = c.PostTweet(msg, v)
	return err
}
