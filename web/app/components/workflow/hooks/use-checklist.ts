import {
  useCallback,
  useMemo,
  useRef,
} from 'react'
import { useTranslation } from 'react-i18next'
import { useStoreApi } from 'reactflow'
import type {
  CommonNodeType,
  Edge,
  Node,
  ValueSelector,
} from '../types'
import { BlockEnum } from '../types'
import { useStore } from '../store'
import {
  getValidTreeNodes,
} from '../utils'
import {
  CUSTOM_NODE,
} from '../constants'
import { useIsChatMode } from './use-workflow'
import { useNodesExtraData } from './use-nodes-data'
import { useToastContext } from '@/app/components/base/toast'
import { useGetLanguage } from '@/context/i18n'
import { MAX_TREE_DEPTH } from '@/config'
import useNodesAvailableVarList from './use-nodes-available-var-list'
import { getNodeUsedVars, isConversationVar, isENV, isSystemVar } from '../nodes/_base/components/variable/utils'

export const useChecklist = (nodes: Node[], edges: Edge[]) => {
  const { t } = useTranslation()
  const language = useGetLanguage()
  const nodesExtraData = useNodesExtraData()
  const isChatMode = useIsChatMode()

  const map = useNodesAvailableVarList(nodes)

  const getCheckData = useCallback((data: CommonNodeType<{}>) => {
    let checkData = data
    return checkData
  }, [])

  const needWarningNodes = useMemo(() => {
    const list = []
    const { validNodes } = getValidTreeNodes(nodes.filter(node => node.type === CUSTOM_NODE), edges)

    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i]
      let moreDataForCheckValid
      let usedVars: ValueSelector[] = []

      usedVars = getNodeUsedVars(node).filter(v => v.length > 0)

      if (node.type === CUSTOM_NODE) {
        const checkData = getCheckData(node.data)
        let { errorMessage } = nodesExtraData[node.data.type].checkValid(checkData, t, moreDataForCheckValid)

        if (!errorMessage) {
          const availableVars = map[node.id].availableVars

          for (const variable of usedVars) {
            const isEnv = isENV(variable)
            const isConvVar = isConversationVar(variable)
            const isSysVar = isSystemVar(variable)
            if (!isEnv && !isConvVar && !isSysVar) {
              const usedNode = availableVars.find(v => v.nodeId === variable?.[0])
              if (usedNode) {
                const usedVar = usedNode.vars.find(v => v.variable === variable?.[1])
                if (!usedVar)
                  errorMessage = t('workflow.errorMsg.invalidVariable')
              }
              else {
                errorMessage = t('workflow.errorMsg.invalidVariable')
              }
            }
          }
        }

        if (errorMessage || !validNodes.find(n => n.id === node.id)) {
          list.push({
            id: node.id,
            type: node.data.type,
            title: node.data.title,
            unConnected: !validNodes.find(n => n.id === node.id),
            errorMessage,
          })
        }
      }
    }

    if (isChatMode && !nodes.find(node => node.data.type === BlockEnum.Answer)) {
      list.push({
        id: 'answer-need-added',
        type: BlockEnum.Answer,
        title: t('workflow.blocks.answer'),
        errorMessage: t('workflow.common.needAnswerNode'),
      })
    }

    if (!isChatMode && !nodes.find(node => node.data.type === BlockEnum.End)) {
      list.push({
        id: 'end-need-added',
        type: BlockEnum.End,
        title: t('workflow.blocks.end'),
        errorMessage: t('workflow.common.needEndNode'),
      })
    }

    return list
  }, [nodes, edges, isChatMode, language, nodesExtraData, t, getCheckData])

  return needWarningNodes
}

export const useChecklistBeforePublish = () => {
  const { t } = useTranslation()
  const language = useGetLanguage()
  const { notify } = useToastContext()
  const isChatMode = useIsChatMode()
  const store = useStoreApi()
  const nodesExtraData = useNodesExtraData()

  const getCheckData = useCallback((data: CommonNodeType<{}>) => {
    let checkData = data

    return checkData
  }, [])

  const handleCheckBeforePublish = useCallback(async () => {
    const {
      getNodes,
      edges,
    } = store.getState()
    const nodes = getNodes().filter(node => node.type === CUSTOM_NODE)
    const {
      validNodes,
      maxDepth,
    } = getValidTreeNodes(nodes.filter(node => node.type === CUSTOM_NODE), edges)

    if (maxDepth > MAX_TREE_DEPTH) {
      notify({ type: 'error', message: t('workflow.common.maxTreeDepth', { depth: MAX_TREE_DEPTH }) })
      return false
    }
    // Before publish, we need to fetch datasets detail, in case of the settings of datasets have been changed
    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i]
      let moreDataForCheckValid

      const checkData = getCheckData(node.data)
      const { errorMessage } = nodesExtraData[node.data.type as BlockEnum].checkValid(checkData, t, moreDataForCheckValid)

      if (errorMessage) {
        notify({ type: 'error', message: `[${node.data.title}] ${errorMessage}` })
        return false
      }

      if (!validNodes.find(n => n.id === node.id)) {
        notify({ type: 'error', message: `[${node.data.title}] ${t('workflow.common.needConnectTip')}` })
        return false
      }
    }

    if (isChatMode && !nodes.find(node => node.data.type === BlockEnum.Answer)) {
      notify({ type: 'error', message: t('workflow.common.needAnswerNode') })
      return false
    }

    if (!isChatMode && !nodes.find(node => node.data.type === BlockEnum.End)) {
      notify({ type: 'error', message: t('workflow.common.needEndNode') })
      return false
    }

    return true
  }, [store, isChatMode, notify, t, language, nodesExtraData, getCheckData])

  return {
    handleCheckBeforePublish,
  }
}
