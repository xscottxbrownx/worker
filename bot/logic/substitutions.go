package logic

import (
	"fmt"
	"strings"

	"github.com/TicketsBot-cloud/worker"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
)

type SubstitutionFunc func(user user.User, member member.Member) string

type Substitutor struct {
	Placeholder string
	NeedsUser   bool
	NeedsMember bool
	F           SubstitutionFunc
}

func NewSubstitutor(placeholder string, needsUser, needsMember bool, f SubstitutionFunc) Substitutor {
	return Substitutor{
		Placeholder: placeholder,
		NeedsUser:   needsUser,
		NeedsMember: needsMember,
		F:           f,
	}
}

func doSubstitutions(worker *worker.Context, s string, userId uint64, guildId uint64, substitutors []Substitutor) (string, error) {
	var needsUser, needsMember bool

	// Determine which objects we need to fetch
	for _, substitutor := range substitutors {
		if substitutor.NeedsUser {
			needsUser = true
		}

		if substitutor.NeedsMember {
			needsMember = true
		}

		if needsUser && needsMember {
			break
		}
	}

	// Retrieve user and member if necessary
	var user user.User
	var member member.Member

	var err error
	if needsUser {
		user, err = worker.GetUser(userId)
	}

	if err != nil {
		return "", err
	}

	if needsMember {
		member, err = worker.GetGuildMember(guildId, userId)
	}

	if err != nil {
		return "", err
	}

	for _, substitutor := range substitutors {
		placeholder := fmt.Sprintf("%%%s%%", substitutor.Placeholder)

		if strings.Contains(s, placeholder) {
			s = strings.ReplaceAll(s, placeholder, substitutor.F(user, member))
		}
	}

	return s, nil
}
