package chain

import (
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
)

var noopCheckTime *func()

type AuthorizationManager struct {
	control *Controller
	db      database.DataBase
}

func newAuthorizationManager(control *Controller) *AuthorizationManager {
	azInstance := &AuthorizationManager{}
	azInstance.control = control
	azInstance.db = control.DB
	return azInstance
}

type PermissionIdType common.IdType

func (a *AuthorizationManager) CreatePermission(account common.AccountName,
	name common.PermissionName,
	parent PermissionIdType,
	auth types.Authority,
	initialCreationTime common.TimePoint,
) *entity.PermissionObject {
	creationTime := initialCreationTime
	if creationTime == 1 {
		creationTime = a.control.PendingBlockTime()
	}

	permUsage := entity.PermissionUsageObject{}
	permUsage.LastUsed = creationTime
	err := a.db.Insert(&permUsage)
	if err != nil {
		log.Error("CreatePermission is error: %s", err)
	}

	perm := entity.PermissionObject{
		UsageId:     permUsage.ID,
		Parent:      common.IdType(parent),
		Owner:       account,
		Name:        name,
		LastUpdated: creationTime,
		Auth:        auth.ToSharedAuthority(),
	}
	err = a.db.Insert(&perm)
	if err != nil {
		log.Error("CreatePermission is error: %s", err)
	}
	return &perm
}

func (a *AuthorizationManager) ModifyPermission(permission *entity.PermissionObject, auth *types.Authority) {
	err := a.db.Modify(permission, func(po *entity.PermissionObject) {
		po.Auth = (*auth).ToSharedAuthority()
		po.LastUpdated = a.control.PendingBlockTime()
	})
	if err != nil {
		log.Error("ModifyPermission is error: %s", err)
	}
}

func (a *AuthorizationManager) RemovePermission(permission *entity.PermissionObject) {
	index, err := a.db.GetIndex("byParent", entity.PermissionObject{})
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	itr, err := index.LowerBound(entity.PermissionObject{Parent: permission.ID})
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	EosAssert(index.CompareEnd(itr), &ActionValidateException{}, "Cannot remove a permission which has children. Remove the children first.")
	usage := entity.PermissionUsageObject{ID: permission.UsageId}
	err = a.db.Find("id", usage, &usage)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	err = a.db.Remove(&usage)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	err = a.db.Remove(permission)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
}

func (a *AuthorizationManager) UpdatePermissionUsage(permission *entity.PermissionObject) {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	err := a.db.Find("id", puo, &puo)
	if err != nil {
		log.Error("UpdatePermissionUsage is error: %s", err)
	}
	err = a.db.Modify(&puo, func(p *entity.PermissionUsageObject) {
		puo.LastUsed = a.control.PendingBlockTime()
	})
	if err != nil {
		log.Error("UpdatePermissionUsage is error: %s", err)
	}
}

func (a *AuthorizationManager) GetPermissionLastUsed(permission *entity.PermissionObject) common.TimePoint {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	err := a.db.Find("id", puo, &puo)
	if err != nil {
		log.Error("GetPermissionLastUsed is error: %s", err)
	}
	return puo.LastUsed
}

func (a *AuthorizationManager) FindPermission(level *types.PermissionLevel) (p *entity.PermissionObject) { //TODO
	Try(func() {
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		err := a.db.Find("byOwner", po, &po)
		if err != nil {
			log.Warn("%v@%v don't find", po.Owner, po.Name)
			p = nil
			return
		}
		p = &po
	}).Catch(func(e Exception) {
		FcRethrowException(&PermissionQueryException{}, log.LvlWarn, "FindPermission is error")
	}).End()
	return p
}

func (a *AuthorizationManager) GetPermission(level *types.PermissionLevel) (p *entity.PermissionObject) {
	Try(func() {
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		err := a.db.Find("byOwner", po, &po)
		if err != nil {
			log.Warn("%v@%v don't find", po.Owner, po.Name)
			p = nil
			return
		}
		p = &po
	}).Catch(func(e Exception) {
		FcRethrowException(&PermissionQueryException{}, log.LvlWarn, "GetPermission is error")
	}).End()
	return p
}

