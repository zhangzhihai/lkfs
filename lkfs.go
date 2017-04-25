package lkfs

type lkfs struct {
	filest  *filesystem
	fileidx *FileIdx
}

func New(path string) (fs *lkfs, err error) {
	fi, err := NewFildIdx(path)
	if err != nil {
		return nil, err
	}

	f, err := Newfs(path, fi.BlockId)
	if err != nil {
		return nil, err
	}
	fs = &lkfs{
		filest:  f,
		fileidx: fi,
	}

	return fs, nil
}

//
func (fs *lkfs) Set(body []byte) (int32, error) {

	size := len(body)
	n, fid, err := fs.filest.Write(int32(size), body)

	if err != nil {
		return 0, err
	}
	//return 0, nil
	//写入索引
	id, err := fs.fileidx.Uuid()

	if err != nil {
		return 0, err
	}
	err = fs.fileidx.Write(id, fid, n, int32(size))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (fs *lkfs) Get(id int32) (res []byte, err error) {
	bkid, start, strlen, err := fs.fileidx.Reader(id)
	if err != nil {
		return []byte(""), err
	}

	out, err := fs.filest.Read(bkid, start, strlen)
	return out, err
}

/**
//返回列表
func (fs *lkfs) List(offset int32, limit int32) ([]int32, error) {

}
**/
