import type { InitOptions } from 'i18next'
import { namespaces } from './resources'

export function getInitOptions(): InitOptions {
  return {
    // We do not have en for fallback
    load: 'currentOnly',
    fallbackLng: 'en-US',
    showSupportNotice: false,
    partialBundledLanguages: true,
    keySeparator: false,
    ns: namespaces,
    interpolation: {
      escapeValue: false,
    },
    // Limit react-i18next re-render triggers. Default `bindI18nStore: 'added'`
    // fires on every resource-add, which under React 19 StrictMode + Suspense
    // caused infinite render loops in components using useTranslation.
    react: {
      useSuspense: false,
      bindI18n: 'languageChanged',
      bindI18nStore: '',
    },
  }
}
