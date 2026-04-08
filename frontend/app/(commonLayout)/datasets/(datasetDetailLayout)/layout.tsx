import * as React from 'react'
import { Outlet } from 'react-router'

const DatasetDetail = () => {
  return (
    <Outlet />
  )
}

export default React.memo(DatasetDetail)
