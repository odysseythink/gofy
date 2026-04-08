import { Outlet, useParams } from 'react-router'
import Main from './layout-main'

const DatasetDetailLayout = () => {
  const { datasetId } = useParams<{ datasetId: string }>()
  return <Main datasetId={datasetId!}><Outlet /></Main>
}
export default DatasetDetailLayout
