package server_config

import (
	"encoding/json"
	"os"
	"time"

	
	ftp_base "github.com/it-shiloheye/ftp_system_v2/lib/base"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
)

var Logger = logging.Logger

var StorageSingleton = NewStorageStruct()

type StorageStruct struct {
	PeerId          pgtype.UUID      `json:"peer_id"`
	IncludeDirs     []string         `json:"include_dirs"`
	ExcludeDirs     []string         `json:"exclude_dirs"`
	IncludeFiles    []string         `json:"include_files"`
	ExcludeRegex    []string         `json:"exclude_regex"`
	IncludeRegex    []string         `json:"include_regex"`
	PollInterval    time.Duration    `json:"poll_interval"`
	SubscribedPeers []SubscribePeers `json:"subscribed_peers"`
}

type SubscribePeers struct {
	PeerUuid               string `json:"peer_uuid"`
	CopyDirectoryStructure bool   `json:"copy_directory_structure"`
	LocalName              string `json:"local_name"`
	PublishOnEdits         string `json:"publish_on_edits"`
	CascadeDeletes         string `json:"cascade_deletes"`
}

func NewStorageStruct() (sts *StorageStruct) {
	sts = &StorageStruct{
		IncludeDirs:     []string{},
		ExcludeDirs:     []string{},
		IncludeFiles:    []string{},
		ExcludeRegex:    []string{},
		IncludeRegex:    []string{},
		PollInterval:    time.Minute * 5,
		SubscribedPeers: []SubscribePeers{},
	}

	return
}

func WriteToDisk(file_path string, storage_struct *StorageStruct) error {
	loc := log_item.Locf(`WriteToDisk(file_path: %s, str_struct *StorageStruct) error`, file_path)

	d, err1 := json.MarshalIndent(storage_struct, " ", "\t")
	if err1 != nil {
		return Logger.LogErr(loc, err1)
	}

	err2 := os.WriteFile(file_path, d, ftp_base.FS_MODE)
	if err2 != nil {
		return Logger.LogErr(loc, err2)
	}
	return nil
}

func ReadFromDisk(file_path string, sts *StorageStruct) (err error) {
	loc := log_item.Locf(`ReadFromDisk(file_path: %s) (sts *StorageStruct, err error)`, file_path)

	f, err1 := os.Open(file_path)
	if err1 != nil {
		err = Logger.LogErr(loc, err1)
		return
	}

	err2 := json.NewDecoder(f).Decode(sts)
	if err2 != nil {
		err = Logger.LogErr(loc, err2)
		return
	}

	return
}
