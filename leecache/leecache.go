package leecache

import (
	"errors"
	pb "leecache/leecache/leecachepb"
	"leecache/leecache/singleflight"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte,error)
}

type GetterFunc func(key string) ([]byte,error)

func (f GetterFunc) Get(key string) ([]byte,error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group
}

var (
	mu						sync.RWMutex
	groups =  make(map[string]*Group,0)
)


func NewGroup(name string,cacheBytes int64,getter Getter) *Group {
	if getter  == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(key string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if g,ok := groups[key];ok {
		return g
	}
	return nil
}

func (g *Group) Get(key string) (ByteView,error) {
	if key == "" {
		return ByteView{},errors.New("fail")
	}
	if v,ok := g.mainCache.get(key);ok {
		log.Println("[leeCache] hit")
		return v,nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView,error){
	if g.peers != nil {
		if peer,ok := g.peers.PickerPeer(key);ok {
			if value,err := g.getFromPeer(peer,key);err == nil {
				return value,nil
			} else {
				log.Println("[leeCache] Failed to get from peer", err)
			}
		}
	}
	return g.getLocally(key)
}


func (g *Group) getLocally(key string) (ByteView,error) {
	bytes,err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value,nil
}

func (g *Group) populateCache(key string,value ByteView) {
	g.mainCache.add(key,value)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter,key string) (ByteView,error) {
	req := &pb.Request{
		Group: g.name,
		Key: key,
	}

	res := &pb.Response{}
	err := peer.Get(req,res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value},nil
}

