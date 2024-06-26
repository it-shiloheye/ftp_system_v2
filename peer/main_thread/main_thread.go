package mainthread

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	ftp_base "github.com/it-shiloheye/ftp_system_v2/lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	db "github.com/it-shiloheye/ftp_system_v2/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system_v2/lib/db_access/generated"

	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system_v2/peer/config"
	db_helpers "github.com/it-shiloheye/ftp_system_v2/peer/main_thread/db_access"
	"github.com/jackc/pgx/v5/pgtype"
)

var Logger = logging.Logger
var ServerConfig = server_config.ServerConfig

func Loop(ctx ftp_context.Context) error {
	loc := log_item.Loc(`mainthread.Loop(ctx ftp_context.Context) error`)
	defer ctx.Finished()
	StorageStruct := server_config.NewStorageStruct()

	tc := time.NewTicker(time.Minute * 5)
	uniq_to_upload := map[string]*FilesListItem{}
	uniq_to_download := map[string]*db_access.GetFilesRow{}
	set_up_storagestruct(StorageStruct)
	for ok := true; ok; {

		Logger.Logf(loc, "new loop: %s", fmt.Sprintln(time.Now()))

		tmp_fileslist, err1 := WalkDir(ctx, StorageStruct)
		if err1 != nil {
			Logger.LogErr(loc, err1)
			continue
		}
		for _, file_item := range tmp_fileslist {
			uniq_to_upload[file_item.Path] = file_item
		}
		db_tmp_fileslist, err2 := db_helpers.GetFiles(ctx, StorageStruct)
		if err2 != nil {
			Logger.LogErr(loc, err2)
			continue
		}

		log.Println("reaching here")
		for _, file_uploaded := range db_tmp_fileslist {
			uniq_to_download[file_uploaded.FilePath] = file_uploaded

			tmp, ok := uniq_to_upload[file_uploaded.FilePath]
			if !ok {
				if file_uploaded.FileState.Valid && db_access.FileStatusTypeDeleted == file_uploaded.FileState.FileStatusType {
					delete(uniq_to_download, file_uploaded.FilePath)
				}
				continue

			}

			if tmp.FD.ModTime() == file_uploaded.ModificationDate.Time {
				delete(uniq_to_download, file_uploaded.FilePath)
				continue
			}
		}

		insert_file_rows := []*db_access.InsertFileRow{}

		for file_path, FD := range uniq_to_upload {
			d, err4 := os.ReadFile(file_path)
			if err4 != nil {
				Logger.LogErr(loc, err4)
				continue
			}

			name_01 := strings.Split(file_path, string(os.PathSeparator))
			name_0l := len(name_01) - 1
			name := name_01[name_0l]
			conn := db.DBPool.GetConn()
			insert_file_row, err7 := db_helpers.DB.InsertFile(ctx.Add(), conn, &db_access.InsertFileParams{
				PeerID:           ServerConfig.PeerId.Bytes,
				FilePath:         file_path,
				FileType:         filepath.Ext(file_path),
				FileName:         name,
				ModificationDate: pgtype.Timestamp{Time: FD.FD.ModTime(), Valid: true},
				FileState:        db_access.NullFileStatusType{FileStatusType: db_access.FileStatusTypeNew},
				FileData:         d,
			})
			db.DBPool.Return(conn)
			if err7 != nil {
				Logger.LogErr(loc, err7)
				continue
			}

			insert_file_rows = append(insert_file_rows, insert_file_row)
		}

		select {
		case <-ctx.Done():
			return Logger.LogErr(loc, ctx.Err())
		case <-tc.C:

			set_up_storagestruct(StorageStruct)
		}
	}

	return nil
}

func set_up_storagestruct(StorageStruct *server_config.StorageStruct) {
	config_filepath := ServerConfig.StorageDirectory + "/config.json"
	loc := log_item.Locf(`set_up_storagestruct := func(StorageStruct :%s) :`, config_filepath)
	f, err1 := os.Open(config_filepath)
	if err1 == nil {
		defer f.Close()
		err2 := json.NewDecoder(f).Decode(StorageStruct)
		if err2 != nil {
			log.Fatalln(Logger.LogErr(loc, err2))
		}

		if len(StorageStruct.IncludeDirs) < 1 {
			log.Printf(`please add at least one directory to included_dirs in: %s\n`, config_filepath)
			<-time.After(time.Minute)
			set_up_storagestruct(StorageStruct)
			return
		}

		return
	}

	if !errors.Is(err1, os.ErrNotExist) {
		log.Fatalln(Logger.LogErr(loc, err1))
	}

	d, err2 := json.MarshalIndent(StorageStruct, " ", "\t")
	if err2 != nil {
		log.Fatalln(Logger.LogErr(loc, err2))
	}
	err3 := os.WriteFile(config_filepath, d, ftp_base.FS_MODE)
	if err3 != nil {
		log.Fatalln(Logger.LogErr(loc, err3))
	}

	log.Printf("please fill: %s appropriately\n", config_filepath)
	<-time.After(time.Minute)
	set_up_storagestruct(StorageStruct)
}

