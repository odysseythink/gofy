import type { Locale } from '.'
import { useState } from 'react'
import { I18nextProvider } from 'react-i18next'
import Cookies from 'js-cookie'
import { createI18nextInstance } from './client'
import { getResources } from './client-resources'
import { i18n } from '.'
import { LanguagesSupported } from './language'

function detectLocale(): Locale {
  const cookieLocale = Cookies.get('locale')
  if (cookieLocale && LanguagesSupported.includes(cookieLocale as Locale)) {
    return cookieLocale as Locale
  }

  const browserLang = navigator.language
  const normalizedLang = browserLang.replace('_', '-')
  if (LanguagesSupported.includes(normalizedLang as Locale)) {
    return normalizedLang as Locale
  }

  const langPrefix = browserLang.split('-')[0]
  const matched = LanguagesSupported.find(l => l.startsWith(langPrefix))
  if (matched) {
    return matched as Locale
  }

  return i18n.defaultLocale as Locale
}

export function I18nProvider({ children }: { children: React.ReactNode }) {
  const [locale] = useState(detectLocale)
  const [resources] = useState(() => getResources(locale))
  const [i18nInstance] = useState(() => createI18nextInstance(locale, resources))

  return (
    <I18nextProvider i18n={i18nInstance}>
      {children}
    </I18nextProvider>
  )
}
