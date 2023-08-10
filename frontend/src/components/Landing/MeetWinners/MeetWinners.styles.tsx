import { Link } from 'react-router-dom'
import styled from 'styled-components'
import tw from 'twin.macro'

export const Title = styled.h2`
  ${tw`text-2xl font-bold text-gray-800 lg:text-4xl`}
`

export const ExploreLink = styled(Link)`
  ${tw`font-extrabold text-lg text-mainColor flex items-center justify-center`}
`

export const WinnersGrid = styled.div`
  ${tw`grid grid-cols-1  gap-4 md:grid-cols-3 lg:grid-cols-4 lg:gap-8`}
`

export const WinnerCard = styled(Link)`
  ${tw`flex flex-col hover:scale-105 transition-all duration-300 items-center cursor-pointer rounded-2xl bg-backgroundGray p-4 hover:bg-gray-100 lg:p-8`}
`

export const ProfilePictureWrapper = styled.div`
  ${tw`mb-2 h-48 w-48 sm:h-24 sm:w-24 overflow-hidden  rounded-full bg-gray-200 shadow-lg md:mb-4 md:h-32 md:w-32 transition-all duration-300`}
`

export const ProfilePicture = styled.img`
  ${tw`h-full w-full object-cover  object-center`}
`

export const Name = styled.div`
  ${tw`text-center font-bold text-mainColor md:text-lg`}
`

export const Position = styled.p`
  ${tw`mb-3 text-center text-sm text-gray-800 md:mb-4 md:text-base`}
`

export const SocialLinksWrapper = styled.div`
  ${tw`flex justify-center`}
`