type FilesListItem struct {
	Path string
	FD   os.FileInfo
	*os.File
}

func (fsi *FilesListItem) Reopen() error {
	loc := log_item.Locf(`func (fsi *FilesListItem) Reopen(%s) error `, fsi.Path)
	var err5, err6 error
	fsi.File, err5 = os.OpenFile(fsi.Path, os.O_RDONLY, ftp_base.FS_MODE)
	if err5 != nil {
		return Logger.LogErr(loc, err5)
	}

	fsi.FD, err6 = fsi.File.Stat()
	if err6 != nil {
		return Logger.LogErr(loc, err6)
	}

	return nil
}

func WalkDir(ctx ftp_context.Context, storage_struct *server_config.StorageStruct) (files_list []*FilesListItem, err error) {
	loc := log_item.Loc(`func WalkDir(ctx ftp_context.Context, storage_struct *server_config.StorageStruct)(files_list []FilesListItem, err  error)`)
	files_list = []*FilesListItem{}

	uniq := map[string]bool{}
	for _, dir_path := range storage_struct.IncludeDirs {
		err1 := filepath.WalkDir(dir_path, func(path string, d fs.DirEntry, err2 error) error {
			loc := log_item.Locf(`err1 := filepath.WalkDir(dir_path: %s,func(path %s, d fs.DirEntry, err2: %s) error `, dir_path, path, fmt.Sprint(err2))
			if err2 != nil {
				return Logger.LogErr(loc, err2)
			}

			if _, ok := uniq[path]; ok || d.IsDir() {
				return nil
			}
			uniq[path] = true

			for _, invalid_dir := range storage_struct.ExcludeDirs {
				if strings.Contains(path, invalid_dir) {
					return nil
				}
			}
			for _, invalid_regex := range storage_struct.ExcludeRegex {
				if not_ok, err3 := regexp.MatchString(invalid_regex, path); not_ok || err3 != nil {
					if err3 != nil {
						Logger.LogErr(loc, &log_item.LogItem{
							After:     fmt.Sprintf(`not_ok, err3 := regexp.MatchString(invalid_regex: %s ,path: %s);not_ok || err3 != nil `, invalid_regex, path),
							Level:     log_item.LogLevelError02,
							CallStack: []error{err3},
						})
					}
					return nil
				}
			}

			add := false

			if len(storage_struct.IncludeRegex) > 0 {
				for _, fpath := range storage_struct.IncludeRegex {
					ok, err4 := regexp.MatchString(fpath, path)
					if ok {
						add = true
						break
					}
					if err4 != nil {
						Logger.LogErr(loc, err4)
					}
				}

			} else {
				add = true
			}
			if add {

				tmp := &FilesListItem{
					Path: path,
				}

				var err5, err6 error
				tmp.File, err5 = os.OpenFile(path, os.O_RDONLY, ftp_base.FS_MODE)
				if err5 != nil {
					return Logger.LogErr(loc, err5)
				}

				tmp.FD, err6 = tmp.File.Stat()
				if err6 != nil {
					return Logger.LogErr(loc, err6)
				}
				tmp.Close()
				files_list = append(files_list, tmp)
			}
			return nil
		})

		if err1 != nil {
			err = Logger.LogErr(loc, err1)

			return
		}
	}

	for _, fpath := range storage_struct.IncludeFiles {
		if _, ok := uniq[fpath]; ok {
			continue
		}
		uniq[fpath] = true

		tmp := &FilesListItem{
			Path: fpath,
		}
		err1 := tmp.Reopen()
		if err1 != nil {
			Logger.LogErr(loc, err1)
			continue
		}
		tmp.Close()
		files_list = append(files_list, tmp)

	}

	return
}
