import { Link } from 'react-router-dom'

import { Button } from './Button.styles'

type ButtonContainerProps = {
  children: React.ReactNode
  color: 'blue' | 'white' | 'red' | 'yellow'
  to?: string
  href?: string
  // add default button props
} & React.ButtonHTMLAttributes<HTMLButtonElement>

const ButtonContainer = ({
  children,
  color,
  to,
  href,
  ...props
}: ButtonContainerProps) => {
  if (to) {
    return (
      <Link to={to}>
        <Button {...props} color={color}>
          {children}
        </Button>
      </Link>
    )
  }

  if (href) {
    return (
      <a href={href} target="_blank">
        <Button {...props} color={color}>
          {children}
        </Button>
      </a>
    )
  }

  return (
    <Button {...props} color={color}>
      {children}
    </Button>
  )
}

export default ButtonContainer
