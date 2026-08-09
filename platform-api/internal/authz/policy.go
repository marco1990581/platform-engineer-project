package authz

// Permission represents an allowed action on a resource.
type Permission struct {
	Resource string
	Action   string
}

// permissions maps a role name to the permissions it grants.
//
// This is a static, in-memory policy for the local development
// authorization model. It can be replaced later by a persisted or
// externally managed policy store without changing the Enforce call
// sites in handlers or middleware.
var permissions = map[string][]Permission{
	"admin": {
		{Resource: "*", Action: "*"},
	},
	"viewer": {
		{Resource: "system", Action: "read"},
		{Resource: "network", Action: "read"},
		{Resource: "filesystem", Action: "read"},
	},
}

// Enforce reports whether role is permitted to perform action on resource.
func Enforce(role, resource, action string) bool {
	for _, permission := range permissions[role] {
		if (permission.Resource == "*" || permission.Resource == resource) &&
			(permission.Action == "*" || permission.Action == action) {
			return true
		}
	}

	return false
}
