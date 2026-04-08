import { lazy, Suspense, type ComponentType } from 'react'

type DynamicOptions = {
  ssr?: boolean
  loading?: () => React.ReactNode
}

export default function dynamic<P extends Record<string, any>>(
  importFn: () => Promise<{ default: ComponentType<P> }>,
  options?: DynamicOptions,
): ComponentType<P> {
  const LazyComponent = lazy(importFn)

  return function DynamicComponent(props: P) {
    const fallback = options?.loading?.() ?? null
    return (
      <Suspense fallback={fallback}>
        <LazyComponent {...props} />
      </Suspense>
    )
  } as ComponentType<P>
}
