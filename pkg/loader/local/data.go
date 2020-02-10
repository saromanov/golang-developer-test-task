package local

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	json "github.com/pquerna/ffjson/ffjson"
	"github.com/saromanov/golang-developer-test-task/pkg/loader"
	"github.com/saromanov/golang-developer-test-task/pkg/models"
)

// Local defines loading data localy
type Local struct {
	path string
}

// New provides initialization of the local storage
func New(path string) (loader.Loader, error) {
	if path == "" {
		return nil, fmt.Errorf("path is not defined")
	}
	return &Local{
		path: path,
	}, nil
}

// Load provides loading of data by the path on storage
func (l *Local) Load() ([]models.Parking, error) {
	d, err := ioutil.ReadFile(l.path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to read data by the path: %s", l.path))
	}
	var result []models.Parking
	if err := json.Unmarshal(d, &result); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal data")
	}
	return result, nil
}
