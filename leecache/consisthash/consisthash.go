package consisthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int						//虚拟节点倍数
	keys     []int					//环形节点
	HashMap  map[int]string			//真实节点和虚拟节点映射
}

func New(fn Hash,replicas int) *Map {
	m := &Map{
		hash:fn,
		replicas: replicas,
		HashMap: make(map[int]string),
	}
	if fn == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.HashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.HashMap[m.keys[idx%len(m.keys)]]
}
