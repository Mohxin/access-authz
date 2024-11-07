package authz

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	cachemanager "github.com/volvo-cars/connect-access-control/internal/pkg/gateway/cache-manager"
	"github.com/volvo-cars/connect-access-control/internal/pkg/gateway/plums"
	"github.com/volvo-cars/connect-access-control/internal/pkg/store"
)

//go:generate mockgen -source=authz.go -destination=mocks/authz_mock.go -package=mocks

const (
	azureIdentityProvider = "AzureAD_VCC"
	cdsIDDelimiter        = "@"
)

var ErrUserNotFound = errors.New("user not found")

//go:generate
type cacheClient interface {
	GetPartnersByCodes(ctx context.Context, partnerCodes []string, partnerType string) ([]*cachemanager.Partner, error)
}

//go:generate
type plumsClient interface {
	GetUserByCDSID(ctx context.Context, cdsid string) (*plums.User, error)
}

//go:generate
type authzStore interface {
	GetRoleMapping(scopeID, roleID string) (store.RoleMapping, error)
	GetRoleMappings(scopeID string) ([]store.RoleMapping, error)
}

type Service struct {
	cache      cacheClient
	plums      plumsClient
	authzStore authzStore
}

func NewService(cache cacheClient, plums plumsClient, authzStore authzStore) *Service {
	return &Service{
		cache:      cache,
		plums:      plums,
		authzStore: authzStore,
	}
}

func (s *Service) GetUserAccess(ctx context.Context, cdsid string, scopes []string) ([]UserAccess, error) {
	user, err := s.GetUserByCDSID(ctx, cdsid)
	if err != nil {
		return nil, fmt.Errorf("GetUserAccess error: %w", err)
	}

	userType := detectUserType(user)

	var accesses []UserAccess
	for _, partner := range user.Partners {
		// Evaluate role mappings
		permissionGroups, err := s.evaluateRoleAccess(partner, scopes, userType)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate role mappings error: %w", err)
		}

		if len(permissionGroups) == 0 {
			continue
		}

		accesses = append(accesses, UserAccess{
			Context: Context{
				ID:   partner.ID,
				Type: partner.Type,
				Tag:  partner.ParmaPartnerCode,
			},
			Roles:            partner.Roles,
			PermissionGroups: permissionGroups,
		})
	}

	return accesses, nil
}

func (s *Service) GetUserByCDSID(ctx context.Context, cdsid string) (User, error) {
	user, err := s.plums.GetUserByCDSID(ctx, cdsid)
	if err != nil {
		if errors.Is(err, plums.ErrUserNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("GetUser error: %w", err)
	}

	return s.buildUserInfo(ctx, user)
}

func (s *Service) buildUserInfo(ctx context.Context, plumsUser *plums.User) (User, error) {
	cdsid := getCdsIDFromUserIdentities(plumsUser.UserIdentities)

	partners, err := s.buildPartners(ctx, plumsUser.Partners)
	if err != nil {
		return User{}, fmt.Errorf("failed to build partners error: %w", err)
	}

	return User{
		ID:          plumsUser.UserID,
		Email:       plumsUser.Email,
		CDSID:       cdsid,
		CountryCode: plumsUser.CountryCode,
		Partners:    partners,
	}, nil
}

func (s *Service) buildPartners(ctx context.Context, partners []plums.Partner) ([]Partner, error) {
	partnersByID := make(map[string]plums.Partner, len(partners))
	partnersByType := make(map[string][]string)

	// Populate partnersByID and partnersByType
	for _, partner := range partners {
		partnersByID[partner.PartnerID] = partner
		partnersByType[partner.PartnerType] = append(partnersByType[partner.PartnerType], partner.PartnerID)
	}

	// Use buffered channel to collect results and errors
	resultChan := make(chan Partner, len(partners))
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	// Concurrently fetch partners by type
	for partnerType, partnerIds := range partnersByType {
		wg.Add(1)
		go func(typ string, ids []string) {
			defer wg.Done()

			// Fetch partner info from cache
			cachedPartners, err := s.cache.GetPartnersByCodes(ctx, ids, typ)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("failed to build partners error: %w", err):
				default:
				}
				return
			}

			// Process cached partners
			for _, cachedPartner := range cachedPartners {
				pp := getPartners(partnersByID, typ, cachedPartner)

				resultChan <- Partner{
					ID:               cachedPartner.ID,
					RoleCode:         cachedPartner.RoleCode,
					Name:             cachedPartner.Name,
					Type:             typ,
					DistributorID:    cachedPartner.DistributorID,
					ParmaPartnerCode: cachedPartner.ParmaPartnerCode,
					Market:           cachedPartner.Market,
					Active:           cachedPartner.Active,
					IsPrimary:        pp.IsPrimary,
					Roles:            pp.Roles,
				}
			}
		}(partnerType, partnerIds)
	}

	// Close the result channel once all goroutines finish
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	// Collect results and check for errors
	var result []Partner
	for partner := range resultChan {
		result = append(result, partner)
	}

	// Check for any errors
	if err := <-errChan; err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) evaluateRoleAccess(partner Partner, scopes []string, userType string) (map[string][]string, error) {
	permissionGroups := make(map[string][]string)
	for _, roleID := range partner.Roles {
		for _, scope := range scopes {
			roleMapping, err := s.authzStore.GetRoleMapping(scope, roleID)
			if err != nil {
				if errors.Is(err, store.ErrRoleMappingNotFound) {
					continue
				}

				return nil, fmt.Errorf("GetRoleMapping error: %w", err)
			}

			for _, mapping := range roleMapping.Mapping {
				if len(mapping.Filter.Market) > 0 && !contains(mapping.Filter.Market, partner.Market) {
					continue
				}

				if len(mapping.Filter.UserType) > 0 && !contains(mapping.Filter.UserType, userType) {
					continue
				}

				if len(mapping.Filter.PartnerType) > 0 && !contains(mapping.Filter.PartnerType, partner.Type) {
					continue
				}

				// TODO: Q: should only show matched roleID(s) in the response?
				// Append permission groups
				permissionGroups[scope] = append(permissionGroups[scope], mapping.PermissionGroups...)
			}
		}
	}

	return permissionGroups, nil
}

func getPartners(partners map[string]plums.Partner, partnerType string, cachedPartner *cachemanager.Partner) plums.Partner {
	switch partnerType {
	case "PARMA":
		if p, ok := partners[cachedPartner.ParmaPartnerCode]; ok {
			return p
		}
	case "NSC":
		if p, ok := partners[cachedPartner.ID]; ok {
			return p
		}
	}

	return plums.Partner{}
}

func getCdsIDFromUserIdentities(identities []plums.UserIdentity) string {
	const minParts = 2
	for _, identity := range identities {
		if identity.Provider == azureIdentityProvider {
			parts := strings.Split(identity.AccountName, cdsIDDelimiter)

			if len(parts) < minParts {
				return ""
			}
			return parts[0]
		}
	}

	return ""
}

func detectUserType(user User) string {
	userType := ""
	if strings.HasSuffix(user.Email, "@volvocars.com") {
		userType = "INTERNAL"
	} else if strings.HasSuffix(user.Email, "@volvocars.biz") {
		userType = "EXTERNAL"
	}
	return userType
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
