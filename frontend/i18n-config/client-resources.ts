import type { Resource, ResourceLanguage } from 'i18next'
import type { Locale } from '.'
import { camelCase } from 'es-toolkit/string'
import { namespacesInFileName } from './resources'

const localeModules = import.meta.glob('../i18n/**/*.json', { eager: true }) as Record<string, { default: Record<string, string> }>

export function getResources(locale: Locale): Resource {
  const messages: ResourceLanguage = {}

  for (const ns of namespacesInFileName) {
    const key = `../i18n/${locale}/${ns}.json`
    const mod = localeModules[key]
    if (mod) {
      messages[camelCase(ns)] = mod.default
    }
  }

  return { [locale]: messages }
}
