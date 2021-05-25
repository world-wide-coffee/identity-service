package memdb

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"io"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"

	"github.com/world-wide-coffee/identity-service/persistence"
)

const (
	writeMode bool = true
	readMode  bool = false

	tableIdentity              string = "identity"
	tableIdentityIndexID       string = "id"
	tableIdentityIndexUsername string = "username"
	tableIdentityFieldID       string = "ID"
	tableIdentityFieldUsername string = "Username"
)

// Storage is a thread safe in-mem storage
type Storage struct {
	db *memdb.MemDB
}

var _ persistence.Persister = &Storage{}

// NewPersister will create an in-mem storage
func NewPersister() *Storage {
	schema := memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableIdentity: {
				Name: tableIdentity,
				Indexes: map[string]*memdb.IndexSchema{
					tableIdentityIndexID: {
						Name:    tableIdentityIndexID,
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: tableIdentityFieldID},
					},
					tableIdentityIndexUsername: {
						Name:    tableIdentityIndexUsername,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: tableIdentityFieldUsername},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(&schema)
	if err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

const saltLength = 64

type identity struct {
	ID       string
	Username string
	Password string
	Salt     []byte
}

// PersistIdentity ...
func (s *Storage) PersistIdentity(ctx context.Context, pid persistence.Identity) error {
	ident := identity{ID: pid.ID.String(), Username: pid.Username}

	ident.Salt = make([]byte, saltLength)
	if _, err := rand.Read(ident.Salt); err != nil {
		return persistence.NewErrorInternal("failed to create salt: %w", err)
	}

	ident.Password = hash(pid.Password, ident.Salt)

	txn := s.db.Txn(writeMode)
	defer txn.Commit()

	if err := txn.Insert(tableIdentity, ident); err != nil {
		txn.Abort()
		return persistence.NewErrorInternal("failed to insert identity: %w", err)
	}

	return nil
}

// ChangePassword ...
func (s *Storage) ChangePassword(ctx context.Context, id uuid.UUID, password string) error {
	ident, err := s.identityByID(ctx, id)
	if err != nil {
		return err
	}

	ident.Password = hash(password, ident.Salt)

	txn := s.db.Txn(writeMode)
	defer txn.Commit()

	if err := txn.Insert(tableIdentity, ident); err != nil {
		txn.Abort()
		return persistence.NewErrorInternal("failed to insert identity: %w", err)
	}

	return nil
}

// IdentityByUsername ...
func (s *Storage) IdentityByUsername(ctx context.Context, username string) (_ persistence.Identity, err error) {
	ident, err := s.identityByUsername(ctx, username)
	if err != nil {
		return
	}

	return persistence.Identity{ID: uuid.MustParse(ident.ID), Username: ident.Username}, err
}

// IdentityByIDAndPassword ...
func (s *Storage) IdentityByIDAndPassword(ctx context.Context, id uuid.UUID, password string) (_ persistence.Identity, err error) {
	ident, err := s.identityByID(ctx, id)
	if err != nil {
		return
	}

	err = s.comparePassword(ctx, ident, password)
	if err != nil {
		return
	}

	return persistence.Identity{ID: uuid.MustParse(ident.ID), Username: ident.Username}, nil
}

// IdentityByUsernameAndPassword ...
func (s *Storage) IdentityByUsernameAndPassword(ctx context.Context, username, password string) (_ persistence.Identity, err error) {
	ident, err := s.identityByUsername(ctx, username)
	if err != nil {
		return
	}

	err = s.comparePassword(ctx, ident, password)
	if err != nil {
		return
	}

	return persistence.Identity{ID: uuid.MustParse(ident.ID), Username: ident.Username}, nil
}

func hash(password string, salt []byte) string {
	hash := sha512.New()
	io.WriteString(hash, password) //nolint: errcheck
	return string(hash.Sum(salt))
}

func (s *Storage) comparePassword(ctx context.Context, storedIdentity identity, userSentPassword string) (err error) {
	hashedPassword := hash(userSentPassword, storedIdentity.Salt)

	if storedIdentity.Password != hashedPassword {
		err = persistence.ErrorWrongPassword{}
		return
	}

	return
}

// identityByUsername ...
func (s *Storage) identityByUsername(ctx context.Context, username string) (ident identity, err error) {
	txn := s.db.Txn(readMode)

	v, err := txn.First(tableIdentity, tableIdentityIndexUsername, username)
	if err != nil {
		err = persistence.NewErrorInternal("failed to find identity for username: %s, err: %w", username, err)
		return
	} else if v == nil {
		err = persistence.ErrorWrongUsername{}
		return
	}

	ident = v.(identity)
	return
}

// identityByID ...
func (s *Storage) identityByID(ctx context.Context, id uuid.UUID) (ident identity, err error) {
	txn := s.db.Txn(readMode)

	v, err := txn.First(tableIdentity, tableIdentityIndexID, id.String())
	if err != nil {
		err = persistence.NewErrorInternal("failed to find identity for id: %s, err: %w", id, err)
		return
	} else if v == nil {
		err = persistence.ErrorNotFound{}
		return
	}

	ident = v.(identity)
	return
}
