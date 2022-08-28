package models

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const ShortLinkType = "shortlink"

type ShortLink struct {
	ID         int    `json:"id,omitempty" mapstructure:"id"`
	Token      string `json:"token,omitempty" mapstructure:"token"`
	LongLinkID int    `json:"long_link_id,omitempty" mapstructure:"long_link_id"`
}

func (s *ShortLink) GetType() string {
	return ShortLinkType
}

func (s *ShortLink) GetList() (lst []interface{}) {
	lst = append(lst, s.Token, s.LongLinkID)
	return
}

func (s *ShortLink) Set(m map[string]interface{}) error {
	if err := mapstructure.Decode(m, &s); err != nil {
		return err
	}
	return nil
}

func (s *ShortLink) Get() map[string]interface{} {
	mShortLinkFields := map[string]interface{}{
		"id":           s.ID,
		"token":        s.Token,
		"long_link_id": s.LongLinkID,
	}
	return mShortLinkFields
}

func (s *ShortLink) String() string {
	return fmt.Sprintf("{\nID: %d\nToken: %s\nLongLinkID: %d\n}",
		s.ID, s.Token, s.LongLinkID)
}
