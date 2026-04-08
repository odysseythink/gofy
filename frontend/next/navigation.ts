import {
  useLocation,
  useNavigate,
  useParams as useRouterParams,
  useSearchParams as useRouterSearchParams,
} from 'react-router'

export function useRouter() {
  const navigate = useNavigate()
  return {
    push: (href: string) => navigate(href),
    replace: (href: string) => navigate(href, { replace: true }),
    back: () => navigate(-1),
    forward: () => navigate(1),
    refresh: () => window.location.reload(),
    prefetch: (_href: string) => { /* no-op in SPA */ },
  }
}

export function usePathname() {
  return useLocation().pathname
}

export function useParams<T extends Record<string, string> = Record<string, string>>() {
  return useRouterParams() as T
}

export function useSearchParams() {
  const [searchParams] = useRouterSearchParams()
  return searchParams
}

export function useSelectedLayoutSegment(): string | null {
  const { pathname } = useLocation()
  const segments = pathname.split('/').filter(Boolean)
  return segments.length > 0 ? segments[segments.length - 1] : null
}

export function useSelectedLayoutSegments(): string[] {
  const { pathname } = useLocation()
  return pathname.split('/').filter(Boolean)
}
