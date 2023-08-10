import { Dialog } from '@headlessui/react'
import styled from 'styled-components'
import tw from 'twin.macro'

export const Overlay = styled.div`
  ${tw`fixed inset-0 bg-black bg-opacity-25`}
`

export const ModalContent = styled(Dialog.Panel)`
  ${tw`w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all`}
`

export const ModalTitle = styled(Dialog.Title)`
  ${tw`text-lg font-medium leading-6 text-gray-900`}
`
