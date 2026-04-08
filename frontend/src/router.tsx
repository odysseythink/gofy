import React, { Suspense } from 'react'
import { Navigate, createBrowserRouter } from 'react-router'
import type { RouteObject } from 'react-router'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/**
 * Wrap a lazy‑loaded **page** component in Suspense.
 */
function page(factory: () => Promise<{ default: React.ComponentType }>) {
  const Lazy = React.lazy(factory)
  return (
    <Suspense fallback={null}>
      <Lazy />
    </Suspense>
  )
}

/**
 * Wrap a lazy‑loaded **layout** component (renders an <Outlet />) in Suspense.
 */
function layout(factory: () => Promise<{ default: React.ComponentType }>) {
  const Lazy = React.lazy(factory)
  return (
    <Suspense fallback={null}>
      <Lazy />
    </Suspense>
  )
}

// ---------------------------------------------------------------------------
// Route tree
// ---------------------------------------------------------------------------

const routes: RouteObject[] = [
  // ── Root redirect ────────────────────────────────────────────────────
  {
    path: '/',
    element: <Navigate to="/apps" replace />,
  },

  // ── Common layout ────────────────────────────────────────────────────
  {
    element: layout(() => import('@/app/(commonLayout)/layout')),
    children: [
      { path: 'apps', element: page(() => import('@/app/(commonLayout)/apps/page')) },
      { path: 'tools', element: page(() => import('@/app/(commonLayout)/tools/page')) },
      { path: 'plugins', element: page(() => import('@/app/(commonLayout)/plugins/page')) },
      { path: 'education-apply', element: page(() => import('@/app/(commonLayout)/education-apply/page')) },

      // ── App detail ─────────────────────────────────────────────────
      {
        path: 'app',
        element: layout(() => import('@/app/(commonLayout)/app/(appDetailLayout)/layout')),
        children: [
          {
            path: ':appId',
            element: layout(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/layout')),
            children: [
              { path: 'overview', element: page(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/overview/page')) },
              { path: 'configuration', element: page(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/configuration/page')) },
              { path: 'develop', element: page(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/develop/page')) },
              { path: 'workflow', element: page(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/workflow/page')) },
              { path: 'logs', element: page(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/logs/page')) },
              { path: 'annotations', element: page(() => import('@/app/(commonLayout)/app/(appDetailLayout)/[appId]/annotations/page')) },
            ],
          },
        ],
      },

      // ── Datasets ───────────────────────────────────────────────────
      {
        path: 'datasets',
        element: layout(() => import('@/app/(commonLayout)/datasets/layout')),
        children: [
          { index: true, element: page(() => import('@/app/(commonLayout)/datasets/page')) },
          { path: 'create', element: page(() => import('@/app/(commonLayout)/datasets/create/page')) },
          { path: 'connect', element: page(() => import('@/app/(commonLayout)/datasets/connect/page')) },
          { path: 'create-from-pipeline', element: page(() => import('@/app/(commonLayout)/datasets/create-from-pipeline/page')) },

          // Dataset detail
          {
            element: layout(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/layout')),
            children: [
              {
                path: ':datasetId',
                element: layout(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/layout')),
                children: [
                  { path: 'documents', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/documents/page')) },
                  { path: 'documents/:documentId', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/documents/[documentId]/page')) },
                  { path: 'documents/:documentId/settings', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/documents/[documentId]/settings/page')) },
                  { path: 'documents/create', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/documents/create/page')) },
                  { path: 'documents/create-from-pipeline', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/documents/create-from-pipeline/page')) },
                  { path: 'api', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/api/page')) },
                  { path: 'pipeline', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/pipeline/page')) },
                  { path: 'hitTesting', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/hitTesting/page')) },
                  { path: 'settings', element: page(() => import('@/app/(commonLayout)/datasets/(datasetDetailLayout)/[datasetId]/settings/page')) },
                ],
              },
            ],
          },
        ],
      },

      // ── Explore ────────────────────────────────────────────────────
      {
        path: 'explore',
        element: layout(() => import('@/app/(commonLayout)/explore/layout')),
        children: [
          { path: 'apps', element: page(() => import('@/app/(commonLayout)/explore/apps/page')) },
          { path: 'installed/:appId', element: page(() => import('@/app/(commonLayout)/explore/installed/[appId]/page')) },
        ],
      },
    ],
  },

  // ── Share layout ─────────────────────────────────────────────────────
  {
    element: layout(() => import('@/app/(shareLayout)/layout')),
    children: [
      { path: 'chat/:token', element: page(() => import('@/app/(shareLayout)/chat/[token]/page')) },
      { path: 'chatbot/:token', element: page(() => import('@/app/(shareLayout)/chatbot/[token]/page')) },
      { path: 'completion/:token', element: page(() => import('@/app/(shareLayout)/completion/[token]/page')) },
      { path: 'workflow/:token', element: page(() => import('@/app/(shareLayout)/workflow/[token]/page')) },

      // Webapp signin
      {
        path: 'webapp-signin',
        element: layout(() => import('@/app/(shareLayout)/webapp-signin/layout')),
        children: [
          { index: true, element: page(() => import('@/app/(shareLayout)/webapp-signin/page')) },
          { path: 'check-code', element: page(() => import('@/app/(shareLayout)/webapp-signin/check-code/page')) },
        ],
      },

      // Webapp reset password
      {
        path: 'webapp-reset-password',
        element: layout(() => import('@/app/(shareLayout)/webapp-reset-password/layout')),
        children: [
          { index: true, element: page(() => import('@/app/(shareLayout)/webapp-reset-password/page')) },
          { path: 'check-code', element: page(() => import('@/app/(shareLayout)/webapp-reset-password/check-code/page')) },
          { path: 'set-password', element: page(() => import('@/app/(shareLayout)/webapp-reset-password/set-password/page')) },
        ],
      },
    ],
  },

  // ── Human input layout ───────────────────────────────────────────────
  {
    path: 'form/:token',
    element: page(() => import('@/app/(humanInputLayout)/form/[token]/page')),
  },

  // ── Auth: Signin ─────────────────────────────────────────────────────
  {
    path: 'signin',
    element: layout(() => import('@/app/signin/layout')),
    children: [
      { index: true, element: page(() => import('@/app/signin/page')) },
      { path: 'check-code', element: page(() => import('@/app/signin/check-code/page')) },
      { path: 'invite-settings', element: page(() => import('@/app/signin/invite-settings/page')) },
    ],
  },

  // ── Auth: Signup ─────────────────────────────────────────────────────
  {
    path: 'signup',
    element: layout(() => import('@/app/signup/layout')),
    children: [
      { index: true, element: page(() => import('@/app/signup/page')) },
      { path: 'check-code', element: page(() => import('@/app/signup/check-code/page')) },
      { path: 'set-password', element: page(() => import('@/app/signup/set-password/page')) },
    ],
  },

  // ── Auth: Reset password ─────────────────────────────────────────────
  {
    path: 'reset-password',
    element: layout(() => import('@/app/reset-password/layout')),
    children: [
      { index: true, element: page(() => import('@/app/reset-password/page')) },
      { path: 'check-code', element: page(() => import('@/app/reset-password/check-code/page')) },
      { path: 'set-password', element: page(() => import('@/app/reset-password/set-password/page')) },
    ],
  },

  // ── Account ──────────────────────────────────────────────────────────
  {
    path: 'account',
    element: layout(() => import('@/app/account/(commonLayout)/layout')),
    children: [
      { index: true, element: page(() => import('@/app/account/(commonLayout)/page')) },
      {
        path: 'oauth/authorize',
        element: layout(() => import('@/app/account/oauth/authorize/layout')),
        children: [
          { index: true, element: page(() => import('@/app/account/oauth/authorize/page')) },
        ],
      },
    ],
  },

  // ── Standalone pages ─────────────────────────────────────────────────
  { path: 'forgot-password', element: page(() => import('@/app/forgot-password/page')) },
  { path: 'init', element: page(() => import('@/app/init/page')) },
  { path: 'install', element: page(() => import('@/app/install/page')) },
  { path: 'activate', element: page(() => import('@/app/activate/page')) },
  { path: 'oauth-callback', element: page(() => import('@/app/oauth-callback/page')) },
]

// ---------------------------------------------------------------------------
// Export
// ---------------------------------------------------------------------------

export const router = createBrowserRouter(routes)
