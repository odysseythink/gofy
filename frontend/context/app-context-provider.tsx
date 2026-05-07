
import type { FC, ReactNode } from 'react'
import type { ICurrentWorkspace, OdysseythinkVersionResponse, UserProfileResponse } from '@/models/common'
import { useQueryClient } from '@tanstack/react-query'
import { useCallback, useEffect, useMemo } from 'react'
import { setUserId, setUserProperties } from '@/app/components/base/amplitude'
import { setZendeskConversationFields } from '@/app/components/base/zendesk/utils'
import MaintenanceNotice from '@/app/components/header/maintenance-notice'
import { ZENDESK_FIELD_IDS } from '@/config'
import {
  AppContext,
  initialOdysseythinkVersionInfo,
  initialWorkspaceInfo,
  userProfilePlaceholder,
  useSelector,
} from '@/context/app-context'
import { env } from '@/env'
import {
  useCurrentWorkspace,
  useOdysseythinkVersion,
  useUserProfile,
} from '@/service/use-common'
import { useGlobalPublicStore } from './global-public-context'

export type AppContextProviderProps = {
  children: ReactNode
}

export const AppContextProvider: FC<AppContextProviderProps> = ({ children }) => {
  const queryClient = useQueryClient()
  const systemFeatures = useGlobalPublicStore(s => s.systemFeatures)
  const { data: userProfileResp } = useUserProfile()
  const { data: currentWorkspaceResp, isPending: isLoadingCurrentWorkspace, isFetching: isValidatingCurrentWorkspace } = useCurrentWorkspace()
  const odysseythinkVersionQuery = useOdysseythinkVersion(
    userProfileResp?.meta.currentVersion,
    !systemFeatures.branding.enabled,
  )

  const userProfile = useMemo<UserProfileResponse>(() => userProfileResp?.profile || userProfilePlaceholder, [userProfileResp?.profile])
  const currentWorkspace = useMemo<ICurrentWorkspace>(() => currentWorkspaceResp || initialWorkspaceInfo, [currentWorkspaceResp])
  const odysseythinkVersionInfo = useMemo<OdysseythinkVersionResponse>(() => {
    if (!userProfileResp?.meta?.currentVersion || !odysseythinkVersionQuery.data)
      return initialOdysseythinkVersionInfo

    const current_version = userProfileResp.meta.currentVersion
    const current_env = userProfileResp.meta.currentEnv || ''
    const versionData = odysseythinkVersionQuery.data
    return {
      ...versionData,
      current_version,
      latest_version: versionData.version,
      current_env,
    }
  }, [odysseythinkVersionQuery.data, userProfileResp?.meta])

  const isCurrentWorkspaceManager = useMemo(() => ['owner', 'admin'].includes(currentWorkspace.role), [currentWorkspace.role])
  const isCurrentWorkspaceOwner = useMemo(() => currentWorkspace.role === 'owner', [currentWorkspace.role])
  const isCurrentWorkspaceEditor = useMemo(() => ['owner', 'admin', 'editor'].includes(currentWorkspace.role), [currentWorkspace.role])
  const isCurrentWorkspaceDatasetOperator = useMemo(() => currentWorkspace.role === 'dataset_operator', [currentWorkspace.role])

  const mutateUserProfile = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ['common', 'user-profile'] })
  }, [queryClient])

  const mutateCurrentWorkspace = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ['common', 'current-workspace'] })
  }, [queryClient])

  // #region Zendesk conversation fields
  useEffect(() => {
    if (ZENDESK_FIELD_IDS.ENVIRONMENT && odysseythinkVersionInfo?.current_env) {
      setZendeskConversationFields([{
        id: ZENDESK_FIELD_IDS.ENVIRONMENT,
        value: odysseythinkVersionInfo.current_env.toLowerCase(),
      }])
    }
  }, [odysseythinkVersionInfo?.current_env])

  useEffect(() => {
    if (ZENDESK_FIELD_IDS.VERSION && odysseythinkVersionInfo?.version) {
      setZendeskConversationFields([{
        id: ZENDESK_FIELD_IDS.VERSION,
        value: odysseythinkVersionInfo.version,
      }])
    }
  }, [odysseythinkVersionInfo?.version])

  useEffect(() => {
    if (ZENDESK_FIELD_IDS.EMAIL && userProfile?.email) {
      setZendeskConversationFields([{
        id: ZENDESK_FIELD_IDS.EMAIL,
        value: userProfile.email,
      }])
    }
  }, [userProfile?.email])

  useEffect(() => {
    if (ZENDESK_FIELD_IDS.WORKSPACE_ID && currentWorkspace?.id) {
      setZendeskConversationFields([{
        id: ZENDESK_FIELD_IDS.WORKSPACE_ID,
        value: currentWorkspace.id,
      }])
    }
  }, [currentWorkspace?.id])
  // #endregion Zendesk conversation fields

  useEffect(() => {
    // Report user and workspace info to Amplitude when loaded
    if (userProfile?.id) {
      setUserId(userProfile.email)
      const properties: Record<string, string | number | boolean> = {
        email: userProfile.email,
        name: userProfile.name,
        has_password: userProfile.is_password_set,
      }

      if (currentWorkspace?.id) {
        properties.workspace_id = currentWorkspace.id
        properties.workspace_name = currentWorkspace.name
        properties.workspace_plan = currentWorkspace.plan
        properties.workspace_status = currentWorkspace.status
        properties.workspace_role = currentWorkspace.role
      }

      setUserProperties(properties)
    }
  }, [userProfile, currentWorkspace])

  return (
    <AppContext.Provider value={{
      userProfile,
      mutateUserProfile,
      odysseythinkVersionInfo,
      useSelector,
      currentWorkspace,
      isCurrentWorkspaceManager,
      isCurrentWorkspaceOwner,
      isCurrentWorkspaceEditor,
      isCurrentWorkspaceDatasetOperator,
      mutateCurrentWorkspace,
      isLoadingCurrentWorkspace,
      isValidatingCurrentWorkspace,
    }}
    >
      <div className="flex h-full flex-col overflow-y-auto">
        {env.VITE_MAINTENANCE_NOTICE && <MaintenanceNotice />}
        <div className="relative flex grow flex-col overflow-y-auto overflow-x-hidden bg-background-body">
          {children}
        </div>
      </div>
    </AppContext.Provider>
  )
}
