import type { ComponentProps } from 'react'

type ScriptProps = ComponentProps<'script'> & {
  strategy?: 'beforeInteractive' | 'afterInteractive' | 'lazyOnload'
}

export default function Script({ strategy: _, ...props }: ScriptProps) {
  return <script {...props} />
}
