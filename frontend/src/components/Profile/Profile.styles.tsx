import { FaUserAlt } from 'react-icons/fa'
import { RiVipDiamondFill } from 'react-icons/ri'
import styled from 'styled-components'
import tw from 'twin.macro'

export const CardContainer = styled.div`
  ${tw`w-full flex flex-col items-center space-y-4 divide-y rounded-2xl bg-backgroundGray hover:transition-all`}
`

export const Wrapper = styled.div`
  ${tw`w-full flex flex-col space-y-5 col-span-2`}
`

export const CardHeader = styled.div`
  ${tw`flex flex-col cursor-pointer space-y-3 pb-1 items-center w-full justify-center`}
`

interface AvatarProps {
  imageUrl: string
  size: 'small' | 'medium' | 'large'
}

export const Avatar = styled.div<AvatarProps>`
  ${tw`cursor-pointer hover:scale-105 transition-all duration-300 rounded-full bg-cover bg-center`}
  background-image: ${({ imageUrl }) => `url(${imageUrl})`};

  ${({ size }) => {
    if (size === 'small') {
      return tw`h-9 w-9 sm:h-9 sm:w-9 md:h-9 md:w-9`
    } else if (size === 'medium') {
      return tw`h-32 w-32 sm:h-40 sm:w-40 md:h-48 md:w-48`
    } else if (size === 'large') {
      return tw`h-40 w-40 sm:h-32 sm:w-32 md:h-40 md:w-40`
    }
  }}
`
export const Name = styled.h2`
  ${tw`text-xl font-bold text-mainColor`}
`

export const Username = styled.p`
  ${tw`text-center font-medium text-gray-500`}
`

export const JobTitle = styled.p`
  ${tw`text-center text-sm font-medium text-backgroundBlack`}
`

export const CardFooter = styled.div`
  ${tw`flex w-10/12 flex-col justify-center`}
`

export const InfoContainer = styled.div`
  ${tw`my-4 space-y-4 text-sm text-backgroundBlack`}
`

export const InfoRow = styled.div`
  ${tw`flex items-center justify-between`}
`

export const AwardsContainer = styled.div`
  ${tw`py-6 flex space-y-4 flex-col w-10/12`}
`

export const AwardsTitle = styled.h2`
  ${tw`font-semibold flex items-center text-mainColor`}
`

export const AwardsText = styled.p`
  ${tw`text-base`}
`

export const DescriptionContainer = styled.div`
  ${tw`pt-6 flex space-y-4 flex-col w-10/12`}
`

export const DescriptionTitle = styled.h2`
  ${tw`font-semibold`}
`

export const DescriptionText = styled.p`
  ${tw`text-base`}
`

export const LanguagesContainer = styled.div`
  ${tw`pt-6 flex space-y-4 flex-col w-10/12`}
`

export const LanguagesTitle = styled.h2`
  ${tw`font-semibold`}
`

export const LanguagesBlock = styled.div`
  ${tw`text-base flex flex-col`}
`

export const LinkedAccountsContainer = styled.div`
  ${tw`pt-6 flex space-y-3 flex-col w-10/12`}
`

export const LinkedAccountsTitle = styled.h2`
  ${tw`font-semibold`}
`

export const LinkedAccount = styled.div`
  ${tw`flex items-center hover:text-mainColor transition-all duration-300`}
`

export const LinkedAccountIcon = styled.div`
  ${tw`mr-2 `}
`

export const LinkedAccountText = styled.a`
  ${tw`ml-2 `}
`

export const CertificationsContainer = styled.div`
  ${tw`py-6 flex space-y-3 flex-col w-10/12`}
`

export const CertificationsTitle = styled.h2`
  ${tw`font-semibold`}
`

export const CertificationsText = styled.p`
  ${tw`text-gray-600 font-medium`}
`

export const CertificationsBlock = styled.div`
  ${tw`text-gray-400 font-medium flex flex-col`}
`

// diamonds

export const DiamondContainer = styled.div`
  ${tw`flex items-center justify-between`}
`

export const Diamond = styled.span`
  ${tw`cursor-pointer px-0.5 transition-all hover:brightness-110`}
`

export const FilledDiamond = styled(RiVipDiamondFill)`
  ${tw`fill-mainColor`}
`

export const EmptyDiamond = styled(RiVipDiamondFill)`
  ${tw`fill-gray-400`}
`

export const AwardCount = styled.div`
  ${tw`ml-1 text-sm font-bold`}
`

export const AwardCountText = styled.span`
  ${tw`font-medium text-gray-500`}
`

type SizeType = 'small' | 'medium' | 'large'

type DefaultAvatarProps = {
  size: SizeType
}

export const AvatarContainer = styled.div`
  ${tw`bg-mainColor rounded-full`}
`

export const UserIcon = styled(FaUserAlt)(({ size }: DefaultAvatarProps) => [
  tw`inline-flex justify-center items-center cursor-pointer fill-white transition-transform duration-300 hover:scale-105`,
  size === 'small' && tw`w-12 h-12 p-2.5`,
  size === 'medium' && tw`w-48 h-48 sm:w-32 sm:h-32 p-9 sm:p-6`,
  size === 'large' && tw`w-40 h-40 p-8`,
])
