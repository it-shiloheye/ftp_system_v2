// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db_access

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FileStatusType string

const (
	FileStatusTypeNew     FileStatusType = "new"
	FileStatusTypeDeleted FileStatusType = "deleted"
	FileStatusTypeUpdated FileStatusType = "updated"
	FileStatusTypeStored  FileStatusType = "stored"
)

func (e *FileStatusType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FileStatusType(s)
	case string:
		*e = FileStatusType(s)
	default:
		return fmt.Errorf("unsupported scan type for FileStatusType: %T", src)
	}
	return nil
}

type NullFileStatusType struct {
	FileStatusType FileStatusType `json:"file_status_type"`
	Valid          bool           `json:"valid"` // Valid is true if FileStatusType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFileStatusType) Scan(value interface{}) error {
	if value == nil {
		ns.FileStatusType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FileStatusType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFileStatusType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FileStatusType), nil
}

type FileTrackerStatus string

const (
	FileTrackerStatusRequested FileTrackerStatus = "requested"
	FileTrackerStatusUploaded  FileTrackerStatus = "uploaded"
	FileTrackerStatusStored    FileTrackerStatus = "stored"
	FileTrackerStatusDeleted   FileTrackerStatus = "deleted"
)

func (e *FileTrackerStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FileTrackerStatus(s)
	case string:
		*e = FileTrackerStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for FileTrackerStatus: %T", src)
	}
	return nil
}

type NullFileTrackerStatus struct {
	FileTrackerStatus FileTrackerStatus `json:"file_tracker_status"`
	Valid             bool              `json:"valid"` // Valid is true if FileTrackerStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFileTrackerStatus) Scan(value interface{}) error {
	if value == nil {
		ns.FileTrackerStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FileTrackerStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFileTrackerStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FileTrackerStatus), nil
}

type PeerRoleType string

const (
	PeerRoleTypeClient  PeerRoleType = "client"
	PeerRoleTypeStorage PeerRoleType = "storage"
	PeerRoleTypeServer  PeerRoleType = "server"
)

func (e *PeerRoleType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PeerRoleType(s)
	case string:
		*e = PeerRoleType(s)
	default:
		return fmt.Errorf("unsupported scan type for PeerRoleType: %T", src)
	}
	return nil
}

type NullPeerRoleType struct {
	PeerRoleType PeerRoleType `json:"peer_role_type"`
	Valid        bool         `json:"valid"` // Valid is true if PeerRoleType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPeerRoleType) Scan(value interface{}) error {
	if value == nil {
		ns.PeerRoleType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PeerRoleType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPeerRoleType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PeerRoleType), nil
}

type FileLog struct {
	ID                int32              `json:"id"`
	PeerID            uuid.UUID          `json:"peer_id"`
	FileHash          string             `json:"file_hash"`
	CurrentFileStatus NullFileStatusType `json:"current_file_status"`
	OldFileStatus     NullFileStatusType `json:"old_file_status"`
	DeltaTime         pgtype.Timestamptz `json:"delta_time"`
}

type FileStorage struct {
	ID               int32              `json:"id"`
	PeerID           uuid.UUID          `json:"peer_id"`
	FileName         string             `json:"file_name"`
	FilePath         string             `json:"file_path"`
	FileType         string             `json:"file_type"`
	FileHash         *string            `json:"file_hash"`
	PrevFileHash     *string            `json:"prev_file_hash"`
	Creation         pgtype.Timestamptz `json:"creation"`
	ModificationDate pgtype.Timestamp   `json:"modification_date"`
	FileState        NullFileStatusType `json:"file_state"`
	FileData         []byte             `json:"file_data"`
}

type PeersTable struct {
	ID        int32            `json:"id"`
	PeerID    pgtype.UUID      `json:"peer_id"`
	IpAddress string           `json:"ip_address"`
	PeerRole  NullPeerRoleType `json:"peer_role"`
	PeerName  *string          `json:"peer_name"`
	Pem       []byte           `json:"pem"`
}
