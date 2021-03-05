package mgmt

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/apps/router/mgmt/entities"
	"github.com/skupperproject/skupper/client"
	"github.com/skupperproject/skupper/test/utils/k8s"
)

var (
	queryCommand = []string{"qdmanage", "query", "--type"}
)

// QdmanageQuery executes a "qdmanager query" command on the provided pod, returning
// a slice of entities of the provided "entity" type.
func QdmanageQuery(client *client.VanClient, namespace string, pod string, container string,
	entity entities.Entity, fn func(entities.Entity) bool) ([]entities.Entity, error) {
	// Preparing command to execute
	command := append(queryCommand, entity.GetEntityId())
	stdout, stderr, err := k8s.Execute(client.KubeClient, client.RestConfig, namespace, pod, container, command)
	if err != nil {
		return nil, fmt.Errorf("error executing: %v - error: %v - stderr: %v", command, err, stderr)
	}

	// Using reflection to get a slice instance of the concrete type
	vo := reflect.TypeOf(entity)
	v := reflect.SliceOf(vo)
	nv := reflect.New(v)
	//fmt.Printf("v    - %T - %v\n", v, v)
	//fmt.Printf("nv   - %T - %v\n", nv, nv)

	// Unmarshalling to a slice of the concrete Entity type provided via "entity" instance
	err = json.Unmarshal(stdout.Bytes(), nv.Interface())
	if err != nil {
		//fmt.Printf("ERROR: %v\n", err)
		return nil, err
	}

	// Adding each parsed concrete Entity to the parsedEntities
	parsedEntities := []entities.Entity{}
	for i := 0; i < nv.Elem().Len(); i++ {
		candidate := nv.Elem().Index(i).Interface().(entities.Entity)

		// If no filter function provided, just add
		if fn == nil {
			parsedEntities = append(parsedEntities, candidate)
			continue
		}

		// Otherwhise invoke to determine whether to include
		if fn(candidate) {
			parsedEntities = append(parsedEntities, candidate)
		}
	}

	return parsedEntities, err
}
