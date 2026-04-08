import { Outlet } from 'react-router'
import WebAppStoreProvider from '@/context/web-app-context'
import Splash from './components/splash'

const Layout = () => {
  return (
    <div className="h-full min-w-[300px] pb-[env(safe-area-inset-bottom)]">
      <WebAppStoreProvider>
        <Splash>
          <Outlet />
        </Splash>
      </WebAppStoreProvider>
    </div>
  )
}

export default Layout
