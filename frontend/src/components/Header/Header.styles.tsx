import { Link } from 'react-router-dom'
import styled from 'styled-components/macro'
import tw from 'twin.macro'

export const HeaderContainer = styled.header`
  ${tw`mb-8 z-20 flex mx-auto bg-white/70 backdrop-blur-md items-center fixed  w-full  justify-between py-4 md:mb-12 md:py-4 xl:mb-16`}
`

export const LogoLink = styled(Link)`
  ${tw`text-gray-800 inline-flex items-center gap-2.5 text-2xl font-extrabold md:text-3xl`}
`

export const LogoIcon = styled.svg`
  ${tw`h-auto w-6 text-mainColor`}
`

export const NavContainer = styled.nav`
  ${tw`hidden gap-2 lg:flex `}
`

export const NavLink = styled(Link)`
  ${tw`text-lg font-semibold text-gray-800 hover:bg-backgroundGray px-4 py-2 rounded-2xl transition-all duration-150`}
`

export const ButtonContainer = styled.div`
  ${tw`hidden space-x-2 lg:flex`}
`

export const MobileButton = styled.button`
  ${tw`inline-flex items-center gap-2 rounded-3xl  px-2.5 py-2 text-sm font-semibold text-gray-500 ring-indigo-300 active:text-gray-700 md:text-base lg:hidden`}
`

export const MobileIcon = styled.svg`
  ${tw`h-6 w-6`}
`
