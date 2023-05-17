package svc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cs3permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	filestore "github.com/owncloud/ocis/v2/services/settings/pkg/store/filesystem"
	metastore "github.com/owncloud/ocis/v2/services/settings/pkg/store/metadata"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service represents a service.
type Service struct {
	id      string
	config  *config.Config
	logger  log.Logger
	manager settings.Manager
}

// NewService returns a service implementation for Service.
func NewService(cfg *config.Config, logger log.Logger) Service {
	service := Service{
		id:     "ocis-settings",
		config: cfg,
		logger: logger,
	}

	switch cfg.StoreType {
	default:
		fallthrough
	case "metadata":
		service.manager = metastore.New(cfg)
	case "filesystem":
		service.manager = filestore.New(cfg)
		// TODO: if we want to further support filesystem store it should use default permissions from store/defaults/defaults.go instead using this duplicate
		service.RegisterDefaultRoles()
	}
	return service
}

// CheckPermission implements the CS3 API Permssions service.
// It's used to check if a subject (user or group) has a permission.
func (g Service) CheckPermission(ctx context.Context, req *cs3permissions.CheckPermissionRequest) (*cs3permissions.CheckPermissionResponse, error) {
	spec := req.SubjectRef.Spec

	var accountID string
	switch ref := spec.(type) {
	case *cs3permissions.SubjectReference_UserId:
		accountID = ref.UserId.OpaqueId
	case *cs3permissions.SubjectReference_GroupId:
		accountID = ref.GroupId.OpaqueId
	}

	assignments, err := g.manager.ListRoleAssignments(accountID)
	if err != nil {
		return &cs3permissions.CheckPermissionResponse{
			Status: status.NewInternal(ctx, err.Error()),
		}, nil
	}

	roleIDs := make([]string, 0, len(assignments))
	for _, a := range assignments {
		roleIDs = append(roleIDs, a.RoleId)
	}

	permission, err := g.manager.ReadPermissionByName(req.Permission, roleIDs)
	if err != nil {
		if !errors.Is(err, settings.ErrNotFound) {
			return &cs3permissions.CheckPermissionResponse{
				Status: status.NewInternal(ctx, err.Error()),
			}, nil
		}
	}

	if permission == nil {
		return &cs3permissions.CheckPermissionResponse{
			Status: &rpcv1beta1.Status{
				Code: rpcv1beta1.Code_CODE_PERMISSION_DENIED,
			},
		}, nil
	}

	return &cs3permissions.CheckPermissionResponse{
		Status: status.NewOK(ctx),
	}, nil
}

// RegisterDefaultRoles composes default roles and saves them. Skipped if the roles already exist.
func (g Service) RegisterDefaultRoles() {
	// FIXME: we're writing default roles per service start (i.e. twice at the moment, for http and grpc server). has to happen only once.
	for _, role := range generateBundlesDefaultRoles() {
		bundleID := role.Extension + "." + role.Id
		// check if the role already exists
		bundle, _ := g.manager.ReadBundle(role.Id)
		if bundle != nil {
			g.logger.Debug().Str("bundleID", bundleID).Msg("bundle already exists. skipping.")
			continue
		}
		// create the role
		_, err := g.manager.WriteBundle(role)
		if err != nil {
			g.logger.Error().Err(err).Str("bundleID", bundleID).Msg("failed to register bundle")
		}
		g.logger.Debug().Str("bundleID", bundleID).Msg("successfully registered bundle")
	}

	for _, req := range generatePermissionRequests() {
		_, err := g.manager.AddSettingToBundle(req.GetBundleId(), req.GetSetting())
		if err != nil {
			g.logger.Error().
				Err(err).
				Str("bundleID", req.GetBundleId()).
				Interface("setting", req.GetSetting()).
				Msg("failed to register permission")
		}
	}

	if g.config.SetupDefaultAssignments {
		for _, req := range g.defaultRoleAssignments() {
			if _, err := g.manager.WriteRoleAssignment(req.AccountUuid, req.RoleId); err != nil {
				g.logger.Error().Err(err).Msg("failed to register role assignment")
			}
		}
	}
}

// TODO: check permissions on every request

// SaveBundle implements the BundleServiceHandler interface
func (g Service) SaveBundle(ctx context.Context, req *settingssvc.SaveBundleRequest, res *settingssvc.SaveBundleResponse) error {
	cleanUpResource(ctx, req.Bundle.Resource)
	if err := g.checkStaticPermissionsByBundleType(ctx, req.Bundle.Type); err != nil {
		return err
	}
	if validationError := validateSaveBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}

	r, err := g.manager.WriteBundle(req.Bundle)
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Bundle = r
	return nil
}

