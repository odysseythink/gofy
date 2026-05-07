import type { FC, PropsWithChildren } from 'react'
import type { SystemFeatures } from '@/types/feature'
import { useQuery } from '@tanstack/react-query'
import { useEffect } from 'react'
import { create } from 'zustand'
import Loading from '@/app/components/base/loading'
import { consoleClient } from '@/service/client'
import { defaultSystemFeatures } from '@/types/feature'
import { fetchSetupStatusWithCache } from '@/utils/setup-status'

type GlobalPublicStore = {
  systemFeatures: SystemFeatures
  setSystemFeatures: (systemFeatures: SystemFeatures) => void
}

export const useGlobalPublicStore = create<GlobalPublicStore>(set => ({
  systemFeatures: defaultSystemFeatures,
  setSystemFeatures: (systemFeatures: SystemFeatures) => set(() => ({ systemFeatures })),
}))

const systemFeaturesQueryKey = ['systemFeatures'] as const
const setupStatusQueryKey = ['setupStatus'] as const

// Pure fetcher: no store side effects here. Writing to zustand inside queryFn
// causes a render-phase cascade under React 19 StrictMode + TanStack Query v5
// that loops forever. Sync the store in a useEffect instead.
async function fetchSystemFeatures() {
  return consoleClient.systemFeatures()
}

export function useSystemFeaturesQuery() {
  return useQuery({
    queryKey: systemFeaturesQueryKey,
    queryFn: fetchSystemFeatures,
  })
}

export function useIsSystemFeaturesPending() {
  const { isPending } = useSystemFeaturesQuery()
  return isPending
}

export function useSetupStatusQuery() {
  return useQuery({
    queryKey: setupStatusQueryKey,
    queryFn: fetchSetupStatusWithCache,
    staleTime: Infinity,
  })
}

const GlobalPublicStoreProvider: FC<PropsWithChildren> = ({
  children,
}) => {
  const { isPending, data } = useSystemFeaturesQuery()
  useSetupStatusQuery()

  useEffect(() => {
    if (data)
      useGlobalPublicStore.getState().setSystemFeatures({ ...defaultSystemFeatures, ...data })
  }, [data])

  if (isPending)
    return <div className="flex h-screen w-screen items-center justify-center"><Loading /></div>
  return <>{children}</>
}
export default GlobalPublicStoreProvider
