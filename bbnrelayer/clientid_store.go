package bbnrelayer

import (
	"fmt"

	"github.com/babylonchain/babylon-relayer/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// setClientID sets the clientID for the IBC light client of a Cosmos zone
// so that when restarting the relayer, it does not need to create another
// IBC light client again.
// key: chainID
// value: client ID of the given chain on Babylon.
func (r *Relayer) setClientID(chainID, clientID string) error {
	dbPath := config.GetDBPath(r.homePath)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return fmt.Errorf("error opening LevelDB at %s: %w", dbPath, err)
	}
	defer db.Close() // Ensure the database is closed when the function exits

	if err := db.Put([]byte(chainID), []byte(clientID), &opt.WriteOptions{Sync: true}); err != nil {
		return fmt.Errorf("error writing to LevelDB at %s: %w", dbPath, err)
	}

	return nil
}

// getClientID retrieves the clientID for the IBC light client of a Cosmos zone
// from the LevelDB database.
// Returns an empty string and nil error if the client ID is not found.
func (r *Relayer) getClientID(chainID string) (string, error) {
	dbPath := config.GetDBPath(r.homePath)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return "", fmt.Errorf("error opening LevelDB at %s: %w", dbPath, err)
	}
	defer db.Close() // Ensure the database is closed when the function exits

	clientID, err := db.Get([]byte(chainID), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", nil // Return nil error if the client ID is not found
		}
		return "", fmt.Errorf("error reading from LevelDB at %s: %w", dbPath, err)
	}

	return string(clientID), nil
}
