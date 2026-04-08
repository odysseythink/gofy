import * as React from 'react'
import { useTranslation } from 'react-i18next'
import { Outlet } from 'react-router'
import ExploreClient from '@/app/components/explore'
import useDocumentTitle from '@/hooks/use-document-title'

const ExploreLayout = () => {
  const { t } = useTranslation()
  useDocumentTitle(t('menus.explore', { ns: 'common' }))
  return (
    <ExploreClient>
      <Outlet />
    </ExploreClient>
  )
}

export default React.memo(ExploreLayout)
