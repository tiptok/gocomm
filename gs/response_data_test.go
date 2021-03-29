package gs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMapData(t *testing.T) {
	m := NewMapData()
	m.AddFiled("user.id", 1)
	m.AddFiled("user.name", "tip")
	m.AddFiled("user.sex", true)

	m.AddFiled("address.lon", 59.2156461)
	m.AddFiled("address.lat", 23.1245648)
	m.AddFiled("phone", "18860183050")

	notification := m.GetFiledMap("notification")
	notification["title"] = "xxx"
	notification["body"] = "body"
	m.SetFieldMap(notification, "url", "http://")
	m.SetFieldMap(notification, "options", nil)

	assert.Equal(t, 1, m.MustFindField("user.id"))
	assert.Equal(t, "tip", m.MustFindField("user.name"))
	assert.Equal(t, true, m.MustFindField("user.sex"))
	assert.Nil(t, m.MustFindField("user.1sex1"))
	assert.Equal(t, "xxx", m.MustFindField("notification.title"))
}
