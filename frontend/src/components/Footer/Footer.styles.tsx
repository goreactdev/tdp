import styled from 'styled-components'
import tw from 'twin.macro'

export const GridContainer = styled.div`
  ${tw`mb-16 grid grid-cols-2 gap-12 border-t pt-10 md:grid-cols-4 lg:grid-cols-6 lg:gap-8 lg:pt-12`}
`

export const LogoContainer = styled.div`
  ${tw`mb-4 lg:-mt-2`}
`

export const LogoLink = styled.a`
  ${tw`text-black inline-flex items-center gap-2 text-xl font-bold md:text-2xl`}
`

export const LogoSvg = styled.svg`
  ${tw`h-auto w-5 text-mainColor`}
`

export const Text = styled.p`
  ${tw`mb-6 text-gray-500 sm:pr-8`}
`

export const NavWrapper = styled.div`
  ${tw`flex flex-col gap-4`}
`

export const NavItem = styled.div`
  ${tw`text-sm`}
`

export const NavHeading = styled.h2`
  ${tw`mb-4 font-bold uppercase tracking-widest text-gray-800`}
`

export const NavLink = styled.a`
  ${tw`text-gray-500 transition duration-100 hover:text-mainColor active:text-indigo-600`}
`
