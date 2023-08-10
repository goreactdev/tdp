import { Dialog, Transition } from '@headlessui/react'
import type { ReactNode } from 'react'
import React, { Fragment, useState } from 'react'

import { ModalContent, ModalTitle, Overlay } from './Modal.styles'

interface ModalProps {
  trigger: ReactNode
  title: string
  children: ReactNode | ((closeModal: () => void) => ReactNode)
}

const Modal: React.FC<ModalProps> = ({ trigger, title, children }) => {
  const [isOpen, setIsOpen] = useState(false)

  function closeModal() {
    setIsOpen(false)
  }

  function openModal() {
    setIsOpen(true)
  }

  return (
    <>
      <div onClick={openModal} className="cursor-pointer">
        {trigger}
      </div>

      <Transition as={Fragment} appear show={isOpen}>
        <Dialog className="relative z-10" onClose={closeModal}>
          <Transition.Child
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Overlay />
          </Transition.Child>

          <div className="fixed inset-0 overflow-y-auto">
            <div className="flex min-h-full items-center justify-center p-4 text-center">
              <Transition.Child
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <ModalContent>
                  <ModalTitle as="h3">{title}</ModalTitle>
                  <div className="mt-2">
                    {typeof children === 'function'
                      ? children(closeModal)
                      : children}
                  </div>
                </ModalContent>
              </Transition.Child>
            </div>
          </div>
        </Dialog>
      </Transition>
    </>
  )
}

export default Modal
