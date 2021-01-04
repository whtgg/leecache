package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes					int64                         //最大使用内存
	nBytes						int64                       //已使用内存
	ll							*list.List                  //链表
	cache						map[string]*list.Element     //缓存队列元素地址
	OnEvicted					func(key string,value Value) //移除元素 回调
}

type Entry struct {
	key   string
	Value Value
}

type Value interface {
	Len()						int
}

func New(max int64,f func(key string,value Value)) *Cache {
	return &Cache{
		maxBytes:  max,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element,0),
		OnEvicted: f,
	}
}

//从缓存中获得
func (c *Cache) Get(key string) (value Value,ok bool) {
	if ele,ok := c.cache[key];ok {
		c.ll.PushFront(ele)
		kv := ele.Value.(*Entry)
		return kv.Value,true
	}
	return
}

//从缓存中删除
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)         //移除链表元素
		kv := ele.Value.(*Entry) //移除映射表元素
		delete(c.cache,kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.Value.Len()) //更新内存占用
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key,kv.Value)
		}
	}
}


//新增&修改
func (c *Cache) Add(key string,value Value)  {
	//判断映射表是否存在
	if ele,ok := c.cache[key];ok {
		c.ll.PushFront(ele) 											//将元素移到队尾
		kv := ele.Value.(*Entry)
		c.nBytes  += int64(value.Len()) - int64(kv.Value.Len()) 		//更新占用内存
		kv.Value = value
	} else {
		ele := c.ll.PushFront(&Entry{key: key,Value: value}) //添加元素
		c.cache[key] = ele                                   //更新映射表
		c.nBytes += int64(len(key)) + int64(value.Len())     //更新内存
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()												//删除溢出元素
	}
}

func (c *Cache) Len() int{
	return c.ll.Len()
}






