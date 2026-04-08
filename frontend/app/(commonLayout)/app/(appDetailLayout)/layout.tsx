import * as React from 'react'
import { useTranslation } from 'react-i18next'
import { Outlet } from 'react-router'
import useDocumentTitle from '@/hooks/use-document-title'

const AppDetail = () => {
  const { t } = useTranslation()
  useDocumentTitle(t('menus.appDetail', { ns: 'common' }))

  return (
    <Outlet />
  )
}

export default React.memo(AppDetail)
