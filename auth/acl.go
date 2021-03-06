package auth

import (
	"encoding/json"
	"errors"

	"github.com/ARGOeu/argo-messaging/stores"
)

// ACL holds the authorized users for a resource (topic/subscription)
type ACL struct {
	AuthUsers []string `json:"authorized_users"`
}

// ExportJSON export topic acl body to json for use in http response
func (acl *ACL) ExportJSON() (string, error) {
	if acl.AuthUsers == nil {
		acl.AuthUsers = make([]string, 0)
	}
	output, err := json.MarshalIndent(acl, "", "   ")
	return string(output[:]), err
}

// GetACLFromJSON retrieves ACL info from JSON
func GetACLFromJSON(input []byte) (ACL, error) {
	acl := ACL{}
	err := json.Unmarshal([]byte(input), &acl)
	if acl.AuthUsers == nil {
		return acl, errors.New("wrong argument")
	}
	return acl, err
}

// ModACL is called to modify an acl
func ModACL(projectUUID string, resourceType string, resourceName string, acl []string, store stores.Store) error {
	// Transform user name to user uuid

	userUUIDs := []string{}
	for _, username := range acl {
		userUUID := GetUUIDByName(username, store)
		userUUIDs = append(userUUIDs, userUUID)
	}

	return store.ModACL(projectUUID, resourceType, resourceName, userUUIDs)
}

// GetACL returns an authorized list of user for the resource (topic or subscription)
func GetACL(projectUUID string, resourceType string, resourceName string, store stores.Store) (ACL, error) {
	result := ACL{}
	acl, err := store.QueryACL(projectUUID, resourceType, resourceName)
	if err != nil {
		return result, err
	}
	for _, item := range acl.ACL {
		// Get Username from user uuid
		username := GetNameByUUID(item, store)
		result.AuthUsers = append(result.AuthUsers, username)
	}
	return result, err
}
