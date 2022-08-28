package models

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const LinkType = "link"

type Link struct {
	ID           int    `json:"id,omitempty" mapstructure:"id"`
	LongLink     string `json:"long_link,omitempty" mapstructure:"long_link"`
	ClickCounter int    `json:"click_counter,omitempty" mapstructure:"click_counter"`
	OwnerID      int    `json:"owner_id,omitempty" mapstructure:"owner_id"`
	IsActive     bool   `json:"is_active,omitempty" mapstructure:"is_active"`
}

func (l *Link) GetType() string {
	return LinkType
}

func (l *Link) GetList() (lst []interface{}) {
	lst = append(lst, l.LongLink, l.ClickCounter, l.OwnerID, l.IsActive)
	return
}

func (l *Link) Set(m map[string]interface{}) error {
	if err := mapstructure.Decode(m, &l); err != nil {
		return err
	}
	//fmt.Println(l)
	return nil
}

func (l *Link) Get() map[string]interface{} {
	mLinkFields := map[string]interface{}{
		"id":            l.ID,
		"long_link":     l.LongLink,
		"click_counter": l.ClickCounter,
		"owner_id":      l.OwnerID,
		"is_active":     l.IsActive,
	}
	return mLinkFields
}

func (l *Link) String() string {
	return fmt.Sprintf("{\nID: %d\nLongLink: %s\nClickCounter: %d\nOwnerID: %v\nIsActive: %t\n}",
		l.ID, l.LongLink, l.ClickCounter, l.OwnerID, l.IsActive)
}
