package main

type FileState struct {
	Source string `json:"source"`
	Offset int64  `json:"offset"`
	Inode  uint64 `json:"inode"`
	Device int32  `json:"device"`
}