func (a *AuthorizationManager) LookupLinkedPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p *common.PermissionName) {
	Try(func() {
		link := entity.PermissionLinkObject{}
		link.Account = authorizerAccount
		link.Code = scope
		link.MessageType = actName
		err := a.db.Find("byActionName", link, &link)
		if err != nil {
			link.MessageType = common.AccountName(common.N(""))
			err = a.db.Find("byActionName", link, &link)
		}
		if err == nil {
			p = &link.RequiredPermission
			return
		}
	}).End()

	return p
}

func (a *AuthorizationManager) LookupMinimumPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p *common.PermissionName) {
	if scope == common.DefaultConfig.SystemAccountName {
		EosAssert(actName != UpdateAuth{}.GetName() &&
			actName != DeleteAuth{}.GetName() &&
			actName != LinkAuth{}.GetName() &&
			actName != UnLinkAuth{}.GetName() &&
			actName != CancelDelay{}.GetName(),
			&UnlinkableMinPermissionAction{}, "cannot call lookup_minimum_permission on native actions that are not allowed to be linked to minimum permissions")
	}
	Try(func() {
		linkedPermission := a.LookupLinkedPermission(authorizerAccount, scope, actName)
		if common.Empty(linkedPermission) {
			p = &common.DefaultConfig.ActiveName
			return
		}

		if *linkedPermission == common.PermissionName(common.DefaultConfig.EosioAnyName) {
			return
		}

		p = linkedPermission
		return
	}).End()
	return p
}

func (a *AuthorizationManager) CheckUpdateAuthAuthorization(update UpdateAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "UpdateAuth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == update.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.FindPermission(&types.PermissionLevel{Actor: update.Account, Permission: update.Permission})
	if minPermission == nil {
		permission := a.GetPermission(&types.PermissionLevel{Actor: update.Account, Permission: update.Parent})
		minPermission = permission
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckUpdateAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission, permissionIndex), &IrrelevantAuthException{},
		"UpdateAuth action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{update.Account, minPermission.Name})
}

func (a *AuthorizationManager) CheckDeleteAuthAuthorization(del DeleteAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "DeleteAuth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == del.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.GetPermission(&types.PermissionLevel{Actor: del.Account, Permission: del.Permission})
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckDeleteAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission, permissionIndex), &IrrelevantAuthException{},
		"DeleteAuth action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{minPermission.Owner, minPermission.Name})
}

func (a *AuthorizationManager) CheckLinkAuthAuthorization(link LinkAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "link action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == link.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")

	EosAssert(link.Type != UpdateAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::UpdateAuth to a minimum permission")
	EosAssert(link.Type != DeleteAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::DeleteAuth to a minimum permission")
	EosAssert(link.Type != LinkAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::LinkAuth to a minimum permission")
	EosAssert(link.Type != UnLinkAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::UnLinkAuth to a minimum permission")
	EosAssert(link.Type != CancelDelay{}.GetName(), &ActionValidateException{}, "Cannot link eosio::CancelDelay to a minimum permission")

	linkedPermissionName := a.LookupMinimumPermission(link.Account, link.Code, link.Type)
	if common.Empty(&linkedPermissionName) {
		return
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckLinkAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&types.PermissionLevel{link.Account, *linkedPermissionName}), permissionIndex), &IrrelevantAuthException{},
		"LinkAuth action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{link.Account, *linkedPermissionName})
}

func (a *AuthorizationManager) CheckUnLinkAuthAuthorization(unlink UnLinkAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "unlink action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == unlink.Account, &IrrelevantAuthException{},
		"the owner of the affected permission needs to be the actor of the declared authorization")

	unlinkedPermissionName := a.LookupLinkedPermission(unlink.Account, unlink.Code, unlink.Type)
	EosAssert(!common.Empty(&unlinkedPermissionName), &TransactionException{},
		"cannot unlink non-existent permission link of account '%v' for actions matching '%v::%v", unlink.Account, unlink.Code, unlink.Type)

	if *unlinkedPermissionName == common.DefaultConfig.EosioAnyName {
		return
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckUnLinkAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&types.PermissionLevel{unlink.Account, *unlinkedPermissionName}), permissionIndex), &IrrelevantAuthException{},
		"unlink action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{unlink.Account, *unlinkedPermissionName})
}

