import { RouterProvider } from 'react-router'
import { Provider as JotaiProvider } from 'jotai/react'
import { ThemeProvider } from 'next-themes'
import { NuqsAdapter } from 'nuqs/adapters/react'
import GlobalPublicStoreProvider from '@/context/global-public-context'
import { TanstackQueryInitializer } from '@/context/query-client'
import { I18nProvider } from '@/i18n-config/provider'
import { ToastProvider } from '@/app/components/base/toast'
import { ToastHost } from '@/app/components/base/ui/toast'
import { TooltipProvider } from '@/app/components/base/ui/tooltip'
import { router } from './router'

export function App() {
  return (
    <JotaiProvider>
      <ThemeProvider
        attribute="data-theme"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
        enableColorScheme={false}
      >
        <NuqsAdapter>
          <TanstackQueryInitializer>
            <I18nProvider>
              <ToastHost timeout={5000} limit={3} />
              <ToastProvider>
                <GlobalPublicStoreProvider>
                  <TooltipProvider delay={300} closeDelay={200}>
                    <RouterProvider router={router} />
                  </TooltipProvider>
                </GlobalPublicStoreProvider>
              </ToastProvider>
            </I18nProvider>
          </TanstackQueryInitializer>
        </NuqsAdapter>
      </ThemeProvider>
    </JotaiProvider>
  )
}
