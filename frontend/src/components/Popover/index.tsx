import { Popover, Transition } from '@headlessui/react'
import type { ReactElement, ReactNode } from 'react'
import React from 'react'
import { Link } from 'react-router-dom'

type PopoverItemBase = {
  name: string
  description: string
  icon: ReactElement | null
  allWithoutIcon?: boolean
}

type PopoverLinkItem = PopoverItemBase & {
  type: 'link'
  href: string
}

type PopoverButtonItem = PopoverItemBase & {
  type: 'button'
  onClick: () => void
}

type PopoverExternalLinkItem = PopoverItemBase & {
  type: 'external_link'
  href: string
}

export type PopoverItem =
  | PopoverLinkItem
  | PopoverButtonItem
  | PopoverExternalLinkItem

const PopoverContent = ({
  name,
  icon,
  description,
  allWithoutIcon,
}: PopoverItemBase) => {
  return (
    <>
      <div
        className={`${
          allWithoutIcon ? 'hidden' : 'h-10 w-10 sm:h-12 sm:w-12 '
        } mx-4 flex shrink-0 items-center justify-center text-white `}
      >
        {icon}
      </div>
      <div>
        <p
          className={
            (allWithoutIcon ? 'px-4 ' : '') +
            'text-sm font-medium text-gray-800'
          }
        >
          {name}
        </p>
        <p className="text-sm text-gray-500">{description}</p>
      </div>
    </>
  )
}

const PopoverElement = ({
  children,
  items,
  header,
}: {
  header?: ReactNode
  children: React.ReactNode
  items: PopoverItem[]
}) => {
  const allWithoutIcon = items.every((item) => item.icon === null)

  return (
    <Popover className="relative">
      {children}
      <Transition
        enter="transition ease-out duration-200"
        enterFrom="opacity-0 translate-y-1"
        enterTo="opacity-100 translate-y-0"
        leave="transition ease-in duration-150"
        leaveFrom="opacity-100 translate-y-0"
        leaveTo="opacity-0 translate-y-1"
      >
        <Popover.Panel
          className={`${
            allWithoutIcon ? 'w-[12rem]' : ' left-1/2 w-screen max-w-xs '
          } absolute  z-10 mt-3 -translate-x-1/2 transform px-4 sm:px-0`}
        >
          {({ close }) => (
            <div
              onClick={() => close()}
              className="overflow-hidden rounded-2xl shadow-xl ring-1 ring-black ring-opacity-5"
            >
              {header && header}
              <div className="relative flex flex-col bg-white  py-2 ">
                {items.map((item) => {
                  if (item.type === 'external_link') {
                    return (
                      <a
                        key={item.name}
                        href={item.href}
                        target="_blank"
                        className="flex items-center py-2 transition duration-150 ease-in-out hover:bg-gray-50 focus:outline-none focus-visible:ring focus-visible:ring-orange-500 focus-visible:ring-opacity-50"
                      >
                        <PopoverContent
                          {...item}
                          allWithoutIcon={allWithoutIcon}
                        />
                      </a>
                    )
                  }
                  if (item.type === 'link' || item.type === undefined) {
                    return (
                      <Link
                        key={item.name}
                        to={item.href}
                        className="flex items-center py-2  transition duration-150 ease-in-out hover:bg-gray-50 focus:outline-none focus-visible:ring focus-visible:ring-orange-500 focus-visible:ring-opacity-50"
                      >
                        <PopoverContent
                          {...item}
                          allWithoutIcon={allWithoutIcon}
                        />
                      </Link>
                    )
                  } else {
                    return (
                      <button
                        key={item.name}
                        onClick={item.onClick}
                        className="flex items-center py-2 transition duration-150 ease-in-out hover:bg-gray-50 focus:outline-none focus-visible:ring focus-visible:ring-orange-500 focus-visible:ring-opacity-50"
                      >
                        <PopoverContent
                          {...item}
                          allWithoutIcon={allWithoutIcon}
                        />
                      </button>
                    )
                  }
                })}
              </div>
            </div>
          )}
        </Popover.Panel>
      </Transition>
    </Popover>
  )
}

export default PopoverElement
