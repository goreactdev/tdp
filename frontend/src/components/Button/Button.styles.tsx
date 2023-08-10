import styled, { css } from 'styled-components'
import tw from 'twin.macro'

export const Button = styled.button`
  ${tw`rounded-3xl hover:scale-105 hover:duration-300 active:scale-95 hover:transition-all px-8 py-3 text-center text-sm font-bold outline-none focus:outline-none md:text-base`}

  ${({ color }) =>
    color === 'blue' &&
    css`
      ${tw`bg-mainColor text-center text-sm font-semibold text-white outline-none transition duration-100 hover:brightness-110 active:brightness-125 md:text-base lg:inline-block`}
    `}


  ${({ color }) =>
    color === 'yellow' &&
    css`
      ${tw`bg-lime-500 text-center text-sm font-semibold text-white outline-none transition duration-100 hover:brightness-110 active:brightness-125 md:text-base lg:inline-block`}
    `}
  ${({ color }) =>
    color === 'red' &&
    css`
      ${tw`bg-red-500 text-center text-sm font-semibold text-white outline-none transition duration-100 hover:brightness-110 active:brightness-125 md:text-base lg:inline-block`}
    `}


  ${({ color }) =>
    color === 'white' &&
    css`
      ${tw`bg-white text-gray-800 border border-gray-200  hover:bg-backgroundGray active:bg-gray-100`}
    `}
`
