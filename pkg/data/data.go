package data

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	json "github.com/pquerna/ffjson/ffjson"
	"github.com/saromanov/golang-developer-test-task/pkg/models"
)

// Load provides loading of data
func Load(path string) ([]models.Parking, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to read data by the path: %s", path))
	}
	var result []models.Parking
	if err := json.Unmarshal(d, &result); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal data")
	}
	return result, nil
}
