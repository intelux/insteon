package plm

// Aliases represents a database of aliases.
type Aliases interface {
	Add(alias string, identity Identity)
	Remove(alias string)
	ParseIdentity(s string) (Identity, error)
}

type aliases map[string]Identity

func (a aliases) Add(alias string, identity Identity) {
	a[alias] = identity
}

func (a aliases) Remove(alias string) {
	delete(a, alias)
}

func (a aliases) ParseIdentity(s string) (Identity, error) {
	if identity, ok := a[s]; ok {
		return identity, nil
	}

	return ParseIdentity(s)
}
