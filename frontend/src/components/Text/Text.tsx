import { HeaderText } from './Text.styles'

type TextProps = {
  children: React.ReactNode
  className?: string
}

export const PageHeaderText = ({ children, className }: TextProps) => {
  return <HeaderText className={className}>{children}</HeaderText>
}

// small text
export const SmallText = ({ children, className }: TextProps) => {
  return <HeaderText className={className}>{children}</HeaderText>
}
