declare module 'react-textfit' {
  import type { ComponentType, CSSProperties, ReactNode } from 'react'

  interface TextfitProps {
    mode?: 'single' | 'multi'
    min?: number
    max?: number
    style?: CSSProperties
    className?: string
    children?: ReactNode
  }

  export const Textfit: ComponentType<TextfitProps>
}
