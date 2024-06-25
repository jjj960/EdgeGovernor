package workflow

import (
	"EdgeGovernor/pkg/models"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/heimdalr/dag"
	"os"
)

func ReadWorkflow(dataType string, data interface{}) (*dag.DAG, models.Adag, error) {
	switch dataType {
	case "string":
		strData, ok := data.(string)
		if !ok {
			return nil, models.Adag{}, errors.New("Error: data is not of type string")
		}

		xmlFile, err := os.Open(strData)
		if err != nil {
			return nil, models.Adag{}, fmt.Errorf("Error opening file: %s", err)
		}
		defer xmlFile.Close()

		var adag models.Adag
		err = xml.NewDecoder(xmlFile).Decode(&adag)
		if err != nil {
			return nil, models.Adag{}, fmt.Errorf("Error decoding XML: %s", err)
		}

		d, err := GenerateDAG(adag)
		if err != nil {
			return nil, models.Adag{}, err
		}

		return d, adag, nil

	case "[]byte":
		xmlData, ok := data.([]byte)
		if !ok {
			return nil, models.Adag{}, errors.New("Error: data is not of type []byte")
		}
		var adag models.Adag
		err := xml.Unmarshal(xmlData, &adag)
		if err != nil {
			return nil, models.Adag{}, fmt.Errorf("Error unmarshalling XML: %s", err)
		}

		d, err := GenerateDAG(adag)
		if err != nil {
			return nil, models.Adag{}, err
		}

		return d, models.Adag{}, nil

	default:
		return nil, models.Adag{}, errors.New("Error: unsupported data type")
	}

}
