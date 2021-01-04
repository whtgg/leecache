package leecache

type ByteView struct {					//缓存数据
	b				[]byte
}

func (v ByteView) Len() int { //缓存数据大小
	return len(v.b)
}

func (v *ByteView) ByteSlice() []byte { //返回复制的字节切片
	return cloneBytes(v.b)
}

func (v *ByteView) String() string { //返回字节字符串
	return string(v.b)
}

func cloneBytes(data []byte) []byte {			//返回字节副本
	c := make([]byte,len(data))
	copy(c,data)
	return c
}

