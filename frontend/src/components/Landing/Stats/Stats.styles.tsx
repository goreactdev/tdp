import styled from 'styled-components'
import tw from 'twin.macro'

export const Section = styled.section`
  ${tw`flex flex-col items-center justify-between gap-10 border-t pt-8 lg:flex-row lg:gap-8`}
`

export const StatsContainer = styled.div`
  ${tw`-mx-6 grid grid-cols-2 gap-4 md:-mx-8 md:flex md:divide-x`}
`

export const Stat = styled.div`
  ${tw`px-6 md:px-8`}
`

export const StatNumber = styled.span`
  ${tw`block text-center font-bold text-mainColor md:text-left md:text-xl`}
`

export const StatLabel = styled.span`
  ${tw`block text-center font-semibold text-gray-800 md:text-left md:text-base`}
`

export const SocialContainer = styled.div`
  ${tw`flex items-center justify-center gap-4 lg:justify-start`}
`

export const SocialLabel = styled.span`
  ${tw`text-sm font-semibold uppercase tracking-widest text-gray-400 sm:text-base`}
`

export const Divider = styled.span`
  ${tw`h-px w-12 bg-gray-200`}
`
