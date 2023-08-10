import styled from 'styled-components'
import tw from 'twin.macro'

export const ErrorPageWrapper = styled.div`
  ${tw`flex justify-center items-center h-screen`}
`

export const ContentWrapper = styled.div`
  ${tw`mx-auto max-w-screen-2xl px-4 md:px-8`}
`

export const LogoLink = styled.a`
  ${tw`text-black mb-8 inline-flex items-center gap-2.5 text-2xl font-bold md:text-3xl`}
  svg {
    ${tw`h-auto w-6 text-mainColor`}
  }
`

export const ErrorText = styled.p`
  ${tw`mb-4 text-sm font-semibold uppercase text-mainColor md:text-base`}
`

export const Heading = styled.h1`
  ${tw`mb-2 text-center text-2xl font-bold text-gray-800 md:text-3xl`}
`

export const Subheading = styled.p`
  ${tw`mb-12 max-w-screen-md text-center text-gray-500 md:text-lg `}
`