func (a *AuthorizationManager) CheckCancelDelayAuthorization(cancel CancelDelay, auths []types.PermissionLevel) common.Microseconds {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "CancelDelay action should only have one declared authorization")
	auth := auths[0]
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckCancelDelayAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&cancel.CancelingAuth), permissionIndex), &IrrelevantAuthException{},
		"CancelDelay action declares irrelevant authority '%v'; specified authority to satisfy is %v", auth, cancel.CancelingAuth)

	generatedTrx := entity.GeneratedTransactionObject{}
	trxId := cancel.TrxId
	generatedIndex, err := a.control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
	if err != nil {
		log.Error("CheckCancelDelayAuthorization is error: %s", err)
	}
	itr, err := generatedIndex.LowerBound(entity.GeneratedTransactionObject{TrxId: trxId})
	if err != nil {
		log.Error("CheckCancelDelayAuthorization is error: %s", err)
	}

	generatedIndex.BeginData(&generatedTrx)
	EosAssert(!generatedIndex.CompareEnd(itr) && generatedTrx.TrxId == trxId, &TxNotFound{},
		"cannot cancel trx_id=%v, there is no deferred transaction with that transaction id", trxId)

	trx := types.Transaction{}
	rlp.DecodeBytes(generatedTrx.PackedTrx, &trx)
	found := false
	for _, act := range trx.Actions {
		for _, auth := range act.Authorization {
			if auth == cancel.CancelingAuth {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	EosAssert(found, &ActionValidateException{}, "canceling_auth in CancelDelay action was not found as authorization in the original delayed transaction")
	return common.Milliseconds(int64(generatedTrx.DelayUntil) - int64(generatedTrx.Published))
}

func (a *AuthorizationManager) CheckAuthorization(actions []*types.Action,
	providedKeys *treeset.Set,
	providedPermissions *treeset.Set,
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	delayMaxLimit := common.Seconds(int64(a.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	var effectiveProvidedDelay common.Microseconds
	if providedDelay >= delayMaxLimit {
		effectiveProvidedDelay = common.MaxMicroseconds()
	} else {
		effectiveProvidedDelay = providedDelay
	}
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority {
		perm := a.GetPermission(p)
		if perm != nil {
			return perm.Auth
		} else {
			return types.SharedAuthority{}
		}
	},
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		providedKeys,
		providedPermissions,
		effectiveProvidedDelay,
		checkTime,
	)
	permissionToSatisfy := make(map[types.PermissionLevel]common.Microseconds)

	for _, act := range actions {
		specialCase := false
		delay := effectiveProvidedDelay

		if act.Account == common.DefaultConfig.SystemAccountName {
			specialCase = true
			switch act.Name {
			case UpdateAuth{}.GetName():
				UpdateAuth := UpdateAuth{}
				rlp.DecodeBytes(act.Data, &UpdateAuth)
				a.CheckUpdateAuthAuthorization(UpdateAuth, act.Authorization)

			case DeleteAuth{}.GetName():
				DeleteAuth := DeleteAuth{}
				rlp.DecodeBytes(act.Data, &DeleteAuth)
				a.CheckDeleteAuthAuthorization(DeleteAuth, act.Authorization)

			case LinkAuth{}.GetName():
				LinkAuth := LinkAuth{}
				rlp.DecodeBytes(act.Data, &LinkAuth)
				a.CheckLinkAuthAuthorization(LinkAuth, act.Authorization)

			case UnLinkAuth{}.GetName():
				UnLinkAuth := UnLinkAuth{}
				rlp.DecodeBytes(act.Data, &UnLinkAuth)
				a.CheckUnLinkAuthAuthorization(UnLinkAuth, act.Authorization)

			case CancelDelay{}.GetName():
				CancelDelay := CancelDelay{}
				rlp.DecodeBytes(act.Data, &CancelDelay)
				a.CheckCancelDelayAuthorization(CancelDelay, act.Authorization)

			default:
				specialCase = false
			}
		}

		for _, declaredAuth := range act.Authorization {
			(*checkTime)()
			if !specialCase {
				minPermissionName := a.LookupMinimumPermission(declaredAuth.Actor, act.Account, act.Name)
				if minPermissionName != nil {
					minPermission := a.GetPermission(&types.PermissionLevel{Actor: declaredAuth.Actor, Permission: *minPermissionName})
					permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
					if err != nil {
						log.Error("CheckAuthorization is error: %s", err)
					}
					EosAssert(a.GetPermission(&declaredAuth).Satisfies(*minPermission, permissionIndex), &IrrelevantAuthException{},
						"action declares irrelevant authority '%v'; minimum authority is %v", declaredAuth, types.PermissionLevel{minPermission.Owner, minPermission.Name})
				}
			}

			isExist := false
			for first, second := range permissionToSatisfy {
				if first == declaredAuth {
					if second > delay {
						second = delay
						isExist = true
						break
					}
				}
			}
			if !isExist {
				permissionToSatisfy[declaredAuth] = delay
			}
		}
	}
	for p, q := range permissionToSatisfy {
		(*checkTime)()
		EosAssert(checker.SatisfiedLoc(&p, q, nil), &UnsatisfiedAuthorization{},
			"transaction declares authority '%v', "+
				"but does not have signatures for it under a provided delay of %v ms, "+
				"provided permissions %v, and provided keys %v", p, providedDelay.Count()/1000, providedPermissions, providedKeys)
	}
	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "transaction bears irrelevant signatures from these keys: %v", checker.GetUnusedKeys())
	}
}

func (a *AuthorizationManager) CheckAuthorization2(account common.AccountName,
	permission common.PermissionName,
	providedKeys *treeset.Set, //flat_set<public_key_type>
	providedPermissions *treeset.Set, //flat_set<permission_level>
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	delayMaxLimit := common.Seconds(int64(a.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	var effectiveProvidedDelay common.Microseconds
	if providedDelay >= delayMaxLimit {
		effectiveProvidedDelay = common.MaxMicroseconds()
	} else {
		effectiveProvidedDelay = providedDelay
	}
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority {
		perm := a.GetPermission(p)
		if perm != nil {
			return perm.Auth
		} else {
			return types.SharedAuthority{}
		}
	},
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		providedKeys,
		providedPermissions,
		effectiveProvidedDelay,
		checkTime)
	EosAssert(checker.SatisfiedLc(&types.PermissionLevel{account, permission}, nil), &UnsatisfiedAuthorization{},
		"permission '%v' was not satisfied under a provided delay of %v ms, provided permissions %v, and provided keys %v",
		types.PermissionLevel{account, permission}, providedDelay.Count()/1000, providedPermissions, providedKeys)

	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "irrelevant keys provided: %v", checker.GetUnusedKeys())
	}
}

func (a *AuthorizationManager) GetRequiredKeys(trx *types.Transaction,
	candidateKeys *treeset.Set,
	providedDelay common.Microseconds) treeset.Set {
	checker := types.MakeAuthChecker(
		func(p *types.PermissionLevel) types.SharedAuthority {
			perm := a.GetPermission(p)
			if perm != nil {
				return perm.Auth
			} else {
				return types.SharedAuthority{}
			}
		},
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		candidateKeys,
		treeset.NewWith(types.PermissionLevelType, types.ComparePermissionLevel),
		providedDelay,
		noopCheckTime,
	)
	for _, act := range trx.Actions {
		for _, declaredAuth := range act.Authorization {
			EosAssert(checker.SatisfiedLc(&declaredAuth, nil), &UnsatisfiedAuthorization{},
				"transaction declares authority '%v', but does not have signatures for it.", declaredAuth)
		}
	}
	return checker.GetUsedKeys()
}
