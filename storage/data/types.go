package data

import "github.com/owncast/owncast/models"

// GetStringSlice will return the string slice value for a key.
func (ds *Store) GetStringSlice(key string) ([]string, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return []string{}, err
	}
	return configEntry.GetStringSlice()
}

// SetStringSlice will set the string slice value for a key.
func (ds *Store) SetStringSlice(key string, value []string) error {
	configEntry := models.ConfigEntry{key, value}
	return ds.Save(configEntry)
}

// GetString will return the string value for a key.
func (ds *Store) GetString(key string) (string, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return "", err
	}
	return configEntry.GetString()
}

// SetString will set the string value for a key.
func (ds *Store) SetString(key string, value string) error {
	configEntry := models.ConfigEntry{key, value}
	return ds.Save(configEntry)
}

// GetNumber will return the numeric value for a key.
func (ds *Store) GetNumber(key string) (float64, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return 0, err
	}
	return configEntry.GetNumber()
}

// SetNumber will set the numeric value for a key.
func (ds *Store) SetNumber(key string, value float64) error {
	configEntry := models.ConfigEntry{Key: key, Value: value}
	return ds.Save(configEntry)
}

// GetBool will return the boolean value for a key.
func (ds *Store) GetBool(key string) (bool, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return false, err
	}
	return configEntry.GetBool()
}

// SetBool will set the boolean value for a key.
func (ds *Store) SetBool(key string, value bool) error {
	configEntry := models.ConfigEntry{Key: key, Value: value}
	return ds.Save(configEntry)
}

// GetStringMap will return the string map value for a key.
func (ds *Store) GetStringMap(key string) (map[string]string, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return map[string]string{}, err
	}
	return configEntry.GetStringMap()
}

// SetStringMap will set the string map value for a key.
func (ds *Store) SetStringMap(key string, value map[string]string) error {
	configEntry := models.ConfigEntry{Key: key, Value: value}
	return ds.Save(configEntry)
}