// GetBundle implements the BundleServiceHandler interface
func (g Service) GetBundle(ctx context.Context, req *settingssvc.GetBundleRequest, res *settingssvc.GetBundleResponse) error {
	if validationError := validateGetBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	bundle, err := g.manager.ReadBundle(req.BundleId)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	filteredBundle := g.getFilteredBundle(g.getRoleIDs(ctx), bundle)
	if len(filteredBundle.Settings) == 0 {
		err = fmt.Errorf("could not read bundle: %s", req.BundleId)
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Bundle = filteredBundle
	return nil
}

// ListBundles implements the BundleServiceHandler interface
func (g Service) ListBundles(ctx context.Context, req *settingssvc.ListBundlesRequest, res *settingssvc.ListBundlesResponse) error {
	// fetch all bundles
	if validationError := validateListBundles(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	bundles, err := g.manager.ListBundles(settingsmsg.Bundle_TYPE_DEFAULT, req.BundleIds)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	roleIDs := g.getRoleIDs(ctx)

	// filter settings in bundles that are allowed according to roles
	var filteredBundles []*settingsmsg.Bundle
	for _, bundle := range bundles {
		filteredBundle := g.getFilteredBundle(roleIDs, bundle)
		if len(filteredBundle.Settings) > 0 {
			filteredBundles = append(filteredBundles, filteredBundle)
		}
	}

	res.Bundles = filteredBundles
	return nil
}

func (g Service) getFilteredBundle(roleIDs []string, bundle *settingsmsg.Bundle) *settingsmsg.Bundle {
	// check if full bundle is whitelisted
	bundleResource := &settingsmsg.Resource{
		Type: settingsmsg.Resource_TYPE_BUNDLE,
		Id:   bundle.Id,
	}
	if g.hasPermission(
		roleIDs,
		bundleResource,
		[]settingsmsg.Permission_Operation{settingsmsg.Permission_OPERATION_READ, settingsmsg.Permission_OPERATION_READWRITE},
		settingsmsg.Permission_CONSTRAINT_OWN,
	) {
		return bundle
	}

	// filter settings based on permissions
	var filteredSettings []*settingsmsg.Setting
	for _, setting := range bundle.Settings {
		settingResource := &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SETTING,
			Id:   setting.Id,
		}
		if g.hasPermission(
			roleIDs,
			settingResource,
			[]settingsmsg.Permission_Operation{settingsmsg.Permission_OPERATION_READ, settingsmsg.Permission_OPERATION_READWRITE},
			settingsmsg.Permission_CONSTRAINT_OWN,
		) {
			filteredSettings = append(filteredSettings, setting)
		}
	}
	bundle.Settings = filteredSettings
	return bundle
}

// AddSettingToBundle implements the BundleServiceHandler interface
func (g Service) AddSettingToBundle(ctx context.Context, req *settingssvc.AddSettingToBundleRequest, res *settingssvc.AddSettingToBundleResponse) error {
	cleanUpResource(ctx, req.Setting.Resource)
	if err := g.checkStaticPermissionsByBundleID(ctx, req.BundleId); err != nil {
		return err
	}
	if validationError := validateAddSettingToBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}

	r, err := g.manager.AddSettingToBundle(req.BundleId, req.Setting)
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Setting = r
	return nil
}

// RemoveSettingFromBundle implements the BundleServiceHandler interface
func (g Service) RemoveSettingFromBundle(ctx context.Context, req *settingssvc.RemoveSettingFromBundleRequest, _ *emptypb.Empty) error {
	if err := g.checkStaticPermissionsByBundleID(ctx, req.BundleId); err != nil {
		return err
	}
	if validationError := validateRemoveSettingFromBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}

	if err := g.manager.RemoveSettingFromBundle(req.BundleId, req.SettingId); err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}

	return nil
}

