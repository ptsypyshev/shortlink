package models

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const LinkType = "link"

type Link struct {
	ID           int    `json:"id,omitempty" mapstructure:"id"`
	ShortLink    string `json:"short_link,omitempty" mapstructure:"short_link"`
	LongLink     string `json:"long_link,omitempty" mapstructure:"long_link"`
	ClickCounter int    `json:"click_counter,omitempty" mapstructure:"click_counter"`
	OwnerID      int    `json:"owner_id,omitempty" mapstructure:"owner_id"`
	IsActive     bool   `json:"is_active,omitempty" mapstructure:"is_active"`
	//Clickers     []Clicker `json:"clickers,omitempty" mapstructure:"id"`
}

func (l *Link) GetType() string {
	return LinkType
}

func (l *Link) GetList() (lst []interface{}) {
	lst = append(lst, l.ShortLink, l.LongLink, l.ClickCounter, l.OwnerID, l.IsActive)
	return
}

func (l *Link) Set(m map[string]interface{}) error {
	if err := mapstructure.Decode(m, &l); err != nil {
		return err
	}
	//fmt.Println(l)
	return nil
}

func (l *Link) String() string {
	return fmt.Sprintf("{\nID: %d\nShortLink: %s\nLongLink: %s\nClickCounter: %d\nOwnerID: %v\nStatus: %t\n}",
		l.ID, l.ShortLink, l.LongLink, l.ClickCounter, l.OwnerID, l.IsActive)
}
