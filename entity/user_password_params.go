package entity

type UserPasswordParams struct {
	TableTemplateCols
	RandSalt string `json:"rand_salt,omitempty"`
	Time     uint32 `json:"time,omitempty"`
	Memory   uint32 `json:"memory,omitempty"`
	Threads  uint8  `json:"threads,omitempty"`
	KeyLen   uint32 `json:"keyLen,omitempty"`
	UsersID  uint64 `json:"users_id,omitempty"`
}
