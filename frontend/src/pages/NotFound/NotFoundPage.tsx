import Button from '../../components/Button'

import {
  ContentWrapper,
  ErrorPageWrapper,
  ErrorText,
  Heading,
  LogoLink,
  Subheading,
} from './NotFoundPage.styles'

export const NotFoundPage = () => {
  return (
    <ErrorPageWrapper>
      <ContentWrapper>
        <div className="flex flex-col items-center">
          <LogoLink href="/" aria-label="logo">
            <svg
              width="95"
              height="94"
              viewBox="0 0 95 94"
              fill="currentColor"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M96 0V47L48 94H0V47L48 0H96Z" />
            </svg>
            TON Developers Platform
          </LogoLink>

          <ErrorText>That’s a 404</ErrorText>
          <Heading>Page not found</Heading>

          <Subheading>The page you’re looking for doesn’t exist.</Subheading>

          <Button to="/" color="blue">
            Go home
          </Button>
        </div>
      </ContentWrapper>
    </ErrorPageWrapper>
  )
}

export default NotFoundPage
