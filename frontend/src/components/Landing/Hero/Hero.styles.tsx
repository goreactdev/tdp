import styled from 'styled-components'
import tw from 'twin.macro'

export const Section = styled.section`
  ${tw`flex flex-col justify-between gap-6 sm:gap-10 md:mb-16 md:gap-16 lg:flex-row`}
`

export const ContentContainer = styled.div`
  ${tw`flex flex-col justify-center sm:text-center lg:py-12 lg:text-left xl:w-6/12`}
`

export const Title = styled.h1`
  ${tw`mb-8 text-center lg:text-left text-4xl font-bold !leading-tight sm:text-5xl md:text-6xl`}

  span {
    ${tw`bg-gradient-to-r from-mainColor to-[#6fb4f4] bg-clip-text text-transparent`}
  }
`

export const ImageContainer = styled.div`
  ${tw`h-[16rem] overflow-hidden rounded-3xl lg:h-[35rem]`}
`

export const Image = styled.img`
  ${tw`h-full w-full object-cover object-center`}
`
