import styled from 'styled-components'
import tw from 'twin.macro'

export const Container = styled.div`
  ${tw`flex flex-col lg:grid grid-cols-7 gap-7 place-content-center`}
`

export const Content = styled.div`
  ${tw`col-span-5`}
`

export const InnerContainer = styled.div`
  ${tw`mx-auto max-w-screen-2xl`}
`

export const Grid = styled.div`
  ${tw`grid gap-4 grid-cols-2 lg:grid-cols-3 `}
`

export const Card = styled.div`${tw`cursor-pointer rounded-2xl `}}`

export const ImageContent = styled.div`
  ${tw`relative block w-full overflow-hidden rounded-t-2xl bg-backgroundGray  pt-[100%]`}
`

export const Img = styled.img`
  ${tw`absolute inset-0 h-full w-full object-cover object-center transition duration-500 hover:scale-110 group-hover:scale-110`}
`

export const Badge = styled.span<{ color: 'red' | 'blue' }>`
  ${tw`absolute cursor-pointer left-0 top-3 rounded-r-2xl select-none hover:scale-110 transition-all duration-300  px-3 py-1.5 text-sm font-semibold uppercase tracking-wider text-white`}
  ${({ color }) => (color === 'red' ? tw`bg-red-500` : tw`bg-blue-500`)}
`

export const CardContent = styled.a`
  ${tw`flex items-start justify-between gap-2 rounded-b-2xl bg-backgroundGray p-4`}
`

export const TextContainer = styled.div`
  ${tw`flex flex-col`}
`

export const CardTitle = styled.a`
  ${tw`font-bold text-gray-800 transition duration-100 hover:text-gray-500 lg:text-base`}
`

export const CardSubtitle = styled.span`
  ${tw`text-sm text-gray-500 lg:text-sm`}
`