// SaveValue implements the ValueServiceHandler interface
func (g Service) SaveValue(ctx context.Context, req *settingssvc.SaveValueRequest, res *settingssvc.SaveValueResponse) error {
	req.Value.AccountUuid = getValidatedAccountUUID(ctx, req.Value.AccountUuid)

	if !g.isCurrentUser(ctx, req.Value.AccountUuid) {
		return merrors.Forbidden(g.id, "can't save value for another user")
	}

	cleanUpResource(ctx, req.Value.Resource)
	// TODO: we need to check, if the authenticated user has permission to write the value for the specified resource (e.g. global, file with id xy, ...)
	if validationError := validateSaveValue(req); validationError != nil {
		return merrors.BadRequest(g.id, validationError.Error())
	}
	r, err := g.manager.WriteValue(req.Value)
	if err != nil {
		return merrors.BadRequest(g.id, err.Error())
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return merrors.NotFound(g.id, err.Error())
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValue implements the ValueServiceHandler interface
func (g Service) GetValue(ctx context.Context, req *settingssvc.GetValueRequest, res *settingssvc.GetValueResponse) error {
	if validationError := validateGetValue(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ReadValue(req.Id)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValueByUniqueIdentifiers implements the ValueService interface
func (g Service) GetValueByUniqueIdentifiers(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, res *settingssvc.GetValueResponse) error {
	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if !g.isCurrentUser(ctx, req.AccountUuid) {
		return merrors.Forbidden(g.id, "can't get value of another user")
	}
	if validationError := validateGetValueByUniqueIdentifiers(req); validationError != nil {
		return merrors.BadRequest(g.id, validationError.Error())
	}
	v, err := g.manager.ReadValueByUniqueIdentifiers(req.AccountUuid, req.SettingId)
	if err != nil {
		return merrors.NotFound(g.id, err.Error())
	}

	if v.BundleId != "" {
		valueWithIdentifier, err := g.getValueWithIdentifier(v)
		if err != nil {
			return merrors.NotFound(g.id, err.Error())
		}

		res.Value = valueWithIdentifier
	}
	return nil
}

// ListValues implements the ValueServiceHandler interface
func (g Service) ListValues(ctx context.Context, req *settingssvc.ListValuesRequest, res *settingssvc.ListValuesResponse) error {
	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if !g.isCurrentUser(ctx, req.AccountUuid) {
		return merrors.Forbidden(g.id, "can't list values of another user")
	}

	if validationError := validateListValues(req); validationError != nil {
		return merrors.BadRequest(g.id, validationError.Error())
	}
	values, err := g.manager.ListValues(req.BundleId, req.AccountUuid)
	if err != nil {
		return merrors.NotFound(g.id, err.Error())
	}
	result := make([]*settingsmsg.ValueWithIdentifier, 0, len(values))
	for _, value := range values {
		valueWithIdentifier, err := g.getValueWithIdentifier(value)
		if err == nil {
			result = append(result, valueWithIdentifier)
		}
	}
	res.Values = result
	return nil
}

// ListRoles implements the RoleServiceHandler interface
func (g Service) ListRoles(c context.Context, req *settingssvc.ListBundlesRequest, res *settingssvc.ListBundlesResponse) error {
	//accountUUID := getValidatedAccountUUID(c, "me")
	if validationError := validateListRoles(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ListBundles(settingsmsg.Bundle_TYPE_ROLE, req.BundleIds)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	// TODO: only allow to list roles when user has account/role/... management permissions
	res.Bundles = r
	return nil
}

// ListRoleAssignments implements the RoleServiceHandler interface
func (g Service) ListRoleAssignments(ctx context.Context, req *settingssvc.ListRoleAssignmentsRequest, res *settingssvc.ListRoleAssignmentsResponse) error {
	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if validationError := validateListRoleAssignments(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ListRoleAssignments(req.AccountUuid)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Assignments = r
	return nil
}

// AssignRoleToUser implements the RoleServiceHandler interface
func (g Service) AssignRoleToUser(ctx context.Context, req *settingssvc.AssignRoleToUserRequest, res *settingssvc.AssignRoleToUserResponse) error {
	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if validationError := validateAssignRoleToUser(req); validationError != nil {
		return merrors.BadRequest(g.id, validationError.Error())
	}

	ownAccountUUID, ok := metadata.Get(ctx, middleware.AccountID)
	if !ok {
		g.logger.Debug().Str("id", g.id).Msg("user not in context")
		return merrors.InternalServerError(g.id, "user not in context")
	}

	switch {
	case ownAccountUUID == req.AccountUuid:
		// Allow users to assign themself to the user role
		// deny any other attempt to change	the user's own assignment
		if r, err := g.manager.ListRoleAssignments(req.AccountUuid); err == nil && len(r) > 0 {
			return merrors.Forbidden(g.id, "Changing own role assignment forbidden")
		}
		if req.RoleId != defaults.BundleUUIDRoleUser {
			return merrors.Forbidden(g.id, "Changing own role assignment forbidden")
		}
		g.logger.Debug().Str("userid", ownAccountUUID).Msg("Self-assignment for default 'user' role permitted")
	case g.canManageRoles(ctx):
	default:
		return merrors.Forbidden(g.id, "user has no role management permission")
	}

	r, err := g.manager.WriteRoleAssignment(req.AccountUuid, req.RoleId)
	if err != nil {
		return merrors.BadRequest(g.id, err.Error())
	}
	res.Assignment = r
	return nil
}

// RemoveRoleFromUser implements the RoleServiceHandler interface
func (g Service) RemoveRoleFromUser(ctx context.Context, req *settingssvc.RemoveRoleFromUserRequest, _ *emptypb.Empty) error {
	if !g.canManageRoles(ctx) {
		return merrors.Forbidden(g.id, "user has no role management permission")
	}

	if validationError := validateRemoveRoleFromUser(req); validationError != nil {
		return merrors.BadRequest(g.id, validationError.Error())
	}

	ownAccountUUID, ok := metadata.Get(ctx, middleware.AccountID)
	if !ok {
		g.logger.Debug().Str("id", g.id).Msg("user not in context")
		return merrors.InternalServerError(g.id, "user not in context")
	}

	al, err := g.manager.ListRoleAssignments(ownAccountUUID)
	if err != nil {
		g.logger.Debug().Err(err).Str("id", g.id).Msg("ListRoleAssignments failed")
		return merrors.InternalServerError(g.id, err.Error())
	}

	for _, a := range al {
		if a.Id == req.Id {
			g.logger.Debug().Str("id", g.id).Msg("Removing own role assignment forbidden")
			return merrors.Forbidden(g.id, "Removing own role assignment forbidden")
		}
	}

	if err := g.manager.RemoveRoleAssignment(req.Id); err != nil {
		return merrors.BadRequest(g.id, err.Error())
	}
	return nil
}

// ListPermissions implements the PermissionServiceHandler interface
func (g Service) ListPermissions(ctx context.Context, req *settingssvc.ListPermissionsRequest, res *settingssvc.ListPermissionsResponse) error {
	ownAccountUUID, ok := metadata.Get(ctx, middleware.AccountID)
	if !ok {
		g.logger.Debug().Str("id", g.id).Msg("user not in context")
		return merrors.InternalServerError(g.id, "user not in context")
	}

	if ownAccountUUID != req.AccountUuid {
		return merrors.NotFound(g.id, "user not found: %s", req.AccountUuid)
	}

	assignments, err := g.manager.ListRoleAssignments(req.AccountUuid)
	if err != nil {
		return err
	}

	// deduplicate role ids
	roleIDs := map[string]struct{}{}
	for _, a := range assignments {
		roleIDs[a.RoleId] = struct{}{}
	}

	// deduplicate permission names
	permissionNames := map[string]struct{}{}
	for roleID := range roleIDs {
		bundle, err := g.manager.ReadBundle(roleID)
		if err != nil {
			if !errors.Is(err, settings.ErrNotFound) {
				return err
			}
			continue
		}

		if bundle != nil {
			for _, setting := range bundle.GetSettings() {
				permissionNames[formatPermissionName(setting)] = struct{}{}
			}
		}
	}

	res.Permissions = make([]string, 0, len(permissionNames))
	for p := range permissionNames {
		res.Permissions = append(res.Permissions, p)
	}

	return nil
}

// ListPermissionsByResource implements the PermissionServiceHandler interface
func (g Service) ListPermissionsByResource(ctx context.Context, req *settingssvc.ListPermissionsByResourceRequest, res *settingssvc.ListPermissionsByResourceResponse) error {
	if validationError := validateListPermissionsByResource(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	permissions, err := g.manager.ListPermissionsByResource(req.Resource, g.getRoleIDs(ctx))
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Permissions = permissions
	return nil
}

// GetPermissionByID implements the PermissionServiceHandler interface
func (g Service) GetPermissionByID(ctx context.Context, req *settingssvc.GetPermissionByIDRequest, res *settingssvc.GetPermissionByIDResponse) error {
	if validationError := validateGetPermissionByID(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	permission, err := g.manager.ReadPermissionByID(req.PermissionId, g.getRoleIDs(ctx))
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	if permission == nil {
		return merrors.NotFound(g.id, "%s", fmt.Errorf("permission %s not found in roles", req.PermissionId))
	}
	res.Permission = permission
	return nil
}

// cleanUpResource makes sure that the account uuid of the authenticated user is injected if needed.
func cleanUpResource(ctx context.Context, resource *settingsmsg.Resource) {
	if resource != nil && resource.Type == settingsmsg.Resource_TYPE_USER {
		resource.Id = getValidatedAccountUUID(ctx, resource.Id)
	}
}

// getValidatedAccountUUID converts `me` into an actual account uuid from the context, if possible.
// the result of this function will always be a valid lower-case UUID or an empty string.
func getValidatedAccountUUID(ctx context.Context, accountUUID string) string {
	if accountUUID == "me" {
		if ownAccountUUID, ok := metadata.Get(ctx, middleware.AccountID); ok {
			accountUUID = ownAccountUUID
		}
	}
	if accountUUID == "me" {
		// no matter what happens above, an accountUUID of `me` must not be passed on. Clear it instead.
		accountUUID = ""
	}
	return accountUUID
}

// getRoleIDs extracts the roleIDs of the authenticated user from the context.
func (g Service) getRoleIDs(ctx context.Context) []string {
	if ownRoleIDs, ok := roles.ReadRoleIDsFromContext(ctx); ok {
		return ownRoleIDs
	}
	if accountID, ok := metadata.Get(ctx, middleware.AccountID); ok {
		assignments, err := g.manager.ListRoleAssignments(accountID)
		if err != nil {
			g.logger.Info().Err(err).Str("userid", accountID).Msg("failed to get roles for user")
			return nil
		}

		ownRoleIDs := make([]string, 0, len(assignments))
		for _, a := range assignments {
			ownRoleIDs = append(ownRoleIDs, a.RoleId)
		}
		return ownRoleIDs
	}
	g.logger.Info().Msg("failed to get accountID from context")
	return nil
}

func (g Service) getValueWithIdentifier(value *settingsmsg.Value) (*settingsmsg.ValueWithIdentifier, error) {
	bundle, err := g.manager.ReadBundle(value.BundleId)
	if err != nil {
		return nil, err
	}
	setting, err := g.manager.ReadSetting(value.SettingId)
	if err != nil {
		return nil, err
	}
	return &settingsmsg.ValueWithIdentifier{
		Identifier: &settingsmsg.Identifier{
			Extension: bundle.Extension,
			Bundle:    bundle.Name,
			Setting:   setting.Name,
		},
		Value: value,
	}, nil
}

func (g Service) hasStaticPermission(ctx context.Context, permissionID string) bool {
	roleIDs, ok := roles.ReadRoleIDsFromContext(ctx)
	if !ok {
		// TODO add system role for internal requests.
		// - at least the proxy needs to look up account info
		// - glauth needs to make bind requests
		// tracked as OCIS-454

		accountID, ok := metadata.Get(ctx, middleware.AccountID)
		if !ok {
			return false
		}
		assignments, err := g.manager.ListRoleAssignments(accountID)
		if err != nil {
			return false
		}

		// deduplicate roleids
		uniqueRoleIds := make(map[string]struct{})
		for _, a := range assignments {
			uniqueRoleIds[a.GetRoleId()] = struct{}{}
		}
		roleIDs = make([]string, 0, len(uniqueRoleIds))
		for a := range uniqueRoleIds {
			roleIDs = append(roleIDs, a)
		}
	}
	p, err := g.manager.ReadPermissionByID(permissionID, roleIDs)
	return err == nil && p != nil
}

func (g Service) checkStaticPermissionsByBundleID(ctx context.Context, bundleID string) error {
	bundle, err := g.manager.ReadBundle(bundleID)
	if err != nil {
		return merrors.NotFound(g.id, "bundle not found: %s", err)
	}
	return g.checkStaticPermissionsByBundleType(ctx, bundle.Type)
}

func (g Service) checkStaticPermissionsByBundleType(ctx context.Context, bundleType settingsmsg.Bundle_Type) error {
	if bundleType == settingsmsg.Bundle_TYPE_ROLE {
		if !g.hasStaticPermission(ctx, RoleManagementPermissionID) {
			return merrors.Forbidden(g.id, "user has no role management permission")
		}
		return nil
	}
	if !g.hasStaticPermission(ctx, SettingsManagementPermissionID) {
		return merrors.Forbidden(g.id, "user has no settings management permission")
	}
	return nil
}

func (g Service) isCurrentUser(ctx context.Context, accountID string) bool {
	ownAccountID, ok := metadata.Get(ctx, middleware.AccountID)
	if !ok {
		return false
	}
	return accountID == ownAccountID
}

func (g Service) canManageRoles(ctx context.Context) bool {
	return g.hasStaticPermission(ctx, RoleManagementPermissionID)
}

func formatPermissionName(setting *settingsmsg.Setting) string {
	constraint := strings.TrimPrefix(setting.GetPermissionValue().GetConstraint().String(), "CONSTRAINT_")
	return setting.Name + "." + strings.ToLower(constraint)
}
