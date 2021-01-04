package leecache

import pb "leecache/leecache/leecachepb"

type PeerPicker interface {
	PickerPeer(key string) (peer PeerGetter,ok bool)
}

type PeerGetter interface {
	//Get(group string,key string) ([]byte,error)
	Get(in *pb.Request,out *pb.Response) error
}